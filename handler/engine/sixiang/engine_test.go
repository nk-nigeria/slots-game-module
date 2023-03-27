package sixiang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlotsEngine(t *testing.T) {
	name := "TestNewSlotsEngine"
	t.Run(name, func(t *testing.T) {
		got := NewEngine()
		assert.NotNil(t, got)
	})
}
