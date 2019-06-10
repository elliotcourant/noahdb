package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/elliotcourant/noahdb/pkg/rpcwire"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/hashicorp/raft"
	"github.com/readystock/golog"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func NewTestColonyEx(t *testing.T, listenAddr string, spawnPg bool, joinAddresses ...string) (core.Colony, func()) {
	// Create a postgres docker image to connect to.
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	golog.Info("things")

	tempPostgresAddress, tempPostgresPort, tempPostgresUser, tempPostgresPassword := "", 0, "", ""

	callbacks := make([]func(), 0)

	if spawnPg {

		imageName := "docker.io/library/postgres:10"

		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, out)

		testNameCleaned := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
		containerName := strings.ReplaceAll(fmt.Sprintf("noahdb-test-database-%s-%d", testNameCleaned, time.Now().Unix()), "/", "_")

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Env: []string{
				fmt.Sprintf("POSTGRES_PASSWORD=%s", testNameCleaned),
			},
		}, &container.HostConfig{
			PublishAllPorts: true,
		}, nil, containerName)
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		info, err := cli.ContainerInspect(ctx, resp.ID)
		if err != nil {
			golog.Errorf("could not inspect container: %v", err)
		}

		postgresPort := info.NetworkSettings.Ports["5432/tcp"][0].HostPort
		port, err := strconv.Atoi(postgresPort)
		if err != nil {
			golog.Fatalf("could not parse temp postgres port [%s]: %v", postgresPort, err)
			panic(err)
		}

		attempts := 0
		for {
			time.Sleep(2 * time.Second)
			address := fmt.Sprintf("%s:%s", "0.0.0.0", postgresPort)
			connStr := fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=disable", testNameCleaned, address)
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				golog.Warnf("failed to connect to test postgres container address [%s]: %v", address, err)
				attempts++
				if attempts > 3 {
					t.Errorf("could not connect to postgres container in 3 attempts: %v", err)
					panic(err)
				}
				continue
			}

			rows, err := db.Query("SELECT 1")
			if err != nil {
				golog.Warnf("failed to execute simple query to test postgres container address [%s]: %v", address, err)
				attempts++
				if attempts > 3 {
					t.Errorf("could not execute simple query to postgres container in 3 attempts: %v", err)
					panic(err)
				}
				continue
			}

			for rows.Next() {
				one := 0
				rows.Scan(&one)
				if one == 1 {
					goto LeaveLoop
				}
			}

		LeaveLoop:
			if err := rows.Err(); err != nil {
				golog.Warnf("failed to execute simple query to test postgres container address [%s]: %v", address, err)
				attempts++
				if attempts > 3 {
					t.Errorf("could not execute simple query to postgres container in 3 attempts: %v", err)
					panic(err)
				}
				continue
			}

			rows.Close()
			db.Close()
			break
		}

		tempPostgresAddress = "0.0.0.0"
		tempPostgresPort = port
		tempPostgresUser = "postgres"
		tempPostgresPassword = testNameCleaned

		callbacks = append(callbacks, func() {
			timeout := time.Second * 5
			if err := cli.ContainerStop(ctx, resp.ID, &timeout); err != nil {
				t.Fail()
				golog.Criticalf("failed to stop docker container at the end of test [%s]: %v", t.Name(), err)
			}
			if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			}); err != nil {
				t.Fail()
				golog.Criticalf("failed to remove docker container at the end of test [%s]: %v", t.Name(), err)
			}
		})
	}

	golog.SetLevel("trace")
	tempdir, err := ioutil.TempDir("", "core-temp")
	if err != nil {
		panic(err)
	}

	callbacks = append(callbacks, func() {
		if err := os.RemoveAll(tempdir); err != nil {
			panic(err)
		}
	})

	// joins := strings.Join(joinAddresses, ",")

	parsedRaftAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		panic(err)
		// return nil, nil, err
	}

	tn := tcp.NewTransport()

	if err := tn.Open(parsedRaftAddr.String()); err != nil {
		panic(err)
	}

	trans := core.NewTransportWrapper(tn)

	colony := core.NewColony()

	go func() {
		if err = pgwire.NewServer(colony, trans); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	go func() {
		if err = rpcwire.NewRpcServer(colony, trans); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	config := core.ColonyConfig{
		DataDirectory:         tempdir,
		JoinAddresses:         make([]raft.Server, 0),
		Transport:             trans,
		LocalPostgresUser:     tempPostgresUser,
		LocalPostgresAddress:  tempPostgresAddress,
		LocalPostgresPassword: tempPostgresPassword,
		LocalPostgresPort:     tempPostgresPort,
	}

	err = colony.InitColony(config)
	if err != nil {
		panic(err)
	}

	// if colony.IsLeader() {
	// 	_, err = colony.DataNodes().NewDataNode("127.0.0.1", os.Getenv("PGPASSWORD"), os.Getenv("PGPORT"))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	_, err = colony.Shards().NewShard()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	return colony, func() {
		for _, callback := range callbacks {
			callback()
		}
	}
}

func NewTestColony(t *testing.T, joinAddresses ...string) (core.Colony, func()) {
	return NewTestColonyEx(t, ":", false, joinAddresses...)
}

func NewPgTestColony(t *testing.T, joinAddresses ...string) (core.Colony, func()) {
	return NewTestColonyEx(t, ":", true, joinAddresses...)
}

func ConnectionString(address net.Addr) string {
	addr, err := net.ResolveTCPAddr(address.Network(), address.String())
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		addr.IP.String(), addr.Port, "noah", "password", "postgres")
}
