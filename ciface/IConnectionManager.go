package ciface

type IConnectionManager interface {
	//添加连接
	Add(IConnection)
	//删除
	Remove(IConnection)
	//根据ID获取连接
	Get(uint32) (IConnection, error)
	//获取连接总数
	Len() int
	//清除并终止所有连接
	Clear()
}
