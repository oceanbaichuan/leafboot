package msg

type GameRegistReq struct {
	Addr     string
	NodeID   string //服务器节点描述
	NodeName string //服务器类型 例如:"login", "xzmj"
	IsGray   int8   //是否灰度
}

type GameRegistRes struct {
	HashCode []uint32
}

type GameFlagNotice struct {
	IsClosed bool
}
