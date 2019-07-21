package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/timber"
	"github.com/hashicorp/raft"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Kubernetes API authentication.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func getAutoJoinAddresses() ([]raft.Server, error) {
	// Check to see if we are currently running inside of Kubernetes
	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if !ok {
		return nil, fmt.Errorf("auto-join is only supported when running within kubernetes")
	}

	timber.Debugf("waiting a few seconds to give Kubernetes a chance to start any other pods")
	time.Sleep(5 * time.Second)

	host := os.Getenv("HOSTNAME")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cn, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, err
	}

	currentNameSpace := string(cn)

	deployments, err := clientSet.AppsV1().Deployments(currentNameSpace).List(metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	noahDeploymentIndex := linq.From(deployments.Items).IndexOf(func(i interface{}) bool {
		deployment, ok := i.(appsV1.Deployment)
		return ok && deployment.Name == "noahdb"
	})

	if noahDeploymentIndex < 0 {
		return nil, fmt.Errorf("could not find a noahdb deployment, auto-join not supported")
	}

	var replicas = 1
	if deployments.Items[noahDeploymentIndex].Spec.Replicas != nil {
		replicas = int(*deployments.Items[noahDeploymentIndex].Spec.Replicas)
	}

	retries, maxRetries := 0, 3
GetPods:
	pods, err := clientSet.CoreV1().Pods(currentNameSpace).List(metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	joinAddresses := make([]raft.Server, 0)
	items := make([]coreV1.Pod, 0)
	linq.From(pods.Items).OrderBy(func(i interface{}) interface{} {
		pod, ok := i.(coreV1.Pod)
		if !ok {
			return "z"
		}
		return pod.Name
	}).ToSlice(&items)

	if len(items) < replicas {
		timber.Warningf("looking for [%d] replicas, but only found [%d]", replicas, len(items))
		if retries < maxRetries {
			time.Sleep(5 * time.Second)
			retries++
			goto GetPods
		}
	}

	for _, pod := range items {
		if pod.Name == host {
			continue
		}

		addr := pod.Status.PodIP
		container, ok := linq.From(pod.Spec.Containers).FirstWith(func(i interface{}) bool {
			container, ok := i.(coreV1.Container)
			return ok && strings.HasPrefix(container.Image, "noahdb/node")
		}).(coreV1.Container)
		if !ok {
			continue
		}
		containerPort, hasNoahPort := linq.From(container.Ports).FirstWith(func(i interface{}) bool {
			port, ok := i.(coreV1.ContainerPort)
			return ok && port.Name == "noahdb"
		}).(coreV1.ContainerPort)

		if !hasNoahPort || addr == "" {
			continue
		}

		processedAddress := fmt.Sprintf("%s:%d", addr, containerPort.ContainerPort)

		timber.Debugf("found pod [%s] address: %s", pod.Name, processedAddress)
		joinAddresses = append(joinAddresses, raft.Server{
			ID:       raft.ServerID(pod.Name),
			Address:  raft.ServerAddress(processedAddress),
			Suffrage: raft.Voter,
		})
	}
	return joinAddresses, nil
}
