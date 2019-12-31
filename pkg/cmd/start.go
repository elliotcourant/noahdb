package cmd

import (
	"github.com/elliotcourant/noahdb/pkg/top"
	"github.com/elliotcourant/timber"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	pgListenAddr    string
	raftListenAddr  string
	joinAddr        string
	autoDataNode    bool
	autoJoinCluster bool
	storeDirectory  string
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
			options := top.Options{
				DataDirectory:     storeDirectory,
				JoinAddresses:     strings.Split(joinAddr, ","),
				PgListenAddress:   pgListenAddr,
				RaftListenAddress: raftListenAddr,
				AutoDataNode:      autoDataNode,
				AutoJoinCluster:   autoJoinCluster,
			}

			top.NoahMain(options)

		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&pgListenAddr, "pg-listen", "L", ":5433", "address that will accept connections")
	startCmd.Flags().StringVarP(&raftListenAddr, "raft-listen", "R", ":5434", "address that will be use for raft")
	startCmd.Flags().StringVarP(&joinAddr, "join", "J", "", "address of another node in the cluster to use to join")
	startCmd.Flags().BoolVarP(&autoDataNode, "auto-data-node", "d", false, "look for a local PostgreSQL instance")
	startCmd.Flags().BoolVarP(&autoJoinCluster, "auto-join", "A", false, "try to auto-join an existing cluster")
	startCmd.Flags().StringVarP(&storeDirectory, "dir", "s", "data", "directory that will be used for Noah's key value store")
	rootCmd.AddCommand(startCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		timber.Fatal(err)
		os.Exit(1)
	}
}
