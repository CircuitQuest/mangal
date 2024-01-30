package main

import (
	"github.com/luevano/mangal/cmd"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/path"
)

func main() {
	if err := config.Load(path.ConfigDir()); err != nil {
		log.L.Fatal().Err(err).Msg("failed to load config")
		panic(err)
	}

	cmd.Execute()
}
