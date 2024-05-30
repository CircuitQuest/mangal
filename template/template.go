package template

import (
	"strings"
	"text/template"

	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/template/funcs"
	"github.com/luevano/mangal/util"
)

func Provider(provider libmangal.ProviderInfo) string {
	var sb strings.Builder

	err := template.Must(template.New("provider").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Provider.NameTemplate.Get())).
		Execute(&sb, provider)
	if err != nil {
		util.Errorf("error during execution of the provider name template: %s\n", err)
	}

	return sb.String()
}

func Manga(_ string, manga mangadata.Manga) string {
	var sb strings.Builder

	plt := config.Config.Download.Manga.NameTemplateFallback.Get()
	// Prioritize the NameTemplate (includes AnilistManga data)
	metadata := manga.Metadata()
	if metadata != nil {
		plt = config.Config.Download.Manga.NameTemplate.Get()
	}

	err := template.Must(template.New("manga").
		Funcs(funcs.FuncMap).
		Parse(plt)).
		Execute(&sb, manga)
	if err != nil {
		util.Errorf("error during execution of the manga name template: %s\n", err)
	}

	return sb.String()
}

func Volume(_ string, manga mangadata.Volume) string {
	var sb strings.Builder

	err := template.Must(template.New("volume").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Volume.NameTemplate.Get())).
		Execute(&sb, manga.Info())
	if err != nil {
		util.Errorf("error during execution of the volume name template: %s\n", err)
	}

	return sb.String()
}

func Chapter(_ string, chapter mangadata.Chapter) string {
	var sb strings.Builder

	err := template.Must(template.New("chapter").
		Funcs(funcs.FuncMap).
		Parse(config.Config.Download.Chapter.NameTemplate.Get())).
		Execute(&sb, chapter.Info())
	if err != nil {
		util.Errorf("Error during execution of the chapter name template: %s\n", err)
	}

	return sb.String()
}
