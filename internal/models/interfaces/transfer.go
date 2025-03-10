package interfaces

type Transfer interface {
	TransferData(data interface{}) error
	TransferText(data string) error
}
