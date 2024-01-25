package path

import (
	"log"

	"github.com/luevano/mangal/afs"
)

const (
	ModeDir  = 0755
	ModeFile = 0644
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

	if err := afs.Afero.MkdirAll(path, ModeDir); err != nil {
		log.Fatal(err)
	}

	return
}
