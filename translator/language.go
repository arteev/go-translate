package translator

// Language - structure defining the language
//of the translation and the direction of possible translations
type Language struct {
	Code string
	Name string
	Dirs []*Language
}

// NewLanguage - returns *Language of language with name and code
func NewLanguage(code, name string) *Language {
	return &Language{
		Code: code,
		Name: name,
	}
}

// AddDir - Adds translation direction
func (l *Language) AddDir(d *Language) {
	if l == d || l.Code == d.Code {
		return
	}
	for _, ld := range l.Dirs {
		if ld == d {
			return
		}
	}
	l.Dirs = append(l.Dirs, d)
}

// String - Stringer
func (l Language) String() string {
	return l.Code
}
