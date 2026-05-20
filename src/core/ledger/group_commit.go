package ledger

import (
	"sync"
	"time"
)

type WriteRequest struct {
	Data  []byte
	Block *Block
	Done  chan error
}

type GroupCommit struct {
	queue  chan *WriteRequest
	writer interface {
		Write(data []byte) (int, error)
		Sync() error
	}
	linger time.Duration
	stop   chan struct{}
	wg     sync.WaitGroup
}

func NewGroupCommit(writer interface {
	Write(data []byte) (int, error)
	Sync() error
}, linger time.Duration) *GroupCommit {
	gc := &GroupCommit{
		queue:  make(chan *WriteRequest, 1000),
		writer: writer,
		linger: linger,
		stop:   make(chan struct{}),
	}
	gc.wg.Add(1)
	go gc.run()
	return gc
}

func (gc *GroupCommit) Submit(data []byte) chan error {
	done := make(chan error, 1)
	gc.queue <- &WriteRequest{Data: data, Done: done}
	return done
}

func (gc *GroupCommit) SubmitBlock(b *Block) chan error {
	done := make(chan error, 1)
	gc.queue <- &WriteRequest{Block: b, Done: done}
	return done
}

func (gc *GroupCommit) run() {
	defer gc.wg.Done()
	for {
		select {
		case req := <-gc.queue:
			batch := []*WriteRequest{req}
			
			timeout := time.NewTimer(gc.linger)
		collect:
			for {
				select {
				case next := <-gc.queue:
					batch = append(batch, next)
					if len(batch) >= 100 {
						break collect
					}
				case <-timeout.C:
					break collect
				case <-gc.stop:
					break collect
				}
			}
			timeout.Stop()
			
			gc.processBatch(batch)
			
		case <-gc.stop:
			for len(gc.queue) > 0 {
				req := <-gc.queue
				gc.processBatch([]*WriteRequest{req})
			}
			return
		}
	}
}

func (gc *GroupCommit) processBatch(batch []*WriteRequest) {
	if len(batch) == 0 {
		return
	}

	var lastErr error
	for _, req := range batch {
		if req.Block != nil {
			err := req.Block.Encode(gc.writer)
			if err != nil {
				lastErr = err
				break
			}
		} else {
			_, err := gc.writer.Write(req.Data)
			if err != nil {
				lastErr = err
				break
			}
		}
	}
	
	if lastErr == nil {
		lastErr = gc.writer.Sync()
	}

	for _, req := range batch {
		req.Done <- lastErr
	}
}

func (gc *GroupCommit) Stop() {
	close(gc.stop)
	gc.wg.Wait()
}
