package main

import (
	"github.com/elliotcourant/noahdb/pkg/cmd"
	"github.com/elliotcourant/noahdb/pkg/util"
	"github.com/elliotcourant/timber"
)

func main() {
	ip, err := util.ExternalIP()
	if err != nil {
		panic(err)
	}

	timber.Warningf("starting noahdb node with IP address: %s", ip)

	// Main entry point for all of noahdb.
	cmd.Execute()
}
