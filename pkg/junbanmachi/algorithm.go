package junbanmachi // 順番待ち - じゅんばんまち - Queuing

import "github.com/bonavadeur/katyusha/pkg/bonalib"

func (q *ExtraQueue) SortAlgorithm(p *Packet) {
	bonalib.Info("SortAlgorithm", "Packet", p)
	// example of adding header
	p.Headers = append(p.Headers, &PushRequest_HeaderSchema{
		Field: "Katyusha-J-Field-1",
		Value: "Katyusha-J-Value-1",
	})

	q.Queue = append([]*Packet{p}, q.Queue...)
	//q.Queues[p.Service] = append([]*Packet{p}, q.Queues[p.Service]...)
}
