package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type distributorGetConfigResult struct {
	config *Config
	err    error
}

func TestConfigDistributorGet(t *testing.T) {
	configV1 := &Config{Listen: "0.0.0.0:4242"}
	configV2 := &Config{Listen: "0.0.0.0:4242", DatabaseURL: "postgresql://"}

	src := make(chan *Config)
	defer close(src)

	p := NewConfigDistributor(src)
	defer p.Close()

	resultCh := make(chan *distributorGetConfigResult)
	defer close(resultCh)

	callGet := func() {
		go func() {
			config, err := p.Get()
			resultCh <- &distributorGetConfigResult{config, err}
		}()
	}

	// Get() waits for initial config indefinitely
	callGet()

	fmt.Printf("waits for config\n")
	select {
	case <-resultCh:
		assert.FailNow(t, "Get() returned before config became available")
	default:
	}

	// Initial config
	src <- configV1
	assert.Same(t, configV1, (<-resultCh).config)

	// Updated config
	src <- configV2

	callGet()
	assert.Same(t, configV2, (<-resultCh).config)

	// Close distributor
	p.Close()

	callGet()
	assert.Equal(t, ErrConfigDistributorClosed, (<-resultCh).err)
}

func TestConfigDistributorSubscribe(t *testing.T) {
	configV1 := &Config{Listen: "0.0.0.0:4242"}
	configV2 := &Config{Listen: "0.0.0.0:4242", DatabaseURL: "postgresql://"}
	configV3 := &Config{Listen: "0.0.0.0:31337", DatabaseURL: "postgresql://"}

	src := make(chan *Config, 1)
	defer close(src)

	// Initial config
	src <- configV1

	allChanges := make(chan *Config, 1)
	defer close(allChanges)

	specificChanges := make(chan *Config, 1)
	defer close(specificChanges)

	p := NewConfigDistributor(src)
	defer p.Close()

	err := p.Subscribe(specificChanges, func(oldConfig, newConfig *Config) bool {
		return oldConfig == nil || oldConfig.Listen != newConfig.Listen
	})
	assert.NoError(t, err)

	err = p.Subscribe(allChanges)
	assert.NoError(t, err)

	// Initial config is delivered
	assert.Equal(t, configV1, <-allChanges)
	assert.Equal(t, configV1, <-specificChanges)

	// Updated config is delivered to only 1 channel
	src <- configV2

	assert.Equal(t, configV2, <-allChanges)

	// Updated config is delivered to both channels
	src <- configV3

	assert.Equal(t, configV3, <-allChanges)
	assert.Equal(t, configV3, <-specificChanges) // this will fail if configV2 was mistakenly delivered

	// Subscription is unavailable after distributor is closed
	p.Close()

	err = p.Subscribe(allChanges) // TODO: should it be an error that we are re-subscribing channel?
	assert.Equal(t, ErrConfigDistributorClosed, err)
}

func TestConfigDistributorUnsubscribe(t *testing.T) {
	configV1 := &Config{Listen: "0.0.0.0:4242"}
	configV2 := &Config{Listen: "0.0.0.0:4242", DatabaseURL: "postgresql://"}

	src := make(chan *Config, 1)
	defer close(src)

	// Initial config
	src <- configV1

	subscriptionCh := make(chan *Config)
	defer close(subscriptionCh)

	p := NewConfigDistributor(src)
	defer p.Close()

	err := p.Subscribe(subscriptionCh)
	assert.NoError(t, err)

	// Initial config is delivered
	assert.Equal(t, configV1, <-subscriptionCh)

	// New changes are not delivered once channel is unsubsribed
	p.Unsubscribe(subscriptionCh)

	// Updated config
	src <- configV2

	select {
	case <-subscriptionCh:
		assert.FailNow(t, "Unsubscribe() failed to unsubscribe the channel")
	default:
	}
}
