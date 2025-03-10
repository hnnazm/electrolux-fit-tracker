package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"fit-tracker/database"

	"github.com/gorilla/websocket"
)

const INTERVAL = 60 * time.Second

type (
	handler struct {
		httpClient *http.Client
		wsClient   *websocket.Dialer
	}

	CheckHealthResult struct {
		IsUp bool `json:"isUp"`
	}

	GenerateTokenInput struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}

	GenerateTokenResult struct {
		AccessToken string `json:"accessToken"`
	}

	GetUserInfoInput struct {
		UserID      string `json:"userId"`
		AccessToken string
	}

	GetUserInfoResult struct {
		ID     string `json:"id"`
		Weight int64  `json:"weight"`
		Height int64  `json:"height"`
	}

	PollTracesInput struct {
		AccessToken string
		DataCh      chan<- *database.SaveTracesInput
	}
)

var (
	VERSION           = "v1"
	DOMAIN_NAME       = "fit-tracker-htmz.onrender.com"
	HTTP_ENDPOINT_URL = "https://" + DOMAIN_NAME + "/api/" + VERSION
	WSS_ENDPOINT_URL  = "wss://" + DOMAIN_NAME + "/api/" + VERSION
)

func New(httpClient *http.Client, wsClient *websocket.Dialer) *handler {
	return &handler{
		httpClient: httpClient,
		wsClient:   wsClient,
	}
}

func (h handler) CheckHealth() bool {
	res, err := h.httpClient.Get(HTTP_ENDPOINT_URL + "/health")
	if err != nil {
		return false
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false
	}

	var result CheckHealthResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	return result.IsUp
}

func (h handler) GenerateAccessToken(input GenerateTokenInput) (string, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", HTTP_ENDPOINT_URL+"/token", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result GenerateTokenResult
	json.Unmarshal(body, &result)

	return result.AccessToken, nil
}

func (h handler) GetUserInfo(ctx context.Context, input *GetUserInfoInput) (*GetUserInfoResult, error) {
	req, err := http.NewRequest("GET", HTTP_ENDPOINT_URL+"/users/"+input.UserID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+input.AccessToken)
	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result GetUserInfoResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (h handler) PollTraces(ctx context.Context, input *PollTracesInput) {
	var headers = http.Header{}
	headers.Set("Authorization", input.AccessToken)
	headers.Set("Accept", "*/*")

	conn, _, err := h.wsClient.Dial(WSS_ENDPOINT_URL+"/traces", headers)
	if err != nil {
		log.Println("error dialing websocket", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		data := decodeBase64(string(message))
		input.DataCh <- data
	}
}

// decode base64 to string
func decodeBase64(data string) *database.SaveTracesInput {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println("error decoding base64", err)
		return nil
	}

	var result *database.SaveTracesInput
	err = json.Unmarshal(decoded, &result)
	if err != nil {
		log.Println("error unmarshalling json", err)
		return nil
	}

	return result
}
