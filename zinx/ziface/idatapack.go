// 拆包封包抽象层接口
package ziface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (IMessage, error)   //解包方法
}
