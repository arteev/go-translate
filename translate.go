package translate

import (
	"errors"
	"fmt"
	"sort"

	"sync"

	"github.com/arteev/go-translate/translator"
)

//Errors
var (
	ErrUnknowProvider = errors.New("Unknown provider")
)

//TranslatorFactory factory translator
type TranslatorFactory interface {
	NewInstance(opts map[string]interface{}) translator.Translator
}

var (
	muTranslators sync.RWMutex
	translators   = make(map[string]TranslatorFactory)
)

type Translate struct {
	translator translator.Translator
	namesLangs map[string][]*translator.Language
	opts       map[string]interface{}
}

func Register(name string, factory TranslatorFactory) {
	muTranslators.Lock()
	defer muTranslators.Unlock()
	if _, dup := translators[name]; dup {
		panic("translator: Register called twice for driver " + name)
	}
	translators[name] = factory
}

// For tests.
func unregisterAllTranslators() {
	muTranslators.Lock()
	defer muTranslators.Unlock()
	translators = make(map[string]TranslatorFactory)
}

// Translators returns a sorted list of the names of the registered translators.
func Translators() []string {
	muTranslators.RLock()
	defer muTranslators.RUnlock()
	var list []string
	for name := range translators {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

type Option func(*Translate)

func WithOption(name string, value interface{}) Option {
	return func(t *Translate) {
		t.opts[name] = value
	}
}

func New(name string, opts ...Option) (*Translate, error) {
	muTranslators.RLock()
	trs, ok := translators[name]
	muTranslators.RUnlock()
	if !ok {
		return nil, fmt.Errorf("translator: unknown translator %q", name)
	}

	tr := &Translate{
		namesLangs: make(map[string][]*translator.Language),
		opts:       make(map[string]interface{}),
	}
	//Fill options
	for _, opt := range opts {
		opt(tr)
	}

	tr.translator = trs.NewInstance(tr.getOptions())
	return tr, nil
}

func (t *Translate) getOptions() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range t.opts {
		m[k] = v
	}
	return m
}

func (t *Translate) GetLangs(code string) ([]*translator.Language, error) {
	if langs, ok := t.namesLangs[code]; ok {
		return langs, nil
	}
	langs, err := t.translator.GetLangs(code)
	if err != nil {
		return nil, err
	}
	if code != "" {
		t.namesLangs[code] = langs
	}
	return langs, nil
}

func (t *Translate) Detect(text string) (*translator.Language, error) {
	l, err := t.translator.Detect(text)
	if err != nil {
		return nil, err
	}
	return l, err
}

func (t *Translate) Translate(text, direction string) *translator.Result {
	return t.translator.Translate(text, direction)
}
