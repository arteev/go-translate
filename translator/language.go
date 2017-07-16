package translator

type Language struct {
	Code string
	Name string
	Dirs []*Language
}

func NewLanguage(code, name string) *Language {
	return &Language{
		Code: code,
		Name: name,
	}
}
func (l *Language) AddDir(d *Language) {
	if l == d || l.Code == d.Code {
		return
	}
	for i := 0; i < len(l.Dirs); i++ {
		if l.Dirs[i] == d {
			return
		}
	}
	l.Dirs = append(l.Dirs, d)
}

func (l Language) String() string {
	return l.Code
}
