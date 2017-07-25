package translate

import (
	"reflect"
	"testing"

	"strings"

	"errors"

	"github.com/arteev/go-translate/translator"
)

type fakeprovider struct {
	invokeGetlang   bool
	invokeDetect    bool
	invokeTranslate bool
}

//Get support languages
func (p *fakeprovider) GetLangs(code string) ([]*translator.Language, error) {
	p.invokeGetlang = true
	if code == "en" {
		return []*translator.Language{translator.NewLanguage("en", "English")}, nil
	}
	return nil, errors.New("Unsupported")
}

func (p *fakeprovider) Detect(text string) (*translator.Language, error) {
	p.invokeDetect = true
	if text == "" {
		return nil, errors.New("Text is empty")
	}
	return translator.NewLanguage("en", "English"), nil
}

func (p *fakeprovider) Translate(text, direction string) *translator.Result {
	p.invokeTranslate = true
	return nil
}

func (fakeprovider) Name() string {
	return ""
}

type testfirst struct{}

func (testfirst) NewInstance(opts map[string]interface{}) translator.Translator {
	return &fakeprovider{}
}

func testsecond(opts map[string]interface{}) translator.Translator {
	return nil
}

func TestDupTranslator(t *testing.T) {
	defer func() {
		dupChecked := false
		if e := recover(); e != nil {
			if str, ok := e.(string); ok && strings.Contains(str, "translator: Register called twice for driver") {
				dupChecked = true
			}
		}
		if !dupChecked {
			t.Error("Expected panic: translator: Register called twice for driver test")
		}
	}()
	unregisterAllTranslators()
	Register("test", &testfirst{})
	Register("test", &testfirst{})
	t.Error("Expected panic: translator: Register called twice for driver test")
}
func TestRegTranslators(t *testing.T) {
	unregisterAllTranslators()
	tf := &testfirst{}
	Register("test", tf)

	p, ok := translators["test"]
	_, is := p.(*testfirst)
	if !ok || !is {
		t.Errorf("Expected %v, got %v", tf, p)
	}

	if len(translators) != 1 {
		t.Errorf("Expected count translators %d,got %d", 1, len(translators))
	}
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(tr.translator).String() != "*translate.fakeprovider" {
		t.Errorf("Expected provider type fakeprovider, got %v", reflect.TypeOf(tr.translator).String())
	}

	if arrTrs := Translators(); len(arrTrs) != 1 {
		t.Errorf("Expected count Translators() %d,got %d", 1, len(arrTrs))
	} else {
		if arrTrs[0] != "test" {
			t.Errorf("Expected count Translators[%d]=%s,got %s", 1, "test", arrTrs[0])
		}
	}

	unregisterAllTranslators()
	if len(translators) != 0 {
		t.Errorf("Expected count providers %d,got %d", 0, len(translators))
	}

	if arrTrs := Translators(); len(arrTrs) != 0 {
		t.Errorf("Expected count Translators() %d,got %d", 0, len(arrTrs))
	}

}

func TestNotExistsProv(t *testing.T) {
	unregisterAllTranslators()
	_, err := New("notexists")
	if err == nil || !strings.Contains(err.Error(), "translator: unknown translator") {
		t.Errorf("Expected error:translator: unknown translator..., got %s", err)
	}
}

func TestInvoke(t *testing.T) {
	unregisterAllTranslators()
	Register("test", &testfirst{})
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	tr.GetLangs("en")
	fp := tr.translator.(*fakeprovider)
	if !fp.invokeGetlang {
		t.Error("Expected call GetLangs")
	}
	_, err = tr.Detect("text")
	if err != nil {
		t.Error(err)
	}
	if !fp.invokeDetect {
		t.Error("Expected call Detect")
	}

	tr.Translate("", "")
	if !fp.invokeTranslate {
		t.Error("Expected call Translate ")
	}
}

func TestDetect(t *testing.T) {
	unregisterAllTranslators()
	Register("test", &testfirst{})
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tr.Detect("")
	if err == nil || err.Error() != "Text is empty" {
		t.Errorf("Expected error: %q, got %s", "Text is empty", err)
	}

	if l, err := tr.Detect("Hello"); err != nil {
		t.Error(err)
	} else {
		if l == nil {
			t.Fatal("Expected Detect() not nil")
		}
		if l.Code != "en" {
			t.Errorf("Expected %s,got %s", translator.NewLanguage("en", "English"), l)
		}
	}

}

func TestGetLang(t *testing.T) {
	unregisterAllTranslators()
	Register("test", &testfirst{})
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	lgs, err := tr.GetLangs("en")
	if err != nil {
		t.Fatal(err)
	}
	if len(lgs) != 1 {
		t.Errorf("Expected count GetLangs %d,got %d", 1, len(lgs))
	}

	_, err = tr.GetLangs("")
	if err == nil {
		t.Fatal("Expected error for GetLangs, got nil")
	}
	if err.Error() != "Unsupported" {
		t.Errorf("Expected error GetLangs %q,got %s", "Unsupported", err)
	}

	lgs2, err := tr.GetLangs("en")
	if err != nil {
		t.Fatal(err)
	}

	if lgs[0] != lgs2[0] {
		t.Errorf("Expected GetLangs() returns ptr[0] %p,got %p", lgs[0], lgs2[0])
	}
}

func TestOptions(t *testing.T) {
	unregisterAllTranslators()
	Register("test", &testfirst{})

	tr, err := New("test",
		WithOption("apikey", "123"),
		WithOption("id", "1"))
	if err != nil {
		t.Fatal(err)
	}

	opts := tr.getOptions()
	if opts == nil {
		t.Fatal("Expected map with options")
	}

	if len(opts) != 2 {
		t.Errorf("Expected len(options) %d, got %d", 2, len(opts))
	}

	if got, ok := opts["apikey"]; !ok || got != "123" {
		t.Errorf("Expected for key:apikey value:%s, got: %s", "123", got)
	}
}
