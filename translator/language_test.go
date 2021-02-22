package translator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLanguage(t *testing.T) {
	const (
		code = "ru"
		name = "Russian"
	)

	lang := New(code, name)

	assert.Equal(t, code, lang.Code)
	assert.Equal(t, name, lang.Name)
	assert.Equal(t, code, lang.String())

	lang2 := New("en", "English")
	assert.False(t, lang2.Equal(lang))
	assert.False(t, lang.Equal(lang2))

	langEqual := New("ru", "ru")
	assert.True(t, langEqual.Equal(lang))
	assert.True(t, lang.Equal(langEqual))
}

func TestLanguage_Empty(t *testing.T) {
	assert.True(t, Language{}.Empty())
	assert.False(t, New("en","en").Empty())
}
