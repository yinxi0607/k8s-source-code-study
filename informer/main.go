package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"strings"
	"time"
)

func process() {
	config, err := clientcmd.BuildConfigFromFlags("", "./.kube/config")
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	stopCh := make(chan struct{})
	defer close(stopCh)
	sharedInformers := informers.NewSharedInformerFactory(clientSet, time.Minute)
	informer := sharedInformers.Core().V1().Pods().Informer()
	_, _ = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s", mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(metav1.Object)
			nObj := newObj.(metav1.Object)
			log.Printf("Pod Updated in Store: Old: %s, New: %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			log.Printf("Pod Deleted from Store: %s", mObj.GetName())
		},
	})

	informer.Run(stopCh)
}

func UsersIndexFunc(obj interface{}) ([]string, error) {
	pod := obj.(*corev1.Pod)
	usersString := pod.Annotations["users"]
	return strings.Split(usersString, ","), nil
}

func IndexTest() {
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"byUser": UsersIndexFunc})
	pod1 := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "one", Annotations: map[string]string{"users": "ernie,bert"}}}
	pod2 := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "two", Annotations: map[string]string{"users": "bert,oscar"}}}
	pod3 := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "tre", Annotations: map[string]string{"users": "ernie,elmo"}}}
	index.Add(pod1)
	index.Add(pod2)
	index.Add(pod3)
	erniePods, err := index.ByIndex("byUser", "ernie")
	if err != nil {
		panic(err)
	}
	for _, erniePod := range erniePods {
		fmt.Println(erniePod.(*corev1.Pod).Name)
	}
}

func main() {
	IndexTest()
}
