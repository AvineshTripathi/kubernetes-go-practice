package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "/home/avinesh/.kube/config", "locaation of config file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Errorf("1 error")
	}

	config.APIPath = "/api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	i64Ptr := func(i int64) *int64 { return &i }

	opts := metav1.ListOptions{
		TimeoutSeconds: i64Ptr(120),
		Watch:          true,
	}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		fmt.Errorf("restclient error")
	}


	watcher, err := restClient.Get().Resource("pods").VersionedParams(&opts, scheme.ParameterCodec).Timeout(time.Duration(*opts.TimeoutSeconds)).Watch(context.TODO())
	if err != nil {
		fmt.Errorf("watcher")
	}

	for event := range watcher.ResultChan() {

		item := event.Object.(*corev1.Pod)

		switch event.Type {
		case watch.Added:
			fmt.Printf("+ '%s' %v \n", item.Name, event.Type)
		case watch.Deleted:
			fmt.Printf("- '%s' %v \n", item.Name, event.Type)
		case watch.Modified:
			fmt.Printf("* '%s' %v \n", item.Name, event.Type)

		}
	}

}
