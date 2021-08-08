package cnet

//用于替代[]byte，包装传输的数据
type Message struct {
	ID     uint32
	Length uint32
	Data   []byte
}

func NewMessage(ID uint32, data []byte) *Message {
	return &Message{
		ID:     ID,
		Length: uint32(len(data)),
		Data:   data,
	}
}

func (msg *Message) GetID() uint32 {
	return msg.ID
}
func (msg *Message) GetLength() uint32 {
	return msg.Length
}
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetID(ID uint32) {
	msg.ID = ID
}
func (msg *Message) SetLength(Length uint32) {
	msg.Length = Length
}
func (msg *Message) SetData(Data []byte) {
	msg.Data = Data
}
