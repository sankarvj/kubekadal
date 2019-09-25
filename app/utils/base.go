package utils

import (
	"encoding/json"
	"fmt"

	dialogflow "google.golang.org/api/dialogflow/v2"
)

//PrintResponse prints the dialogflow response
func PrintResponse(response dialogflow.GoogleCloudDialogflowV2WebhookResponse) {
	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
