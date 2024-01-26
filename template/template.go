package template

import (
	"log"
	"strings"
	"text/template"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/template/funcs"
)

// TODO: change logger?

func Provider(provider libmangal.ProviderInfo) string {
	var sb strings.Builder

	err := template.Must(template.New("provider").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Provider.NameTemplate.Get())).
		Execute(&sb, provider)
	if err != nil {
		log.Fatal("error during execution of the provider name template", "err", err)
	}

	return sb.String()
}

func Manga(_ string, manga libmangal.Manga) string {
	// TODO: will change to need MangaWithAnilist in the future, to have standardized manga dir names
	// _, ok := manga.(libmangal.MangaWithAnilist)
	// if !ok {
	// 	panic("[Manga template]manga doesn't have anilist data")
	// }
	var sb strings.Builder

	err := template.Must(template.New("manga").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Manga.NameTemplate.Get())).
		Execute(&sb, manga.Info())
	if err != nil {
		log.Fatal("error during execution of the manga name template", "err", err)
	}

	return sb.String()
}

func Volume(_ string, manga libmangal.Volume) string {
	var sb strings.Builder

	err := template.Must(template.New("volume").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Volume.NameTemplate.Get())).
		Execute(&sb, manga.Info())
	if err != nil {
		log.Fatal("error during execution of the volume name template", "err", err)
	}

	return sb.String()
}

func Chapter(_ string, chapter libmangal.Chapter) string {
	var sb strings.Builder

	err := template.Must(template.New("chapter").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Chapter.NameTemplate.Get())).
		Execute(&sb, chapter.Info())
	if err != nil {
		log.Fatal("error during execution of the chapter name template", "err", err)
	}

	return sb.String()
}

// Currently unused.
func Config(field config.Entry) string {
	var sb strings.Builder

	err := template.Must(template.New("field").
		Funcs(funcs.FuncMap).
		Parse(`
{{ .Description }}

Key: {{ .Key }}
Value: {{ getConfig .Key }}
Default: {{ .Default }}
		`)).Execute(&sb, field)
	if err != nil {
		log.Fatal("error during execution of the config name template", "err", err)
	}

	return sb.String()
}
