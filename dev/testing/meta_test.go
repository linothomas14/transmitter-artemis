package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"transmitter-artemis/dto"
)

func MetaHandler(w http.ResponseWriter, r *http.Request) {
	const bearerPrefix = "Bearer "
	authHeader := r.Header.Get("Authorization")

	token := authHeader[len(bearerPrefix):]
	// params := r.URL.Path[0:]
	// phoneNumberID := strings.TrimLeft(params, "/")
	status := 200

	var response dto.ResponseFromMeta
	var request dto.RequestToMeta

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("requestnya = ", request)

	if request.To == "invalid_phone_number" {
		response = dto.ResponseFromMeta{
			Error: &dto.ErrorRes{
				Message: "(#131030) Recipient phone number not in allowed list",
				Type:    "OAuthException",
				Code:    131030,
				ErrorData: map[string]interface{}{
					"messaging_product": "whatsapp",
					"details":           "Recipient phone number not in allowed list: Add recipient phone number…",
				},
				FbTraceID: "INVALID PHONE NUMBER",
			},
		}
		status = 400
	} else if request.Text.Body == "" {
		response = dto.ResponseFromMeta{
			Error: &dto.ErrorRes{
				Message: "(#100) The parameter text['body'] is required.",
				Type:    "OAuthException",
				Code:    100,

				FbTraceID: "BODY MSG WAS EMPTY",
			},
		}
		status = 400

	} else if request.Type != "text" {

		response = dto.ResponseFromMeta{
			Error: &dto.ErrorRes{
				Message: "(#100) Param type must be one of {AUDIO, CONTACTS, DOCUMENT, IMAGE, INTERACTIVE, LINK_PREVIEW, LOCATION, REACTION, STICKER, TEMPLATE, TEXT, VIDEO} - got \"not-text\".",
				Type:    "OAuthException",
				Code:    100,

				FbTraceID: "Type must be Text",
			},
		}
		status = 400
	} else if request.To == "" {
		response = dto.ResponseFromMeta{
			Error: &dto.ErrorRes{
				Message:   "The parameter to is required.",
				Type:      "OAuthException",
				Code:      100,
				FbTraceID: "NO \"TO\" FIELD",
			},
		}
		status = 400
	} else if token == "invalid_token" {
		response = dto.ResponseFromMeta{
			Error: &dto.ErrorRes{
				Message: "Invalid OAuth access token - Cannot parse access token",
				Type:    "OAuthException",
				Code:    190,
				ErrorData: map[string]interface{}{
					"messaging_product": "whatsapp",
					"details":           "Invalid OAuth access token - Cannot parse access token",
				},
				FbTraceID: "INVALID TOKEN",
			},
		}
		status = 401
	} else {
		response = dto.ResponseFromMeta{
			MessagingProduct: "whatsapp",
			Contacts: []struct {
				Input string `json:"input" bson:"input"`
				WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
			}{{
				Input: "valid_phone_number",
				WAID:  "valid_phone_number",
			}},
			Messages: []struct {
				ID string `json:"id"`
			}{{ID: "wamid.123"}},
		}

	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
	fmt.Println("Response from meta :", response)
	response = dto.ResponseFromMeta{}
}

func MetaServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MetaHandler(w, r)
	}))
	return server
}
