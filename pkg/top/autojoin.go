package top

import (
	"fmt"
	"github.com/readystock/golinq"
	"github.com/readystock/golog"
	"io/ioutil"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func getAutoJoinAddresses() ([]string, error) {
	// Check to see if we are currently running inside of Kubernetes
	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if !ok {
		return nil, fmt.Errorf("auto-join is only supported when running within kubernetes")
	}

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

	pods, err := clientSet.CoreV1().Pods(currentNameSpace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	joinAddresses := make([]string, 0)
	items := make([]v1.Pod, 0)
	linq.From(pods.Items).OrderBy(func(i interface{}) interface{} {
		pod, ok := i.(v1.Pod)
		if !ok {
			return "z"
		} else {
			return pod.Name
		}
	}).ToSlice(&items)

	for _, pod := range items {
		if pod.Name == host {
			// This is a bit weird, but basically if this is the first node in the cluster
			// alphabetically then we want to have it try to be the leader.
			if len(joinAddresses) == 0 {
				golog.Warnf("this node should be the default leader for the cluster")
				return joinAddresses, nil
			}
			continue
		}
		addr := pod.Status.PodIP
		container, ok := linq.From(pod.Spec.Containers).FirstWith(func(i interface{}) bool {
			container, ok := i.(v1.Container)
			return ok && strings.HasPrefix(container.Image, "noahdb/node")
		}).(v1.Container)
		if !ok {
			continue
		}
		hasNoahPort := linq.From(container.Ports).AnyWith(func(i interface{}) bool {
			port, ok := i.(v1.ContainerPort)
			return ok && port.ContainerPort == 5433
		})
		if !hasNoahPort {
			continue
		}
		golog.Debugf("found pod [%s] address: %s", pod.Name, pod.Status.PodIP)
		joinAddresses = append(joinAddresses, fmt.Sprintf("%s:%d", addr, 5433))
	}
	return joinAddresses, nil
}
