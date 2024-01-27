package config

import "github.com/luevano/libmangal"

// DownloadOptions returns the libmangal.DownloadOptions fields populated by the set Config.
func (c config) DownloadOptions() libmangal.DownloadOptions {
	return libmangal.DownloadOptions{
		Format:              c.Download.Format.Get(),
		Directory:           c.Download.Path.Get(),
		CreateProviderDir:   c.Download.Provider.CreateDir.Get(),
		CreateMangaDir:      c.Download.Manga.CreateDir.Get(),
		CreateVolumeDir:     c.Download.Volume.CreateDir.Get(),
		Strict:              c.Download.Strict.Get(),
		SkipIfExists:        c.Download.SkipIfExists.Get(),
		DownloadMangaCover:  c.Download.Manga.Cover.Get(),
		DownloadMangaBanner: c.Download.Manga.Banner.Get(),
		WriteSeriesJson:     c.Download.Metadata.SeriesJSON.Get(),
		WriteComicInfoXml:   c.Download.Metadata.ComicInfoXML.Get(),
		ReadAfter:           false,
		ReadOptions: libmangal.ReadOptions{
			SaveHistory: c.Read.History.Local.Get(),
			SaveAnilist: c.Read.History.Anilist.Get(),
		},
		ComicInfoXMLOptions: libmangal.DefaultComicInfoOptions(),
		ImageTransformer: func(bytes []byte) ([]byte, error) {
			return bytes, nil
		},
	}
}
