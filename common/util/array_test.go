package util

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestUnique(t *testing.T) {
	slices := []string{"123", "123", "456"}
	slices = Unique(slices)
	assert.Equal(t, len(slices), 2)
}
