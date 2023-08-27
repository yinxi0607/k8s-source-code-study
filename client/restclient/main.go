package main

/*
	restClient是最基础的客户端，能够处理多种类型的调用，返回不同的数据格式
*/

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./.kube/config")
	if err != nil {
		panic(err)
	}
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	result := &corev1.PodList{}
	ctx := context.TODO()
	err = restClient.Get().
		Namespace("nhb").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		panic(err)
	}
	for _, pod := range result.Items {
		fmt.Printf("Namespace:%v \t Name:%v \t statu:%v \n", pod.Namespace, pod.Name, pod.Status.Phase)
	}

}
