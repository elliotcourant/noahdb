package cmd

import (
	"github.com/elliotcourant/noahdb/pkg/top"
	"github.com/elliotcourant/timber"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	PgListenAddr   string
	RaftListenAddr string
	JoinAddr       string
	AutoDataNode   bool
	AutoJoin       bool
	UseTmpDir      bool
	StoreDirectory string
	LogLevel       string
)

var (
	rootCmd = &cobra.Command{
		Use:   "noahdb",
		Short: "noahdb is a multi-tenant new-SQL database.",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	startCmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			StartDB(StoreDirectory, JoinAddr, PgListenAddr, UseTmpDir, AutoDataNode, AutoJoin)
		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&PgListenAddr, "pg-listen", "L", ":5433", "address that will accept connections")
	startCmd.Flags().StringVarP(&RaftListenAddr, "raft-listen", "R", ":5434", "address that will be use for raft")
	startCmd.Flags().StringVarP(&JoinAddr, "join", "J", "", "address of another node in the cluster to use to join")
	startCmd.Flags().BoolVarP(&AutoDataNode, "auto-data-node", "d", false, "look for a local PostgreSQL instance")
	startCmd.Flags().BoolVarP(&AutoJoin, "auto-join", "A", false, "try to auto-join an existing cluster")
	startCmd.Flags().BoolVarP(&UseTmpDir, "temp", "t", false, "use temp directory each time")
	startCmd.Flags().StringVarP(&StoreDirectory, "store", "s", "data", "directory that will be used for Noah's key value store")
	startCmd.Flags().StringVarP(&LogLevel, "log", "l", "verbose", "log output level, valid values: trace, verbose, debug, info, warn, error, fatal")
	rootCmd.AddCommand(startCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		timber.Fatal(err)
		os.Exit(1)
	}
}

func StartDB(storeDirectory, joinAddr, listenAddr string, useTempDir, autoDataNode, autoJoin bool) {
	if useTempDir {
		tempdir, err := ioutil.TempDir("", "noahdb")
		if err != nil {
			panic(err)
		}
		storeDirectory = tempdir
		defer func() {
			timber.Infof("cleaning up temp directory: %s", tempdir)
			os.RemoveAll(tempdir)
		}()
	}

	options := top.Options{
		DataDirectory:     storeDirectory,
		JoinAddresses:     nil,
		PgListenAddress:   "",
		RaftListenAddress: "",
		AutoDataNode:      false,
		AutoJoinCluster:   false,
	}

	top.NoahMain(options)
}
