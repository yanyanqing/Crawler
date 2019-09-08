package distribute

import (
	"Crawler/distribution/spider"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type Slave struct {
	spider    spider.Spider
	workerNum int
	conn      redis.Conn
	wg        sync.WaitGroup
}

func NewSlave(name string) *Slave {
	slave := &Slave{
		workerNum: 10,
		conn:      GetRedisConn(),
		wg:        sync.WaitGroup{},
	}

	slave.spider = spider.NewSpider(name, slave.conn, slave.workerNum)
	return slave
}

func (s *Slave) Run() {
	s.wg.Add(1)
	go s.spider.Run()
	s.wg.Wait()
}
