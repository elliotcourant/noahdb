package top

import (
	"fmt"
	"github.com/readystock/golinq"
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

	for _, pod := range pods.Items {
		if pod.Name == host {
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
		joinAddresses = append(joinAddresses, fmt.Sprintf("%s:%d", addr, 5433))
	}
	return joinAddresses, nil
}
