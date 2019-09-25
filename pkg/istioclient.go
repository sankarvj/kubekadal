package pkg

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

var (
	virtualServiceName = "default-virtual-service-name"
	weight1            = uint32(50)
	weight2            = uint32(50)
)

//  patchStringValue specifies a patch operation for a string.
type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

//  patchStringValue specifies a patch operation for a uint32.
type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

func setVirtualServiceWeights(client dynamic.Interface, namespace, virtualServiceName string, weight1 uint32, weight2 uint32) error {
	//  Create a GVR which represents an Istio Virtual Service.
	virtualServiceGVR := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "virtualservices",
	}

	//  Weight the two routes - 50/50.
	patchPayload := make([]patchUInt32Value, 2)
	patchPayload[0].Op = "replace"
	patchPayload[0].Path = "/spec/http/0/route/0/weight"
	patchPayload[0].Value = weight1
	patchPayload[1].Op = "replace"
	patchPayload[1].Path = "/spec/http/0/route/1/weight"
	patchPayload[1].Value = weight2
	patchBytes, _ := json.Marshal(patchPayload)

	//  Apply the patch to the 'service2' service.
	_, err := client.Resource(virtualServiceGVR).Namespace(namespace).Patch(virtualServiceName, types.JSONPatchType, patchBytes, metav1.PatchOptions{})
	return err
}
