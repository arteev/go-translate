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
	ErrUnknownProvider = errors.New("unknown provider")
)

//TranslatorFactory - factory translator
type TranslatorFactory interface {
	NewInstance(opts map[string]interface{}) translator.Translator
}

var (
	muTranslators sync.RWMutex
	translators   = make(map[string]TranslatorFactory)
)

//Translate - wrappers for Translator interface
type Translate struct {
	translator    translator.Translator
	nameLanguages map[string][]*translator.Language
	opts          map[string]interface{}
}

var _ translator.Translator = (*Translate)(nil)

//Register - registers translator with name and factory function
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

//Option for storing translator options
type Option func(*Translate)

//WithOption - Adds an optional parameter for the translator at the time of creation
func WithOption(name string, value interface{}) Option {
	return func(t *Translate) {
		t.opts[name] = value
	}
}

//New - Creates an translator with a name and opts options
func New(name string, opts ...Option) (*Translate, error) {
	muTranslators.RLock()
	trs, ok := translators[name]
	muTranslators.RUnlock()
	if !ok {
		return nil, fmt.Errorf("translator: unknown translator %q", name)
	}

	tr := &Translate{
		nameLanguages: make(map[string][]*translator.Language),
		opts:          make(map[string]interface{}),
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

//GetLangs - returns supported languages
func (t *Translate) GetLangs(code string) ([]*translator.Language, error) {
	if langs, ok := t.nameLanguages[code]; ok {
		return langs, nil
	}
	langs, err := t.translator.GetLangs(code)
	if err != nil {
		return nil, err
	}
	if code != "" {
		t.nameLanguages[code] = langs
	}
	return langs, nil
}

//Detect - returns automatically detected text language
func (t *Translate) Detect(text string) (*translator.Language, error) {
	l, err := t.translator.Detect(text)
	if err != nil {
		return nil, err
	}
	return l, err
}

//Translate - returns the translated text to the language direction
func (t *Translate) Translate(text, direction string) *translator.Result {
	return t.translator.Translate(text, direction)
}

func (t *Translate) Name() string {
	return t.translator.Name()
}
