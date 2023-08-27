package main

/*
  dynamicClient.Resource 函数用于设置请求的资源组，资源版本，资源名称，可以用于crd的自定义资源，
  unstructured.Unstructured 用于解析请求的资源对象
*/
import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./.kube/config")
	if err != nil {
		panic(err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	ctx := context.TODO()
	unstructObj, err := dynamicClient.Resource(gvr).Namespace("nhb").List(ctx, metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
	if err != nil {
		panic(err)
	}
	for _, pod := range podList.Items {
		fmt.Printf("Namespace:%v \t Name:%v \t statu:%v \n", pod.Namespace, pod.Name, pod.Status.Phase)
	}
}
