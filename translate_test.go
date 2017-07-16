package translate

import (
	"reflect"
	"testing"

	"strings"

	"github.com/arteev/go-translate/translator"
)

//TODO: test opts

type fakeprovider struct {
	getlang   bool
	detect    bool
	translate bool
}

//Get support languages
func (p *fakeprovider) GetLangs(code string) ([]*translator.Language, error) {
	p.getlang = true
	return nil, nil
}

func (p *fakeprovider) Detect(text string) (*translator.Language, error) {
	p.detect = true
	return nil, nil
}

func (p *fakeprovider) Translate(text, direction string) *translator.Result {
	p.translate = true
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
func TestProv(t *testing.T) {
	unregisterAllTranslators()
	tf := &testfirst{}
	Register("test", tf)

	p, ok := translators["test"]
	_, is := p.(*testfirst)
	if !ok || !is {
		t.Errorf("Expected %v, got %v", tf, p)
	}

	if len(translators) != 1 {
		t.Errorf("Expected count providers %d,got %d", 1, len(translators))
	}
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(tr.translator).String() != "*translate.fakeprovider" {
		t.Errorf("Expected provider type fakeprovider, got %v", reflect.TypeOf(tr.translator).String())
	}

}

func TestNotExistsProv(t *testing.T) {

	unregisterAllTranslators()
	_, err := New("notexists")
	if err == nil || !strings.Contains(err.Error(), "translator: unknown translator") {
		t.Errorf("Expected error:translator: unknown translator..., got %s", err)
	}
}

func TestCalled(t *testing.T) {
	unregisterAllTranslators()
	Register("test", &testfirst{})
	tr, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	tr.GetLangs("en")
	fp := tr.translator.(*fakeprovider)
	if !fp.getlang {
		t.Error("Expected call GetLangs")
	}
	tr.Detect("")
	if !fp.detect {
		t.Error("Expected call Detect")
	}
	tr.Translate("", "")
	if !fp.translate {
		t.Error("Expected call Translate ")
	}

}
