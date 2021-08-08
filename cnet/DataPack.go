package cnet

import (
	"bytes"
	"cinx/ciface"
	"cinx/utils"
	"encoding/binary"
	"errors"
)

//只用于风暴和拆包的类，自身不存储数据
type DataPack struct{}

//获取实例
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头的长度
func (dataPack *DataPack) GetHeadLength() uint32 {
	//ID(uint32 4字节)+Length(uint32 4字节)
	return 8
}

//封包
func (dataPack *DataPack) Pack(msg ciface.IMessage) ([]byte, error) {
	//创建一个byte缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//写入ID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, err
	}
	//写入Length
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetLength()); err != nil {
		return nil, err
	}
	//写入Data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包（按照头部的固定长度解析消息头部）
func (dataPack *DataPack) UnPack(binaryDaya []byte) (ciface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryDaya)

	msg := &Message{}

	//读ID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	//读Length
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Length); err != nil {
		return nil, err
	}
	//判断包体是否超长
	if utils.GlobalCfg.MaxPackageSize > 0 && msg.Length > utils.GlobalCfg.MaxPackageSize {
		return nil, errors.New("[Cnet]Package too large")
	}

	return msg, nil
}
