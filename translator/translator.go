package translator

import "errors"

var (
	ErrWrongApiKey          = errors.New("Wrong API key")
	ErrBlockedApiKey        = errors.New("The API key is blocked")
	ErrUnsupported          = errors.New("Unsupported")
	ErrLimitDayExceeded     = errors.New("Day limit exceeded")
	ErrLimitMonthExceeded   = errors.New("Month limit exceeded")
	ErrLimitTextExceeded    = errors.New("Exceeded the maximum size of the text")
	ErrTextNotTranslated    = errors.New("Text can not be translated")
	ErrDirectionUnsupported = errors.New("Set the direction of translation is not supported")
)

type Translator interface {
	//Get support languages
	GetLangs(code string) ([]*Language, error)
	//Detect language
	Detect(text string) (*Language, error)
	//Translate text
	Translate(text, direction string) *Result
	//Name of translator
	Name() string
}

type Result struct {
	Text     string
	FromLang *Language
	ToLang   *Language
	Detected *Language
	Err      error
}
