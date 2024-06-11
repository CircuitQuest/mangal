package path

import (
	"log"

	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/util/afs"
)

func createDirIfAbsent(path string) {
	exists, err := afs.Afero.Exists(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	if exists {
		return
	}

	if err := afs.Afero.MkdirAll(path, config.Download.ModeDir.Get()); err != nil {
		log.Fatal(err)
	}

	return
}
