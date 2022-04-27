package yx

const (
	BP_MIN_BUFF_SIZE = 32
	BP_MAX_BUFF_SIZE = 1024
	BP_QUEUE_STEP    = 32
)

type BuffPool struct {
	buffQueues  []Queue
	maxReuseCnt int
}

func NewBuffPool(maxReuseCnt uint64) *BuffPool {
	queueCnt := (BP_MAX_BUFF_SIZE-BP_MIN_BUFF_SIZE)/BP_QUEUE_STEP + 1
	p := &BuffPool{
		buffQueues:  make([]Queue, queueCnt),
		maxReuseCnt: int(maxReuseCnt),
	}

	for i := 0; i < queueCnt; i++ {
		p.buffQueues[i] = NewSyncLimitQueue(maxReuseCnt)
	}

	return p
}

func (p *BuffPool) CreateBuff(size uint16) []byte {
	if (size > BP_MAX_BUFF_SIZE) || (p.maxReuseCnt == 0) {
		return make([]byte, size)
	}

	idx := 0
	if size > BP_MIN_BUFF_SIZE {
		idx = (int(size) - BP_MIN_BUFF_SIZE) / BP_QUEUE_STEP
		if (size-BP_MIN_BUFF_SIZE)%BP_QUEUE_STEP != 0 {
			idx++
		}
	}

	item, err := p.buffQueues[idx].Dequeue()
	if err != nil {
		fixSize := BP_MIN_BUFF_SIZE + idx*BP_QUEUE_STEP
		return make([]byte, fixSize)
	}

	return item.([]byte)
}

func (p *BuffPool) ReuseBuff(b []byte) {
	if p.maxReuseCnt == 0 {
		return
	}

	size := len(b)
	if size < BP_MIN_BUFF_SIZE || size > BP_MAX_BUFF_SIZE {
		return
	}

	if (size-BP_MIN_BUFF_SIZE)%BP_QUEUE_STEP != 0 {
		return
	}

	idx := (size - BP_MIN_BUFF_SIZE) / BP_QUEUE_STEP
	p.buffQueues[idx].Enqueue(b)
}
