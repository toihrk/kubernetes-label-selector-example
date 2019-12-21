package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func main() {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myapp-pod-version3",
			Namespace: metav1.NamespaceDefault,
			Labels: map[string]string{
				"app":        "myapp",
				"appVersion": "3",
			},
		},
		Spec: corev1.PodSpec{},
	}

	labelSet := labels.Set(pod.GetLabels())
	fmt.Printf("Target pod: %s/%s\n", pod.GetNamespace(), pod.GetName())
	fmt.Println("Labels:")
	for key, value := range labelSet {
		fmt.Printf("%s = %s\n", key, value)
	}
	fmt.Println()

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
		&metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      "appVersion",
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"1", "2"},
				},
			},
		},
		&metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      "appVersion",
					Operator: metav1.LabelSelectorOpExists,
				},
			},
		},
		&metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      "foo",
					Operator: metav1.LabelSelectorOpDoesNotExist,
				},
			},
		},
	}

	fmt.Println("Label selectors:")
	for _, labelSelector := range labelSelectors {
		selector, err := metav1.LabelSelectorAsSelector(labelSelector)
		if err != nil {
			panic(err.Error())
		}

		result := "Not match"
		if selector.Matches(labelSet) {
			result = "Match"
		}
		fmt.Printf("%s => %s\n", metav1.FormatLabelSelector(labelSelector), result)
	}
}
