package translator

import (
	"fmt"
	"testing"
)

func TestLang(t *testing.T) {
	const (
		code = "ru"
		name = "Russian"
	)

	lang := NewLanguage(code, name)

	if lang.Code != code {
		t.Errorf("Excepted %s, got %s", code, lang.Code)
	}
	if lang.Name != name {
		t.Errorf("Excepted %s, got %s", name, lang.Name)
	}

	if fmt.Sprint(lang) != code {
		t.Errorf("Excepted %s, got %s", code, fmt.Sprint(lang))
	}

	if len(lang.Dirs) != 0 {
		t.Errorf("Excepted len dir %d, got %d", 0, len(lang.Dirs))
	}
	lang.AddDir(lang)
	if len(lang.Dirs) != 0 {
		t.Errorf("Excepted len dir %d, got %d", 0, len(lang.Dirs))
	}

	en := NewLanguage("en", "English")
	lang.AddDir(en)
	if len(lang.Dirs) != 1 {
		t.Errorf("Excepted len dir %d, got %d", 1, len(lang.Dirs))
	}

	if lang.Dirs[0] != en {
		t.Errorf("Excepted %v, got %v", en, lang.Dirs[0])
	}
	lang.AddDir(en)
	if len(lang.Dirs) != 1 {
		t.Errorf("Excepted len dir %d, got %d", 1, len(lang.Dirs))
	}

}
