package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"transmitter-artemis/dto"
)

type MetaClient interface {
	SendRequestToMeta(ctx context.Context, URL string, token string, payload dto.RequestToMeta) (res dto.ResponseFromMeta, httpCode int, err error)
}

type metaClient struct {
	client *http.Client
}

func NewMetaClient() *metaClient {

	client := &http.Client{}

	return &metaClient{
		client: client,
	}
}

func (meta *metaClient) SendRequestToMeta(ctx context.Context, URL string, token string, payload dto.RequestToMeta) (res dto.ResponseFromMeta, httpCode int, err error) {

	httpCode = http.StatusInternalServerError
	res = dto.ResponseFromMeta{}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return dto.ResponseFromMeta{}, httpCode, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewReader(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := meta.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	httpCode = resp.StatusCode

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return
	}

	return
}
