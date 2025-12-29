package junbanmachi // 順番待ち - じゅんばんまち - Queuing

type Packet struct {
	ID       uint32
	SourceIP string
	Domain   string
	URI      string
	Method   string
	//Service string
	Headers  []*PushRequest_HeaderSchema
}
