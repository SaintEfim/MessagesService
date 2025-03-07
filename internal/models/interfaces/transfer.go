package interfaces

type Transfer interface {
	TransferData(data interface{}) error
	TransferDataText(text string) error
}
