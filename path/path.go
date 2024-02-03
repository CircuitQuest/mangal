package path

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/luevano/mangal/meta"
)

func CacheDir() string {
	dir := filepath.Join(xdg.CacheHome, meta.AppName)
	createDirIfAbsent(dir)
	return dir
}

func ConfigDir() string {
	var dir string

	if runtime.GOOS == "darwin" {
		dir = filepath.Join(xdg.Home, ".config", meta.AppName)
	} else {
		dir = filepath.Join(xdg.ConfigHome, meta.AppName)
	}

	createDirIfAbsent(dir)
	return dir
}

func DownloadsDir() string {
	dir := xdg.UserDirs.Download
	createDirIfAbsent(dir)
	return dir
}

func TempDir() string {
	dir := filepath.Join(os.TempDir(), meta.AppName)
	createDirIfAbsent(dir)
	return dir
}

// TODO: this references ConfigDir() it should instead use the
// configured config dir in case that it is passed as a flag; e.g. use a parameter
func ProvidersDir() string {
	dir := filepath.Join(ConfigDir(), "providers")
	createDirIfAbsent(dir)
	return dir
}

func LogDir() string {
	dir := filepath.Join(TempDir(), "logs")
	createDirIfAbsent(dir)
	return dir
}
