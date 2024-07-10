package mangas

import "github.com/luevano/libmangal/mangadata"

type searchMetadataMsg struct {
	item *item
}

type searchVolumesMsg struct {
	item *item
}

type searchAllChaptersMsg struct {
	manga   mangadata.Manga
	volumes []mangadata.Volume
}

type searchChaptersMsg struct {
	manga  mangadata.Manga
	volume mangadata.Volume
}
