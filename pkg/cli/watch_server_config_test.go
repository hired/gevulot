package cli

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hired/gevulot/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestWatchServerConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "config")

	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString("listen = '0.0.0.0:4242'\n"); err != nil {
		t.Fatal(err)
	}

	configChan := make(chan *server.Config, 1)

	err = watchServerConfig(tmpfile.Name(), configChan)

	assert.NoError(t, err)
	assert.Empty(t, configChan)

	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.WriteString("listen = '0.0.0.0:31337'\n"); err != nil {
		t.Fatal(err)
	}

	assert.Eventually(t, func() bool { return len(configChan) == 1 }, time.Second, 10*time.Millisecond)

	// Previous assertion checked that the channel has at least one message so this will never block
	assert.Equal(t, "0.0.0.0:31337", (<-configChan).Listen)
}
