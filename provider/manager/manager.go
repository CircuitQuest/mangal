package manager

import (
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/provider/bundle"
	"github.com/luevano/libmangal"
)

func Loaders() ([]libmangal.ProviderLoader, error) {
	return bundle.Loaders(path.ProvidersDir())
}
