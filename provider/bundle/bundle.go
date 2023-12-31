package bundle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/provider/info"
	"github.com/luevano/mangal/provider/lua"
)

func Loaders(dir string) ([]libmangal.ProviderLoader, error) {
	return getLoaderBundles("", dir)
}

// getLoaderBundles returns a linear list of all providers with their parent bundle id attached.
// A bundle is a collection of providers, any amount of nested providers (or even bundles) is valid.
func getLoaderBundles(bundleID, dir string) ([]libmangal.ProviderLoader, error) {
	dirEntries, err := afs.Afero.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var loaderBundles []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		// skip non directories
		if !dirEntry.IsDir() {
			continue
		}

		dirEntryPath := filepath.Join(dir, dirEntry.Name())

		isProvider, err := afs.Afero.Exists(filepath.Join(dirEntryPath, info.Filename))
		if err != nil {
			return nil, err
		}

		// only interested on metadata files 'mangal.toml'
		if !isProvider {
			continue
		}

		loaders, err := getLoaders(bundleID, dirEntryPath)
		if err != nil {
			return nil, err
		}

		loaderBundles = append(loaderBundles, loaders...)
	}

	return loaderBundles, nil
}

// getLoaders returns either the providerloader or the next subset of bundles.
// It all depends if the provider type is Lua (actual provider loader) or bundle.
func getLoaders(bundleID, dir string) ([]libmangal.ProviderLoader, error) {
	infoFile, err := afs.Afero.OpenFile(
		filepath.Join(dir, info.Filename),
		os.O_RDONLY,
		0755,
	)
	if err != nil {
		return nil, err
	}

	defer infoFile.Close()

	providerInfo, err := info.New(infoFile)
	if err != nil {
		return nil, err
	}

	// to make provider ids unique, attach the bundle id
	if bundleID != "" {
		providerInfo.ID = fmt.Sprint(bundleID, "-", providerInfo.ID)
	}

	switch providerInfo.Type {
	case info.TypeLua:
		loader, err := lua.NewLoader(providerInfo.ProviderInfo, dir)
		if err != nil {
			return nil, err
		}

		return []libmangal.ProviderLoader{loader}, nil
	case info.TypeBundle:
		return getLoaderBundles(providerInfo.ID, dir)
	default:
		return nil, fmt.Errorf("unkown provider type: %#v", providerInfo.Type)
	}
}
