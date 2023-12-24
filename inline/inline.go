package inline

import (
	"context"
	"fmt"

	"github.com/mangalorg/libmangal"
)

type Options struct {
	Client    *libmangal.Client
	Anilist   *libmangal.Anilist
	Title     string
	Download  bool
	JSON      bool
}

func Run(ctx context.Context, options Options) error {
	fmt.Println("DEBUG: inline.Run()")
	fmt.Println(options)

	return nil
}
