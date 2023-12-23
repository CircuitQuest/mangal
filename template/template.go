package template

import (
	"log"
	"strings"
	"text/template"

	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/template/util"
	"github.com/mangalorg/libmangal"
)

func Chapter(_ string, chapter libmangal.Chapter) string {
	var sb strings.Builder

	err := template.Must(template.New("chapter").
		Funcs(util.FuncMap).
		Parse(config.Config.Download.Chapter.NameTemplate.Get())).
		Execute(&sb, chapter.Info())

	if err != nil {
		log.Fatal("error during execution of the chapter name template", "err", err)
	}

	return sb.String()
}

func Manga(_ string, manga libmangal.Manga) string {
	var sb strings.Builder

	err := template.Must(template.New("manga").
		Funcs(util.FuncMap).
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
		Funcs(util.FuncMap).
		Parse(config.Config.Download.Volume.NameTemplate.Get())).
		Execute(&sb, manga.Info())

	if err != nil {
		log.Fatal("error during execution of the volume name template", "err", err)
	}

	return sb.String()
}
