package kube

import (
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"os"
)

func GetPeerAddresses(podLabel string) ([]string, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	host := os.Getenv("HOSTNAME")

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cn, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, err
	}

	currentNameSpace := string(cn)

	pods, err := clientset.CoreV1().Pods(currentNameSpace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	peers := make([]string, 0)

	for _, pod := range pods.Items {
		if pod.Name == host {
			continue
		}

		peers = append(peers, pod.Status.PodIP)
	}

	return peers, nil
}
