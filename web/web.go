package web

import "github.com/skratchdot/open-golang/open"

type Args struct {
	Open bool
	Port string
}

func Run(args Args) error {
	// TODO: should this be done after server.Start? Maybe as a defer?
	if args.Open {
		open.Start("http://localhost:" + args.Port)
	}

	server, err := NewServer()
	if err != nil {
		return err
	}

	return server.Start(":" + args.Port)
}
