package translate

import (
	"errors"
	"fmt"
	"github.com/arteev/go-translate/language"
	"sort"

	"sync"
)

//Errors
var (
	ErrUnknownProvider      = errors.New("unknown provider")
	ErrWrongAPIKey          = errors.New("wrong API key")
	ErrBlockedAPIKey        = errors.New("the API key is blocked")
	ErrUnsupported          = errors.New("unsupported")
	ErrLimitDayExceeded     = errors.New("day limit exceeded")
	ErrLimitMonthExceeded   = errors.New("month limit exceeded")
	ErrLimitTextExceeded    = errors.New("exceeded the maximum size of the text")
	ErrTextNotTranslated    = errors.New("text can not be translated")
	ErrDirectionUnsupported = errors.New("set the direction of translation is not supported")
)

//Translator -this interface defines the basic
//translation methods for specific translation providers
type Translator interface {
	//Get support languages
	GetLanguages(code string) ([]language.Language, error)
	//Detect language
	Detect(text string) (language.Language, error)
	//Translate text
	Translate(text, direction string) (*Result, error)
	//Name of translator
	Name() string
}

//A Result of translation
type Result struct {
	Text     string
	FromLang *language.Language
	ToLang   *language.Language
	Detected *language.Language
}

//TranslatorFactory - factory translator
type TranslatorFactory interface {
	NewInstance(opts map[string]interface{}) Translator
}

var (
	muTranslators sync.RWMutex
	translators   = make(map[string]TranslatorFactory)
)

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

type Options map[string]interface{}

//Option for storing translator options
type Option func(*Options)

//WithOption - Adds an optional parameter for the translator at the time of creation
func WithOption(name string, value interface{}) Option {
	return func(t *Options) {
		(*t)[name] = value
	}
}

//New - Creates an translator with a name and opts options
func New(name string, opts ...Option) (Translator, error) {
	muTranslators.RLock()
	trs, ok := translators[name]
	muTranslators.RUnlock()
	if !ok {
		return nil, fmt.Errorf("translator: unknown translator %q", name)
	}

	options := &Options{}
	//Fill options
	for _, opt := range opts {
		opt(options)
	}

	translator := trs.NewInstance(*options)
	return translator, nil
}
