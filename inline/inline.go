package inline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/luevano/libmangal"
)

type InlineArgs struct {
	Query    string
	Provider string
	Download bool
	JSON     bool
}

type Options struct {
	InlineArgs
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
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

		if len(mangas) == 0 {
			return fmt.Errorf("no mangas found with provider ID %q and query %q", options.Provider, options.Query)
		}

		mangasJSON, err := json.Marshal(mangas)
		fmt.Println(string(mangasJSON))
	}

	return nil
}
