package scheduler

import (
	"Crawler/single/common"

	"github.com/willf/bloom"
)

// Scheduler filter request and push request to spider downloader
type Scheduler struct {
	reqChan chan *common.Request
	filter  *bloom.BloomFilter
}

func NewScheduler(reqchan chan *common.Request) *Scheduler {
	return &Scheduler{
		reqChan: reqchan,
		filter:  bloom.New(100000, 5),
	}
}

func (s *Scheduler) Submit(request *common.Request) {
	go func() {
		if !s.filter.TestString(request.Url) {
			s.reqChan <- request
			s.filter.AddString(request.Url)
		}
	}()
}
