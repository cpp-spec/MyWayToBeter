package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// 定义了对象实例，但没有填入具体成员？
type DataPack struct {
}

// 初始化实例
// 这里很有意思	DataPack{}功能为初始化了一个0值实例
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度 DataLen+Id即4+4=8bite
func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

// 打包函数，传入IMessage对象，返回...
// 制造一个缓冲，分别按顺序写入数据
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	databuff := bytes.NewBuffer([]byte{}) //创建一个存放bytes字节的缓冲
	//向缓冲写入datalen,后面是类似的目的
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetDateLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return databuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	databuff := bytes.NewReader(binaryData) //这里和Pack函数对比，前者创建写，后者自然是读
	//下面是解压的过程
	msg := &Message{}
	//第三个参数是一个地址，这样可以将对应的数据读到合适的地址去，这很重要，因为整个流程是约定好的
	if err := binary.Read(databuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(databuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received") //go语言约定俗成，错误信息用小写开头
	}
	return msg, nil
}
