package main

import (
	"github.com/luevano/mangal/cmd"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
)

func main() {
	if err := config.Load(); err != nil {
		log.L.Fatal().Err(err).Msg("failed to load config")
		panic(err)
	}

	cmd := cmd.NewCmd()
	cmd.Execute()
}
