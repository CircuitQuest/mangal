package template

import (
	"strings"
	"text/template"

	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/template/funcs"
	"github.com/luevano/mangal/util"
)

type mangaTemplateData struct {
	Manga    mangadata.MangaInfo
	Metadata metadata.Metadata
}

type volumeTemplateData struct {
	Volume mangadata.VolumeInfo
	Manga  mangadata.MangaInfo
}

type chapterTemplateData struct {
	Chapter mangadata.ChapterInfo
	Volume  mangadata.VolumeInfo
	Manga   mangadata.MangaInfo
}

func Provider(provider libmangal.ProviderInfo) string {
	var sb strings.Builder

	err := template.Must(template.New("provider").
		Funcs(funcs.FuncMap).
		Parse(config.Download.Provider.NameTemplate.Get())).
		Execute(&sb, provider)
	if err != nil {
		util.Errorf("error during execution of the provider name template: %s\n", err)
	}

	return sb.String()
}

func Manga(_ string, manga mangadata.Manga) string {
	var sb strings.Builder

	plt := config.Download.Manga.NameTemplateFallback.Get()
	// Prioritize the NameTemplate (includes AnilistManga data)
	meta := manga.Metadata()
	if metadata.Validate(meta) == nil {
		plt = config.Download.Manga.NameTemplate.Get()
	}

	err := template.Must(template.New("manga").
		Funcs(funcs.FuncMap).
		Parse(plt)).
		Execute(&sb, mangaTemplateData{
			Manga:    manga.Info(),
			Metadata: meta,
		})
	if err != nil {
		util.Errorf("error during execution of the manga name template: %s\n", err)
	}

	return sb.String()
}

func Volume(_ string, volume mangadata.Volume) string {
	var sb strings.Builder

	err := template.Must(template.New("volume").
		Funcs(funcs.FuncMap).
		Parse(config.Download.Volume.NameTemplate.Get())).
		Execute(&sb, volumeTemplateData{
			Volume: volume.Info(),
			Manga:  volume.Manga().Info(),
		})
	if err != nil {
		util.Errorf("error during execution of the volume name template: %s\n", err)
	}

	return sb.String()
}

func Chapter(_ string, chapter mangadata.Chapter) string {
	var sb strings.Builder

	err := template.Must(template.New("chapter").
		Funcs(funcs.FuncMap).
		Parse(config.Download.Chapter.NameTemplate.Get())).
		Execute(&sb, chapterTemplateData{
			Chapter: chapter.Info(),
			Volume:  chapter.Volume().Info(),
			Manga:   chapter.Volume().Manga().Info(),
		})
	if err != nil {
		util.Errorf("Error during execution of the chapter name template: %s\n", err)
	}

	return sb.String()
}
