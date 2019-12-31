package kube

import (
	"github.com/ahmetb/go-linq/v3"
	"k8s.io/api/core/v1"
	"os"
	"strings"
	"time"

	"fmt"
	"github.com/readystock/golog"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// I'm the eye hole man
func RunEyeholes(colony core.Colony) {
	return
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	host := os.Getenv("HOSTNAME")

	golog.Infof("running watcher from: %s", host)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	cn, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		panic(err)
	}

	currentNameSpace := string(cn)
	golog.Infof("current name space: %s", currentNameSpace)
	for {
		time.Sleep(1 * time.Minute)
		if colony == nil {
			continue
		}
		if !colony.IsLeader() {
			golog.Debugf("looking for friendly neighbors")
			neighborItems, err := colony.Neighbors()
			if err != nil {
				panic(err)
			}

			neighborMap := map[string]interface{}{}
			for _, neighbor := range neighborItems {
				neighborMap[neighbor.Addr] = nil
			}

			pods, err := clientset.CoreV1().Pods(currentNameSpace).List(metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}

			golog.Debugf("found %d pod(s)", len(pods.Items))

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
				podAddress := fmt.Sprintf("%s:%d", addr, 5433)
				_, ok = neighborMap[podAddress]
				if ok {
					continue
				}

				// The pod is not in our friendly neighborhood. Add it.
				golog.Infof("pod [%s] is not in our friendly neighborhood, asking it to move in with address: %s", pod.Name, podAddress)
				if err := colony.Join(podAddress, podAddress); err != nil {
					golog.Errorf("pod [%s] with address %s could not join the neighborhood", pod.Name, podAddress)
				}
			}
		}
	}
}
