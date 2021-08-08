package ciface

//用于替代[]byte，包装传输的数据
type IMessage interface {
	GetID() uint32
	GetLength() uint32
	GetData() []byte

	SetID(uint32)
	SetLength(uint32)
	SetData([]byte)
}
