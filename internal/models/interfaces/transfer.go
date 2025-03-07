package interfaces

type Transfer interface {
	TransferData(data interface{}) error
}
