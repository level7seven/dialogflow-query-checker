package query

import (
	"github.com/yoichiro/dialogflow-query-checker/config"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
)

type RequestBody struct {
	Contexts []string `json:"contexts"`
	Language string `json:"lang"`
	Query string `json:"query"`
	SessionId string `json:"sessionId"`
}

func Execute(test *config.Test, clientAccessToken string) (*Response, error) {
	res, err := send(test, clientAccessToken)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func send(test *config.Test, clientAccessToken string) (*http.Response, error) {
	requestBody := RequestBody{
		Contexts: test.Condition.Contexts,
		Language: "ja",
		Query: test.Condition.Query,
		SessionId: "1",
	}
	body, err := json.Marshal(&requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.dialogflow.com/v1/query?v=20150910", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer " + clientAccessToken)
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}