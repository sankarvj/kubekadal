package controller

import (
	"errors"
	"fmt"

	"github.com/sankarvj/kubekadal/app/utils"
	"github.com/sankarvj/kubekadal/pkg"
	dialogflow "google.golang.org/api/dialogflow/v2"
)

func kubectlMainInfo(dfReq *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (dialogflow.GoogleCloudDialogflowV2WebhookResponse, error) {
	response := "Hello... Vijay, I am your kubernetes assistant. Do you want check my skills?" //pass yes and wait for the follow up intent to fire from dialog flow
	return utils.MakeWebhookResponse(true, response, nil), nil
}

// Get the services running in the cluster
func kubectlServiceInfo(dfReq *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (dialogflow.GoogleCloudDialogflowV2WebhookResponse, error) {
	parametersMap, err := utils.DecodeParameters(dfReq.QueryResult.Parameters)
	if err != nil {
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, err
	}

	response, err := pkg.RunningServices("default")
	if err != nil {
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, err
	}

	fmt.Println("deploying the app using context parameters", parametersMap)
	return utils.MakeWebhookResponse(true, response, nil), nil
}

// The idea is as follows
// If the text has the operator equally, then route traffic 50/50 to v1,v2
// Otherwise route 100% to either v1 or v2 based on the value in the parameter appversion
// Use the application name to find out the virtual service name
func kubectlVirtualServiceRoute(dfReq *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (dialogflow.GoogleCloudDialogflowV2WebhookResponse, error) {
	parametersMap, err := utils.DecodeParameters(dfReq.QueryResult.Parameters)
	if err != nil {
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, err
	}

	var virtualServiceName string
	var weight1 uint32
	var weight2 uint32
	if val, ok := parametersMap["application"]; ok {
		if val == "demo_app" {
			virtualServiceName = "demo-app-vs" // replace this with your virtual service name
		}
	}

	if operatorVal, ok := parametersMap["operator"]; ok {
		if operatorVal == "equally" {
			weight1 = 50
			weight2 = 50
		} else {
			if appversionVal, ok := parametersMap["appversion"]; ok {
				if appversionVal == "v1" {
					weight1 = 100
					weight2 = 0
				} else if appversionVal == "v2" {
					weight1 = 0
					weight2 = 100
				} else {
					return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, errors.New("Sorry. Could you please specify correct app version. It should be either v1 or v2")
				}
			} else {
				return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, errors.New("Sorry. Could you please specify the app version?")
			}
		}
	} else {
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, errors.New("Do you want me to allow all the traffic?")
	}

	if virtualServiceName == "" {
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, errors.New("Sorry. Could you please specify the app name?")
	}

	response, err := pkg.ChangeIstioWeight("default", virtualServiceName, weight1, weight2)
	if err != nil {
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, err
	}

	fmt.Println("deploying the app using context parameters", parametersMap)
	return utils.MakeWebhookResponse(true, response, nil), nil
}
