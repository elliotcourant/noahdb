package main

import (
	"github.com/elliotcourant/noahdb/pkg/cmd"
	"github.com/elliotcourant/noahdb/pkg/util"
	"github.com/readystock/golog"
)

func main() {
	golog.SetLevel("trace")

	ip, err := util.ExternalIP()
	if err != nil {
		panic(err)
	}

	golog.Warnf("starting noahdb node with IP address: %s", ip)

	// Main entry point for all of noahdb.
	cmd.Execute()
}
