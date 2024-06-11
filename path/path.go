package path

import (
	"os"
	"path/filepath"

	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
)

func CacheDir() string {
	dir := config.Cache.Path.Get()
	createDirIfAbsent(dir)
	return dir
}

func ConfigDir() string {
	var dir string = filepath.Base(config.Path)
	createDirIfAbsent(dir)
	return dir
}

func DownloadsDir() string {
	dir := config.Download.Path.Get()
	createDirIfAbsent(dir)
	return dir
}

func TempDir() string {
	dir := filepath.Join(os.TempDir(), meta.AppName)
	createDirIfAbsent(dir)
	return dir
}

func ProvidersDir() string {
	dir := config.Providers.Path.Get()
	createDirIfAbsent(dir)
	return dir
}

func LogDir() string {
	dir := filepath.Join(TempDir(), "logs")
	createDirIfAbsent(dir)
	return dir
}
