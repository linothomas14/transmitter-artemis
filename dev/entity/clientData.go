package entity

type ClientData struct {
	ClientName    string `json:"client_name" bson:"client_name"`
	Token         string `json:"token" bson:"token"`
	PhoneNumberID string `json:"phone_number_id" bson:"phone_number_id"`
	WAHost        string `json:"wa_host" bson:"wa_host"`
}
