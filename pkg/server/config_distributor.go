package server

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

// ConfigStore represents API that server package using to get configuration parameters.
type ConfigStore interface {
	// Get returns current active configuration.
	Get() (*Config, error)

	// Subscribe subscribes given channel to configuration updates.
	// If checks are specified, new config will be sent to the channel only if one of them returns true.
	Subscribe(ch chan<- *Config, checks ...ConfigChangedCheckFn) error

	// Unsubscribe cancels subscription for the given channel.
	Unsubscribe(ch chan<- *Config)
}

// ConfigChangedCheckFn is a predicate function that must return true when previous config differs from the new one.
// We use this to notify config subscribers only about relevant changes.
type ConfigChangedCheckFn func(oldConfig, newConfig *Config) bool

// ConfigDistributor is an implementation of ConfigStore. It takes a channel that receives config updates and
// multiplexes it into multiple subscribed channels.
type ConfigDistributor struct {
	// Guards following
	mu sync.RWMutex

	// Fired after config monitoring goroutine is started
	started *Event

	// Fired after Close is called
	closed *Event

	// Channel with Config updates
	source <-chan *Config

	// Last received Config
	lastConfig *Config

	// Subscriptions to Config updates
	subscriptions map[chan<- *Config][]ConfigChangedCheckFn
}

var (
	// ErrConfigDistributorClosed is returned by the ConfigDistributor's Get and Subscribe methods after a call to Close.
	ErrConfigDistributorClosed = errors.New("config_distributor: ConfigDistributor is closed")
)

var _ ConfigStore = &ConfigDistributor{}

// NewConfigDistributor initializes a new ConfigDistributor.
func NewConfigDistributor(configChan <-chan *Config) *ConfigDistributor {
	return &ConfigDistributor{
		source:  configChan,
		started: NewEvent(),
		closed:  NewEvent(),
	}
}

// Get returns current active configuration. Get will block execution until new configuration arrive.
func (p *ConfigDistributor) Get() (*Config, error) {
	// Check that distributor is still active
	if p.closed.HasFired() {
		return nil, ErrConfigDistributorClosed
	}

	// Take lastConfig snapshot
	var lastConfigSnapshot *Config

	p.mu.RLock()
	{
		lastConfigSnapshot = p.lastConfig
	}
	p.mu.RUnlock()

	// Shortcut: return cached config
	if lastConfigSnapshot != nil {
		return lastConfigSnapshot, nil
	}

	// Wait for the next config update from the source
	return p.waitForNextConfig()
}

// Subscribe subscribes given channel to configuration updates.
// If checks are specified, new config will be sent to the channel only if one of them returns true.
func (p *ConfigDistributor) Subscribe(ch chan<- *Config, checks ...ConfigChangedCheckFn) error {
	// Check that distributor is still active
	if p.closed.HasFired() {
		return ErrConfigDistributorClosed
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Send current config immediately even before we registered a new subscription
	if p.lastConfig != nil && p.shouldSendUpdatedConfig(checks, nil, p.lastConfig) {
		ch <- p.lastConfig
	}

	// Register subscription
	if p.subscriptions == nil {
		p.subscriptions = make(map[chan<- *Config][]ConfigChangedCheckFn)
	}

	p.subscriptions[ch] = checks

	// Start monitoring goroutine
	if p.started.Fire() {
		go p.monitorConfig()
	}

	return nil
}

// Unsubscribe cancels subscription for the given channel.
func (p *ConfigDistributor) Unsubscribe(ch chan<- *Config) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.unsubscribeLocked(ch)
}

// Close stops the monitoring goroutine and cancels all subscriptions.
func (p *ConfigDistributor) Close() {
	if !p.closed.Fire() {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for ch := range p.subscriptions {
		p.unsubscribeLocked(ch)
	}
}

// unsubscribeLocked unsubscribes the channel without locking the mutex.
func (p *ConfigDistributor) unsubscribeLocked(ch chan<- *Config) {
	delete(p.subscriptions, ch)
}

// waitForNextConfig waits for the next received Config update from the source and then returns it.
func (p *ConfigDistributor) waitForNextConfig() (*Config, error) {
	nextConfigChan := make(chan *Config)
	defer close(nextConfigChan)

	// Subscribe to ourselves for the duration of this method
	err := p.Subscribe(nextConfigChan)

	if err != nil {
		return nil, err
	}

	defer p.Unsubscribe(nextConfigChan)

	// Wait for the config or Cancel signal
	select {
	case config := <-nextConfigChan:
		return config, nil

	case <-p.closed.Done():
		return nil, ErrConfigDistributorClosed
	}
}

// monitorConfig waits for Config updates from the source and notifies subscribers.
func (p *ConfigDistributor) monitorConfig() {
	// For debugging purposes
	log.Debug("config_distributor: starting monitorConfig() loop")
	defer log.Debug("config_distributor: monitorConfig() loop exited")

	for {
		select {
		// Wait for next config
		case newConfig, ok := <-p.source:
			if !ok {
				log.Debug("config_distributor: source channel is closed, closing the distributor")
				p.Close()

				return
			}

			// Update lastConfig and notify subscribers
			p.mu.Lock()
			{
				oldConfig := p.lastConfig
				p.lastConfig = newConfig

				p.notifySubscribersLocked(oldConfig, newConfig)
			}
			p.mu.Unlock()

		// Stop the loop if distributor is closed
		case <-p.closed.Done():
			return
		}
	}
}

// notifySubscribersLocked notifies subscribers about config changes.
func (p *ConfigDistributor) notifySubscribersLocked(oldConfig, newConfig *Config) {
	for ch, checks := range p.subscriptions {
		if p.shouldSendUpdatedConfig(checks, oldConfig, newConfig) {
			ch <- newConfig
		}
	}
}

// shouldSendUpdatedConfig runs the subscription checks to determine whether or not
// subscribed channels should receive update.
func (p *ConfigDistributor) shouldSendUpdatedConfig(checks []ConfigChangedCheckFn, oldConfig, newConfig *Config) bool {
	if len(checks) == 0 {
		return true
	}

	// Search for any check that returns true
	for _, fn := range checks {
		if fn(oldConfig, newConfig) {
			return true
		}
	}

	return false
}
