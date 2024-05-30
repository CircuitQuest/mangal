package config

import (
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/metadata"
)

// DownloadOptions constructs the libmangal.DownloadOptions populated by the Config.
func (c config) DownloadOptions() libmangal.DownloadOptions {
	// start from the defaults in case of new additions and build on top of it.
	o := libmangal.DefaultDownloadOptions()
	o.Format = c.Download.Format.Get()
	o.Directory = c.Download.Path.Get()
	o.CreateProviderDir = c.Download.Provider.CreateDir.Get()
	o.CreateMangaDir = c.Download.Manga.CreateDir.Get()
	o.CreateVolumeDir = c.Download.Volume.CreateDir.Get()
	o.Strict = c.Download.Strict.Get()
	o.SkipIfExists = c.Download.SkipIfExists.Get()
	o.SearchMissingMetadata = c.Download.Metadata.SearchMissingMetadata.Get()
	o.DownloadMangaCover = c.Download.Manga.Cover.Get()
	o.DownloadMangaBanner = c.Download.Manga.Banner.Get()
	o.WriteSeriesJSON = c.Download.Metadata.SeriesJSON.Get()
	o.SkipSeriesJSONIfOngoing = c.Download.Metadata.SkipSeriesJSONIfOngoing.Get()
	o.WriteComicInfoXML = c.Download.Metadata.ComicInfoXML.Get()
	o.ReadAfter = false
	o.ReadOptions = libmangal.ReadOptions{
		SaveHistory: c.Read.History.Local.Get(),
		SaveAnilist: c.Read.History.Anilist.Get(),
	}
	o.ComicInfoXMLOptions = metadata.DefaultComicInfoOptions()
	o.ImageTransformer = func(bytes []byte) ([]byte, error) {
		return bytes, nil
	}
	return o
}
