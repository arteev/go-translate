package translator

import "errors"

//Errors
var (
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
	GetLangs(code string) ([]Language, error)
	//Detect language
	Detect(text string) (Language, error)
	//Translate text
	Translate(text, direction string) *Result
	//Name of translator
	Name() string
}

//A Result of translation
type Result struct {
	Text     string
	FromLang *Language
	ToLang   *Language
	Detected *Language
	Err      error
}
