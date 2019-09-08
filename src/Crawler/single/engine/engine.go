package engine

import (
	"Crawler/single/common"
	scheduler "Crawler/single/scheduler"
	spider "Crawler/single/spider"
	"sync"
)

type Engine struct {
	req       chan *common.Request
	resp      chan *common.ParseResult
	startUrl  string
	scheduler *scheduler.Scheduler
	spider    spider.Spider

	wg sync.WaitGroup
}

func NewEngine(name, startUrl string) *Engine {
	engine := &Engine{
		startUrl: startUrl,
		req:      make(chan *common.Request),
		resp:     make(chan *common.ParseResult),
		wg:       sync.WaitGroup{},
	}
	engine.scheduler = scheduler.NewScheduler(engine.req)
	engine.spider = spider.NewSpider(name, engine.req, engine.resp)
	return engine
}

func (e *Engine) Run() {
	//log.Error("e.req:%+v", e.req)
	//startUrl enqueue
	//e.spider.ConfigureMasterWorkerChan(e.req)
	e.wg.Add(1)
	go e.spider.Run()
	// e.req <- &common.Request{
	// 	Url:        e.startUrl,
	// 	ParserFunc: e.spider.Parse,
	// }
	e.scheduler.Submit(&common.Request{
		Url:        e.startUrl,
		ParserFunc: e.spider.Parse,
	})
	e.wg.Add(1)
	go func() {
		for {
			select {
			case result := <-e.resp:
				for _, request := range result.Requests {
					e.scheduler.Submit(request)
				}
			default:
				continue
			}
		}
	}()
	e.wg.Wait()
}

// type Scheduler interface {
// 	Submit(request *Request)
// 	ConfigureMasterWorkerChan(chan *Request)
// 	// WorkerReady(chan Request)
// 	// Run()
// }

// type ConcurrentEngine struct {
// 	Scheduler   Scheduler //Sheduler
// 	WorkerCount int       //worker的数量
// 	ItemChan    chan interface{}
// }

// func (e *ConcurrentEngine) Run(seeds ...*Request) {

// 	//worker公用一个in，out
// 	in := make(chan *Request)
// 	out := make(chan *ParseResult)

// 	e.Scheduler.ConfigureMasterWorkerChan(in)

// 	for i := 0; i < e.WorkerCount; i++ {
// 		createWorker(in, out) //创建worker
// 	}

// 	//参数seeds的request，要分配任务
// 	for _, r := range seeds {
// 		e.Scheduler.Submit(r)
// 	}

// 	//从out中获取result，对于item就打印即可，对于request，就继续分配
// 	for {
// 		result := <-out
// 		for _, item := range result.Items {
// 			e.ItemChan <- item
// 		}

// 		for _, request := range result.Requests {
// 			e.Scheduler.Submit(request)
// 		}
// 	}
// }

// func createWorker(in chan *Request, out chan *ParseResult) {
// 	go func() {
// 		for {
// 			request := <-in
// 			result, err := Worker(request)
// 			if err != nil {
// 				continue
// 			}
// 			out <- &result
// 		}
// 	}()
// }

// func Worker(r *Request) (ParseResult, error) {
// 	log.Printf("Fetching %s", r.Url)
// 	body, err := fetcher.Fetch(r.Url)
// 	if err != nil {
// 		log.Printf("Fetcher: error fetching url %s %v", r.Url, err)
// 		return ParseResult{}, fmt.Errorf("fetching url error :%v", err)
// 	}
// 	return r.ParserFunc(body), nil
// }
