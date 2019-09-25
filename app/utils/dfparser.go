package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	dialogflow "google.golang.org/api/dialogflow/v2"
)

// ResponsePayload struct
type ResponsePayload struct {
	Google ResponseGoogle `json:"google"`
}

// ResponseGoogle struct
type ResponseGoogle struct {
	ExpectUserResponse bool                  `json:"expectUserResponse"`
	RichResponse       RichResponse          `json:"richResponse"`
	SystemIntent       *ResponseSystemIntent `json:"systemIntent,omitempty"`
}

// RichResponse struct
type RichResponse struct {
	Items []Item `json:"items"`
}

// Item struct
type Item struct {
	SimpleResponse SimpleResponse `json:"simpleResponse"`
}

// SimpleResponse struct
type SimpleResponse struct {
	TextToSpeech string `json:"textToSpeech"`
}

// ResponseSystemIntent struct
type ResponseSystemIntent struct {
	Intent string                   `json:"intent"`
	Data   ResponseSystemIntentData `json:"data"`
}

// ResponseSystemIntentData struct
type ResponseSystemIntentData struct {
	Type        string   `json:"@type"`
	OptContext  string   `json:"optContext"`
	Permissions []string `json:"permissions"`
}

// FollowupEventInput struct
type FollowupEventInput struct {
	Name         string `json:"name"`
	LanguageCode string `json:"languageCode"`
	Parameters   struct {
		Param string `json:"param"`
	} `json:"parameters"`
}

//MakeWebhookResponse make dialogflow webhook response
func MakeWebhookResponse(expectUserResponse bool, msg string, outputCOntexts []*dialogflow.GoogleCloudDialogflowV2Context) dialogflow.GoogleCloudDialogflowV2WebhookResponse {
	payload := ResponsePayload{
		Google: ResponseGoogle{
			ExpectUserResponse: expectUserResponse,
			RichResponse: RichResponse{
				Items: []Item{
					{
						SimpleResponse: SimpleResponse{
							TextToSpeech: msg,
						},
					},
				},
			},
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	return dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		Payload:        payloadBytes,
		OutputContexts: outputCOntexts,
	}
}

//MakeErrorWebhookResponse makes error response if something fails
func MakeErrorWebhookResponse(msg string) dialogflow.GoogleCloudDialogflowV2WebhookResponse {
	payload := ResponsePayload{
		Google: ResponseGoogle{
			ExpectUserResponse: true,
			RichResponse: RichResponse{
				Items: []Item{
					{
						SimpleResponse: SimpleResponse{
							TextToSpeech: msg,
						},
					},
				},
			},
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	return dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		Payload: payloadBytes,
	}
}

//DecodeDialogFlowRequest parses the dialog flow
func DecodeDialogFlowRequest(r *http.Request) (*dialogflow.GoogleCloudDialogflowV2WebhookRequest, error) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("body -- ", string(body))
	r.ParseForm()
	var t *dialogflow.GoogleCloudDialogflowV2WebhookRequest
	err := json.Unmarshal(body, &t)
	return t, err
}

// DecodeParameters decodes the parameters from the webhook request
func DecodeParameters(jsonBytes []byte) (map[string]string, error) {
	fmt.Println("parameters -- ", string(jsonBytes))
	jsonMap := make(map[string]string)
	err := json.Unmarshal(jsonBytes, &jsonMap)
	return jsonMap, err
}
