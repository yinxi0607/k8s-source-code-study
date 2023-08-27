package main

/*
	用于发现Kubernetes API Server所支持的资源组，资源版本，资源信息的客户端
*/

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./.kube/config")
	if err != nil {
		panic(err)
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	_, apiResourceLists, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			panic(err)
		}
		for _, apiResource := range apiResourceList.APIResources {
			fmt.Printf("name: %v,group: %v,version: %v\n", apiResource.Name, gv.Group, gv.Version)
		}
	}
}
