package mangas

import "github.com/luevano/libmangal/mangadata"

type Item struct {
	manga *mangadata.Manga
}

func (i Item) FilterValue() string {
	return (*i.manga).String()
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return (*i.manga).Info().URL
}
