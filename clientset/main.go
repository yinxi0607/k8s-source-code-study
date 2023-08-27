package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./.kube/config")
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	podClient := clientSet.CoreV1().Pods("nhb")
	ctx := context.Background()
	list, err := podClient.List(ctx, metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}
	for _, pod := range list.Items {
		fmt.Printf("Namespace:%v \t Name:%v \t statu:%v \n", pod.Namespace, pod.Name, pod.Status.Phase)
	}

}
