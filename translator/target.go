package translator

// Targets map Source -> [Target1,Target2]
type Targets map[Language][]Language

func exists(l Language, items []Language) bool {
	for _, item := range items {
		if item.Equal(l) {
			return true
		}
	}
	return false
}

func (t *Targets) Append(source Language, languages ...Language) {
	listSource, ok := (*t)[source]
	if !ok {
		listSource = make([]Language, 0)
		(*t)[source] = listSource
	}
	for _, language := range languages {
		if !exists(language, listSource) {
			listSource = append(listSource, language)
		}
	}
	(*t)[source] = listSource
}
