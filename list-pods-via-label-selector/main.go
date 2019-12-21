package main

import (
	"fmt"

	flag "github.com/spf13/pflag"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var (
		apiServerHost, kubeconfigFilePath string
	)

	flag.StringVar(&apiServerHost, "master", "", "Address of the Kubernetes API server.")
	flag.StringVar(&kubeconfigFilePath, "kubeconfig", "", "Path to a kubeconfig file containing authorization and API server information.")
	flag.Parse()

	kubeconfig, err := clientcmd.BuildConfigFromFlags(apiServerHost, kubeconfigFilePath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	versionInfo, err := clientset.ServerVersion()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Server version: %s\n", versionInfo.String())

	fmt.Println("All pods")
	podListAll, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if len(podListAll.Items) == 0 {
		fmt.Println("Not Found")
	}
	for _, pod := range podListAll.Items {
		fmt.Printf("%s/%s\n", pod.GetNamespace(), pod.GetName())
	}
	fmt.Println("")

	labelSelectors := []*metav1.LabelSelector{
		&metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "myapp",
			},
		},
		&metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      "appVersion",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"1", "2"},
				},
			},
		},
	}

	for _, labelSelector := range labelSelectors {
		selector := metav1.FormatLabelSelector(labelSelector)
		fmt.Printf("Selector: %s\n", selector)
		podList, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			panic(err.Error())
		}
		if len(podList.Items) == 0 {
			fmt.Println("Not Found")
		}
		for _, pod := range podList.Items {
			fmt.Printf("%s/%s\n", pod.GetNamespace(), pod.GetName())
		}
		fmt.Println("")
	}
}
