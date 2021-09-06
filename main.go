package main

import (
	"bytes"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

// resolveJsonpath resolves jsonpath value from the k8s resource object
func resolveJsonpath(obj runtime.Object, jsonpathStr string) (string, error) {
	var buff bytes.Buffer
	jp, err := printers.NewJSONPathPrinter(jsonpathStr)
	if err != nil {
		return "", nil
	}
	err = jp.PrintObj(obj, &buff)
	return buff.String(), err
}

func printJsonpathValues(obj runtime.Object, jsonpathStr string) {
	value, err := resolveJsonpath(obj, jsonpathStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonpathStr, "->", value)
}

func main() {
	// Tests
	printJsonpathValues(getDeploy(), "{.metadata.name}")
	printJsonpathValues(getDeploy(), "{.spec.template.spec.containers[0].image}")
	printJsonpathValues(getDeploy(), "{.spec.replicas}")
	printJsonpathValues(getDeploy(), "{.spec}")
}

func getDeploy() *appsv1.Deployment {
	replicas := int32(2)
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "http",
									HostPort:      0,
									ContainerPort: 80,
									Protocol:      corev1.Protocol("TCP"),
								},
							},
							Resources:       corev1.ResourceRequirements{},
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						},
					},
				},
			},
			Strategy:        appsv1.DeploymentStrategy{},
			MinReadySeconds: 0,
		},
	}
}
