package cmd

import (
	"github.com/readystock/golog"
	"github.com/spf13/cobra"
	"os"
)

var (
	PGWireAddr     string
	GrpcAddr       string
	RaftAddr       string
	JoinAddr       string
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

		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&PGWireAddr, "pg", "p", ":5433", "address that will accept PostgreSQL connections")
	startCmd.Flags().StringVarP(&GrpcAddr, "grpc", "g", ":5434", "address that will be used for Noah's internal gRPC interface")
	startCmd.Flags().StringVarP(&RaftAddr, "raft", "r", ":5435", "address that will be used for Noah's raft protocol")
	startCmd.Flags().StringVarP(&JoinAddr, "join", "j", "", "address and gRPC port of another node in a cluster to join")
	startCmd.Flags().StringVarP(&StoreDirectory, "store", "s", "data", "directory that will be used for Noah's key value store")
	startCmd.Flags().StringVarP(&LogLevel, "log", "l", "verbose", "log output level, valid values: trace, verbose, debug, info, warn, error, fatal")
	rootCmd.AddCommand(startCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		golog.Fatal(err)
		os.Exit(1)
	}
}
