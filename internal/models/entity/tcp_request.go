package entity

type TCPRequest struct {
	UserCredential UserCredential `json:"user_credential" binding:"required"`
	Message        string         `json:"message" binding:"required"`
}
