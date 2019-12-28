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
	"github.com/elliotcourant/noahdb/pkg/transport"
	"github.com/elliotcourant/noahdb/pkg/util"
	"github.com/elliotcourant/timber"
	"github.com/hashicorp/raft"
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

type TestDataNode struct {
	Address  string
	Port     int32
	User     string
	Password string
}

func NewDataNode(t *testing.T) (TestDataNode, func(), error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	testNameCleaned := GetCleanTestName(t)

	imageName := "docker.io/library/postgres:12"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	containerName := strings.ReplaceAll(fmt.Sprintf("noahdb-test-db-%s-%d", testNameCleaned, time.Now().Unix()), "/", "_")

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
		timber.Errorf("could not inspect container: %v", err)
	}
	netInfo := info.NetworkSettings.Ports["5432/tcp"][0]
	postgresPort := netInfo.HostPort
	postgresAddress := netInfo.HostIP

	port, err := strconv.ParseInt(postgresPort, 10, 32)
	if err != nil {
		timber.Fatalf("could not parse temp postgres port [%s]: %v", postgresPort, err)
		panic(err)
	}
	address := fmt.Sprintf("%s:%s", postgresAddress, postgresPort)
	timber.Warningf("USING [%s] AS POSTGRES TEMP DB", address)

	attempts := 0
	maxAttempts := 10
	connStr := fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=disable", testNameCleaned, address)
	for {
		time.Sleep(5 * time.Second)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			timber.Warningf("failed to connect to test postgres container address [%s]: %v", address, err)
			attempts++
			if attempts > maxAttempts {
				t.Errorf("could not connect to postgres container in %d attempts: %v", maxAttempts, err)
				panic(err)
			}
			continue
		}

		rows, err := db.Query("SELECT 1")
		if err != nil {
			timber.Warningf("failed to execute simple query to test postgres container address [%s]: %v", address, err)
			attempts++
			if attempts > maxAttempts {
				t.Errorf("could not execute simple query to postgres container in %d attempts: %v", maxAttempts, err)
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
			timber.Warningf("failed to execute simple query to test postgres container address [%s]: %v", address, err)
			attempts++
			if attempts > maxAttempts {
				t.Errorf("could not execute simple query to postgres container in %d attempts: %v", maxAttempts, err)
				panic(err)
			}
			continue
		}

		rows.Close()
		db.Close()
		break
	}

	timber.Warningf("Temp DB Address: %s", connStr)

	node := TestDataNode{
		Address:  "0.0.0.0",
		Port:     int32(port),
		User:     "postgres",
		Password: testNameCleaned,
	}
	callbacks := make([]func(), 0)
	callbacks = append(callbacks, func() {
		timber.Infof("cleaning up postgres test node")
		timeout := time.Second * 5
		if err := cli.ContainerStop(ctx, resp.ID, &timeout); err != nil {
			t.Fail()
			timber.Criticalf("failed to stop docker container at the end of test [%s]: %v", t.Name(), err)
		}
		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}); err != nil {
			t.Fail()
			timber.Criticalf("failed to remove docker container at the end of test [%s]: %v", t.Name(), err)
		}
	})

	return node, func() {
		for _, callback := range callbacks {
			callback()
		}
	}, nil
}

func GetCleanTestName(t *testing.T) string {
	return strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
}

func NewTestColonyEx(t *testing.T, listenAddr string, spawnPg bool, joinAddresses ...string) (core.Colony, func()) {
	log := timber.New()

	tempPostgresAddress, tempPostgresPort, tempPostgresUser, tempPostgresPassword := "", int32(0), "", ""
	testNameCleaned := fmt.Sprintf("%s-%d", GetCleanTestName(t), time.Now().Unix())

	callbacks := make([]func(), 0)

	if spawnPg {
		node, pgCleanup, err := NewDataNode(t)
		if err != nil {
			panic(err)
		}
		callbacks = append(callbacks, pgCleanup)
		tempPostgresAddress, tempPostgresPort, tempPostgresUser, tempPostgresPassword =
			node.Address, node.Port, node.User, node.Password
	}

	tempdir, err := ioutil.TempDir("", testNameCleaned)
	if err != nil {
		panic(err)
	}

	callbacks = append(callbacks, func() {
		log.Warningf("removing temporary directory at: %s", tempdir)
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

	trans := transport.NewTransportWrapper(tn)

	colony := core.NewColony()

	go func() {
		if err = pgwire.RunServer(colony, trans); err != nil {
			log.Errorf(err.Error())
		}
	}()

	go func() {
		if err = rpcwire.NewRpcServer(colony, trans); err != nil {
			log.Errorf(err.Error())
		}
	}()

	joins := make([]raft.Server, 0)
	for _, joinAddress := range joinAddresses {
		if addr, err := util.ResolveAddress(joinAddress); err != nil {
			log.Errorf("failed to parse join address [%s]: %v", joinAddress, err)
			t.Errorf("failed to parse join address [%s]: %v", joinAddress, err)
		} else {
			joins = append(joins, raft.Server{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(addr),
				Address:  raft.ServerAddress(addr),
			})
		}
	}

	config := core.ColonyConfig{
		DataDirectory:         tempdir,
		JoinAddresses:         joins,
		Transport:             trans,
		LocalPostgresUser:     tempPostgresUser,
		LocalPostgresAddress:  tempPostgresAddress,
		LocalPostgresPassword: tempPostgresPassword,
		LocalPostgresPort:     tempPostgresPort,
	}

	err = colony.InitColony(config, log)
	if err != nil {
		panic(err)
	}

	callbacks = append(callbacks, func() {
		colony.Close()
	})

	log.Infof("finished starting noahdb coordinator")
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

	host := addr.IP.String()
	if host == "::" {
		host, _ = util.ExternalIP()
	}

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, addr.Port, "noah", "password", "postgres")
}
