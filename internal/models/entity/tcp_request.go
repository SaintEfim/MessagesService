package entity

type TCPRequest struct {
	UserCredential UserCredential `json:"user_credential"`
	Message        string         `json:"message"`
}
