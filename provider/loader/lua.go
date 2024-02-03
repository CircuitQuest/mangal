package loader

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/luevano/libmangal"
	"github.com/luevano/luaprovider"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/provider/info"
	"github.com/luevano/mangal/util/afs"
)

const mainLua = "main.lua"

func LuaLoaders() ([]libmangal.ProviderLoader, error) {
	return getLoaderBundles("", config.Config.Providers.Path.Get())
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
		path.ModeDir,
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
		loader, err := newLoader(providerInfo.ProviderInfo, dir)
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

func newLoader(info libmangal.ProviderInfo, dir string) (libmangal.ProviderLoader, error) {
	providerMainFilePath := filepath.Join(dir, mainLua)
	exists, err := afs.Afero.Exists(providerMainFilePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("%s is missing", providerMainFilePath)
	}

	providerMainFileContents, err := afs.Afero.ReadFile(providerMainFilePath)
	if err != nil {
		return nil, err
	}

	options := luaprovider.Options{
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		HTTPStore:    httpStore,
		PackagePaths: []string{dir},
	}

	return luaprovider.NewLoader(providerMainFileContents, info, options)
}
