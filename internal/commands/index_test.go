package commands

import (
	"testing"
	"time"

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/pkg/capture"
)

func TestIndexCommand(t *testing.T) {
	var err error

	ctx := config.CliTestContext()

	s := event.Subscribe("log.info")
	defer event.Unsubscribe(s)

	var l string

	assert.IsType(t, hub.Subscription{}, s)

	go func() {
		for msg := range s.Receiver {
			l += msg.Fields["message"].(string) + "\n"
		}
	}()

	stdout := capture.Output(func() {
		err = IndexCommand.Run(ctx)
	})

	if err != nil {
		t.Fatal(err)
	}

	if stdout != "" {
		t.Logf("stdout: %s", stdout)
	}

	time.Sleep(time.Second)

	if l != "" {
		// Expected index command output.
		assert.Contains(t, l, "classify: loading labels")
		assert.NotContains(t, l, "error")
		assert.Contains(t, l, "found no .ppignore file")
	} else {
		t.Fatal("log output missing")
	}
}
