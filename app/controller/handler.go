package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sankarvj/kubekadal/app/utils"
	dialogflow "google.golang.org/api/dialogflow/v2"
)

//DialogFlowHandler handles webhook call from the dialog flow
func DialogFlowHandler(w http.ResponseWriter, r *http.Request) {
	var response dialogflow.GoogleCloudDialogflowV2WebhookResponse
	dfReq, err := utils.DecodeDialogFlowRequest(r)
	if err != nil {
		response = utils.MakeErrorWebhookResponse(err.Error())
	}
	response, err = processInput(dfReq)
	if err != nil {
		response = utils.MakeErrorWebhookResponse(err.Error())
	}
	utils.PrintResponse(response)
	json.NewEncoder(w).Encode(response)
}

func processInput(dfReq *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (dialogflow.GoogleCloudDialogflowV2WebhookResponse, error) {
	switch dfReq.QueryResult.Intent.Name {
	case utils.KubeCtlMain:
		return kubectlMainInfo(dfReq)
	case utils.KubeCtlServiceInfo:
		return kubectlServiceInfo(dfReq)
	case utils.KubeCtlVirtualService:
		return kubectlVirtualServiceRoute(dfReq)
	default:
		return dialogflow.GoogleCloudDialogflowV2WebhookResponse{}, errors.New("No intent matches. Please add/change intents in kubekadal/app/utils ")
	}
}
