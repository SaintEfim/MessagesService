package entity

type OperationType string

const (
	OperationInit        OperationType = "init"
	OperationSendMessage               = "send_message"
)

type TCPRequest struct {
	UserCredential UserCredential `json:"user_credential"`
	Message        string         `json:"message"`
	Operation      OperationType  `json:"operation"`
}
