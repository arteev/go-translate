package translator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTargets_Append(t *testing.T) {
	targets := make(Targets)
	ru := New("ru", "Russian")
	en := New("en", "English")
	fr := New("fr", "French")
	targets.Append(ru, en, fr)
	targets.Append(ru, fr)
	targets.Append(ru, en)

	src, ok := targets[ru]
	assert.True(t, ok)
	assert.EqualValues(t, []Language{en, fr}, src)
}
