package config

import (
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/metadata"
)

// DownloadOptions constructs the libmangal.DownloadOptions populated by the Config.
func DownloadOptions() libmangal.DownloadOptions {
	// start from the defaults in case of new additions and build on top of it.
	o := libmangal.DefaultDownloadOptions()
	o.Format = Download.Format.Get()
	o.Directory = Download.Path.Get()
	o.CreateProviderDir = Download.Provider.CreateDir.Get()
	o.CreateMangaDir = Download.Manga.CreateDir.Get()
	o.CreateVolumeDir = Download.Volume.CreateDir.Get()
	o.SkipIfExists = Download.SkipIfExists.Get()
	o.Strict = Download.Metadata.Strict.Get()
	o.SearchMetadata = Download.Metadata.Search.Get()
	o.DownloadMangaCover = Download.Manga.Cover.Get()
	o.DownloadMangaBanner = Download.Manga.Banner.Get()
	o.WriteSeriesJSON = Download.Metadata.SeriesJSON.Get()
	o.SkipSeriesJSONIfOngoing = Download.Metadata.SkipSeriesJSONIfOngoing.Get()
	o.WriteComicInfoXML = Download.Metadata.ComicInfoXML.Get()
	o.ComicInfoXMLOptions = metadata.DefaultComicInfoOptions()
	o.ImageTransformer = func(bytes []byte) ([]byte, error) {
		return bytes, nil
	}
	return o
}

func ReadOptions() libmangal.ReadOptions {
	o := libmangal.DefaultReadOptions()
	o.SaveHistory = Read.History.Local.Get()
	o.SaveAnilist = Read.History.Anilist.Get()
	return o
}
