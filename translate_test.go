package translate

import (
	"github.com/arteev/go-translate/language"
	"reflect"
	"testing"

	"strings"

	"errors"
)

type fakeprovider struct {
	invokeGetlang   bool
	invokeDetect    bool
	invokeTranslate bool
	opts            map[string]interface{}
}

//Get support languages
func (p *fakeprovider) GetLanguages(code string) ([]language.Language, error) {
	p.invokeGetlang = true
	if code == "en" {
		return []language.Language{language.New("en", "English")}, nil
	}
	return nil, errors.New("Unsupported")
}

func (p *fakeprovider) Detect(text string) (language.Language, error) {
	p.invokeDetect = true
	if text == "" {
		return language.Language{}, errors.New("Text is empty")
	}
	return language.New("en", "English"), nil
}

func (p *fakeprovider) Translate(text, direction string) (*Result, error) {
	p.invokeTranslate = true
	return nil, nil
}

func (fakeprovider) Name() string {
	return "fakeprovider"
}

type testfirst struct{}

func (testfirst) NewInstance(opts map[string]interface{}) Translator {
	return &fakeprovider{
		opts: opts,
	}
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
	if reflect.TypeOf(tr).String() != "*translate.fakeprovider" {
		t.Errorf("Expected provider type fakeprovider, got %v", reflect.TypeOf(tr).String())
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
	_, err = tr.GetLanguages("en")
	if err != nil {
		t.Fatal(err)
	}
	fp := tr.(*fakeprovider)
	if !fp.invokeGetlang {
		t.Error("Expected call GetLanguages")
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

func TestName(t *testing.T) {
	unregisterAllTranslators()
	Register("test", &testfirst{})
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	if name := tr.Name(); name != "fakeprovider" {
		t.Errorf("Expected Name fakeprovider, got %s", name)
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
		if l.Empty() {
			t.Fatal("Expected Detect() not nil")
		}
		if l.Code != "en" {
			t.Errorf("Expected %s,got %s", language.New("en", "English"), l)
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
	lgs, err := tr.GetLanguages("en")
	if err != nil {
		t.Fatal(err)
	}
	if len(lgs) != 1 {
		t.Errorf("Expected count GetLanguages %d,got %d", 1, len(lgs))
	}

	_, err = tr.GetLanguages("")
	if err == nil {
		t.Fatal("Expected error for GetLanguages, got nil")
	}
	if err.Error() != "Unsupported" {
		t.Errorf("Expected error GetLanguages %q,got %s", "Unsupported", err)
	}

	lgs2, err := tr.GetLanguages("en")
	if err != nil {
		t.Fatal(err)
	}

	if lgs[0] != lgs2[0] {
		t.Errorf("Expected GetLanguages() returns ptr[0] %v,got %v", lgs[0], lgs2[0])
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

	opts := tr.(*fakeprovider).opts
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
