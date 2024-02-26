package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"transmitter-artemis/dto"
)

func MetaHandler(w http.ResponseWriter, r *http.Request) {

	// params := r.URL.Path[0:]
	// phoneNumberID := strings.TrimLeft(params, "/")

	fmt.Println("Masuk request meta")
	var request dto.RequestToMeta
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("requestnya = ", request)
	response := dto.ResponseFromMeta{
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
	// response := map[string]interface{}{
	// 	"id": p.RetailerID,
	// }
	// status := http.StatusOK
	// if p.RetailerID == "no_name" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    100,
	// 			"message": "(#100) The parameter name is required",
	// 		},
	// 	}
	// } else if p.RetailerID == "no_currency" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    100,
	// 			"message": "(#100) The parameter currency is required",
	// 		},
	// 	}
	// } else if p.RetailerID == "no_price" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    100,
	// 			"message": "(#100) The parameter price is required",
	// 		},
	// 	}
	// } else if p.RetailerID == "no_image_url" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    100,
	// 			"message": "(#10801) Either \"uploaded_image_id\" or \"image_url\" must be specified.",
	// 		},
	// 	}
	// } else if p.RetailerID == "no_description" || p.RetailerID == "no_url" || p.RetailerID == "no_brand" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    100,
	// 			"message": "Invalid parameter",
	// 		},
	// 	}
	// } else if p.RetailerID == "session_token_expired" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    500,
	// 			"message": "internal server error, please try again later",
	// 		},
	// 	}
	// } else if p.RetailerID == "retryable_error" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    123,
	// 			"message": "infinite retry",
	// 		},
	// 	}
	// } else if p.RetailerID == "max_retry" {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"code":    124,
	// 			"message": "internal server error, please try again later",
	// 		},
	// 	}
	// } else if p.RetailerID == "no_product_id" {
	// 	response = map[string]interface{}{}
	// }

	// if strings.Contains(productFBId, "error") {
	// 	response = map[string]interface{}{
	// 		"error": map[string]interface{}{
	// 			"message":       "Error from FB",
	// 			"type":          "GraphMethodException",
	// 			"code":          100,
	// 			"error_subcode": 33,
	// 			"fbtrace_id":    "Az8or2yhqkZfEZ-_4Qn_Bam",
	// 		},
	// 	}

	// 	status = http.StatusBadRequest
	// }

	jsonResponse, _ := json.Marshal(response)
	status := http.StatusOK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}

func MetaServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MetaHandler(w, r)
	}))
	return server
}
