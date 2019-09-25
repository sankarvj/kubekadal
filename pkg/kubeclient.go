package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// UpdateDeployment update the deployment with the image tag
func UpdateDeployment(namespace, deploymentName, imageTag string) error {
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	client, err := dynamicClient()
	if err != nil {
		return err
	}
	result, err := client.Resource(deploymentRes).Namespace(namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// extract spec containers
	containers, found, err := unstructured.NestedSlice(result.Object, "spec", "template", "spec", "containers")
	if err != nil || !found || containers == nil {
		fmt.Printf("deployment containers not found or error in spec: %v", err)
		return err
	}

	// update container[0] image
	if err := unstructured.SetNestedField(containers[0].(map[string]interface{}), imageTag, "image"); err != nil {
		return err
	}
	if err := unstructured.SetNestedField(result.Object, containers, "spec", "template", "spec", "containers"); err != nil {
		return err
	}

	updateResult, err := client.Resource(deploymentRes).Namespace(namespace).Update(result, metav1.UpdateOptions{})
	fmt.Println("result ", updateResult)
	return err
}

//RunningServices fetches the services running for the specific namespace in current cluster
func RunningServices(namespace string) (string, error) {
	client, err := staticClient()
	if err != nil {
		return "", err
	}
	result, err := client.CoreV1().Services(namespace).List(metav1.ListOptions{})
	fmt.Printf("result %+v\n", result)
	totalSvcs := len(result.Items)
	fmt.Println("totalSvcs", totalSvcs)

	var serviceNames string
	for _, service := range result.Items {
		serviceNames = service.ObjectMeta.Name + "," + serviceNames
		fmt.Println("service name", service.ObjectMeta.Name)
		fmt.Println("service namespace", service.ObjectMeta.Namespace)
	}

	s := fmt.Sprintf("you have %d running service in kubernetes cluster. The service name is %s", totalSvcs, serviceNames)

	return s, nil
}

//ChangeIstioWeight modify weight of to routes in a virtual service
func ChangeIstioWeight(namespace, virtualServiceName string, weight1 uint32, weight2 uint32) (string, error) {
	client, err := dynamicClient()
	if err != nil {
		return "", err
	}

	err = setVirtualServiceWeights(client, namespace, virtualServiceName, weight1, weight2)
	if err != nil {
		return "", err
	}

	var response string
	if weight1 == 100 {
		response = "All your traffic has been routed to V1 version of demo app"
	} else if weight2 == 100 {
		response = "All your traffic has been routed to V2 version of demo app"
	} else if weight1 == 50 && weight2 == 50 {
		response = "Your traffic is now routed to both V1 and V2 versions of your demo app. "
	} else {
		response = fmt.Sprintf("You have set the weight %d for v1 and weight %d for v2", weight1, weight2)
	}

	return response, nil
}

func clientSet() (*rest.Config, error) {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func dynamicClient() (dynamic.Interface, error) {
	config, err := clientSet()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func staticClient() (*kubernetes.Clientset, error) {
	config, err := clientSet()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
