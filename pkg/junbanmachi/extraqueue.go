package junbanmachi // 順番待ち - じゅんばんまち - Queuing

import (
	"reflect"
	"sync"

	"github.com/bonavadeur/katyusha/pkg/bonalib"
	"github.com/bonavadeur/katyusha/pkg/hashi"
)

type ExtraQueue struct {
	Queue           []*Packet
	//Queues map[string][]*Packet
	Next            chan bool
	NextQueueSize   int
	NextQueueLength int
	popLock         sync.Mutex
	pushBridge      *hashi.Hashi
	popBridge       *hashi.Hashi
	MsgIDLock       sync.Mutex
	MsgIDCount      uint32
	sortLock        sync.Mutex
	queueLock       sync.Mutex
}

func NewExtraQueue() *ExtraQueue {
	containerConcurrency := bonalib.Cm2Int("katyusha-junbanmachi-concurrent-request")
	katyushaThreads := bonalib.Cm2Int("katyusha-threads")

	newExtraQueue := &ExtraQueue{
		Queue:           make([]*Packet, 0),
		//Queues: make(map[string][]*Packet),
		Next:            make(chan bool, katyushaThreads),
		NextQueueSize:   containerConcurrency,
		NextQueueLength: 0,
		popLock:         sync.Mutex{},
		MsgIDLock:       sync.Mutex{},
		MsgIDCount:      0,
		sortLock:        sync.Mutex{},
		queueLock:       sync.Mutex{},
	}
	newExtraQueue.pushBridge = hashi.NewHashi(
		"PushBridge",
		hashi.HASHI_TYPE_SERVER,
		BASE_PATH+"/push-bridge",
		katyushaThreads,
		reflect.TypeOf(PushRequest{}),
		reflect.TypeOf(PushResponse{}),
		newExtraQueue.PushResponseAdapter,
	)
	newExtraQueue.popBridge = hashi.NewHashi(
		"PopBridge",
		hashi.HASHI_TYPE_SERVER,
		BASE_PATH+"/pop-bridge",
		katyushaThreads,
		reflect.TypeOf(PopRequest{}),
		reflect.TypeOf(PopResponse{}),
		newExtraQueue.PopResponseAdapter,
	)

	return newExtraQueue
}

func (q *ExtraQueue) PushResponseAdapter(params ...interface{}) (interface{}, error) {
	pushRequest := params[0].(*PushRequest)
	q.Push(pushRequest)
	return &PushResponse{SymbolizeResponse: Status_Success}, nil
}

func (q *ExtraQueue) PopResponseAdapter(params ...interface{}) (interface{}, error) {
	// _ = params[0].(*PopRequest)
	popPacket := q.Pop()
	return packet2PopResponse(popPacket), nil
}

func (q *ExtraQueue) Push(pushPacket *PushRequest) {
	q.queueLock.Lock()
	q.sort(pushRequest2Packet(pushPacket))
	q.queueLock.Unlock()
	q.Next <- true
}

//pop request

func (q *ExtraQueue) Pop() *Packet {
	q.popLock.Lock()
	defer q.popLock.Unlock()

	<-q.Next // hangout until len of Queue > 0
	q.queueLock.Lock()
	popPacket := q.Queue[len(q.Queue)-1]
	q.HeaderModifier(popPacket)
	q.Queue = q.Queue[:len(q.Queue)-1]
	q.queueLock.Unlock()

	return popPacket
}




// func (q *ExtraQueue) Pop(service string) *Packet {
// 	q.popLock.Lock()
// 	defer q.popLock.Unlock()

// 	<-q.Next // chờ có request bất kỳ

// 	q.queueLock.Lock()
// 	queue := q.Queues[service]
// 	if len(queue) == 0 {
// 		q.queueLock.Unlock()
// 		return nil
// 	}

// 	popPacket := queue[len(queue)-1]
// 	q.HeaderModifier(popPacket)
// 	q.Queues[service] = queue[:len(queue)-1]
// 	q.queueLock.Unlock()

// 	return popPacket
// }





// func (q *ExtraQueue) Pop() *Packet {
// 	q.popLock.Lock()
// 	defer q.popLock.Unlock()

// 	// ====== BATCH GATE (sửa ở đây) ======
// 	for {const MIN_QUEUE_BEFORE_POP = 3
// 		q.queueLock.Lock()
// 		queueLen := len(q.Queue)
// 		q.queueLock.Unlock()

// 		if queueLen >= MIN_QUEUE_BEFORE_POP {
// 			break
// 		}

// 		<-q.Next // chờ thêm request
// 	}
// 	// ====== HẾT PHẦN SỬA ======

// 	q.queueLock.Lock()
// 	popPacket := q.Queue[len(q.Queue)-1]
// 	q.HeaderModifier(popPacket)
// 	q.Queue = q.Queue[:len(q.Queue)-1]
// 	q.queueLock.Unlock()

// 	return popPacket
// }


// custom this
func (q *ExtraQueue) sort(p *Packet) {
	q.sortLock.Lock()
	defer q.sortLock.Unlock()

	q.SortAlgorithm(p)
}
