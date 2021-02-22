package language

// Language - structure defining the language (ISO 639-2)
//of the translation and the direction of possible translations
type Language struct {
	Code string
	Name string
}

// New - returns *Language of language with name and code
func New(code, name string) Language {
	return Language{
		Code: code,
		Name: name,
	}
}

// String - Stringer
func (l Language) String() string {
	return l.Code
}

// Equal returns true if language codes match
func (l Language) Equal(language Language) bool {
	return l.Code == language.Code
}

func (l Language) Empty() bool {
	return l.Code == ""
}