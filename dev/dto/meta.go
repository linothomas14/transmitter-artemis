package dto

type ResponseFromMeta struct {
	MessagingProduct string `json:"messaging_product,omitempty" bson:"messaging_product,omitempty"`
	Contacts         []struct {
		Input string `json:"input" bson:"input"`
		WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
	} `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages,omitempty" bson:"messages,omitempty"`
	Error *ErrorRes `json:"error,omitempty" bson:"error,omitempty"`
}

type ErrorRes struct {
	Message   string                 `json:"message,omitempty" bson:"message,omitempty"`
	Type      string                 `json:"type,omitempty" bson:"type,omitempty"`
	Code      int                    `json:"code,omitempty" bson:"code,omitempty"`
	ErrorData map[string]interface{} `json:"error_data,omitempty" bson:"error_data,omitempty"`
	FbTraceID string                 `json:"fbtrace_id,omitempty" bson:"fb_trace_id,omitempty"`
}

type RequestToMeta struct {
	MessagingProduct string `json:"messaging_product" bson:"messaging_product"`
	RecipientType    string `json:"recipient_type" bson:"recipient_type"`
	To               string `json:"to" bson:"to"`
	Type             string `json:"type" bson:"type"`
	Text             struct {
		PreviewURL bool   `json:"preview_url" bson:"preview_url"`
		Body       string `json:"body" bson:"body"`
	} `json:"text"`
}
