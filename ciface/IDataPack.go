package ciface

type IDataPack interface {
	//获取包头的长度
	GetHeadLength() uint32
	//封包
	Pack(IMessage) ([]byte, error)
	//拆包
	UnPack([]byte) (IMessage, error)
}
