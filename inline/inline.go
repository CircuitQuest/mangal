package inline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mangalorg/libmangal"
)

type Options struct {
	Client   *libmangal.Client
	Anilist  *libmangal.Anilist
	Query    string
	Download bool
	JSON     bool
}

func Run(ctx context.Context, options Options) error {
	switch {
	case options.Download:
		return fmt.Errorf("unimplemented")
	case options.JSON:
		mangas, err := options.Client.SearchMangas(ctx, options.Query)
		if err != nil {
			return err
		}

		for _, manga := range mangas {
			mangaJSON, err := json.Marshal(manga.Info())
			if err != nil {
				return err
			}

			fmt.Println(string(mangaJSON))

			break
		}
	}

	return nil
}
