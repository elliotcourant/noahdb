package main

import (
	"github.com/elliotcourant/noahdb/pkg/cmd"
	"github.com/readystock/golog"
)

func main() {
	panic("make a man out of you")
	golog.SetLevel("trace")
	// Main entry point for all of noahdb.
	cmd.Execute()
}
