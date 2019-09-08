package distribute

import (
	"Crawler/distribution/common"
	"encoding/json"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/liangdas/mqant/log"
	"github.com/willf/bloom"
)

type Master struct {
	conn      redis.Conn
	filter    *bloom.BloomFilter
	startUrl  string
	workerNum int
	wg        sync.WaitGroup
}

func NewMaster() *Master {
	return &Master{
		conn:      GetRedisConn(),
		filter:    bloom.New(100000, 5),
		workerNum: 10,
		wg:        sync.WaitGroup{},
	}
}

func (m *Master) Run(startReq []byte) {
	_, err := m.conn.Do("lpush", "REQUEST_KEY", startReq)
	if err != nil {
		log.Error("master lpush start request error:%v", err)
		return
	}
	m.wg.Add(1)
	// m.wg.Add(m.workerNum)
	//	for i := 0; i < m.workerNum; i++ {
	go func() {
		for {
			// 多个 gorutine 同时对一张表执行操作会出现 use of closed network connection
			// https://blog.csdn.net/chenbaoke/article/details/39899177
			vals, err := redis.Values(m.conn.Do("brpop", "RESPONSE_KEY", 0))
			if err != nil {
				log.Error("master brpop RESPONSE_KEY error:%v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			for i, v := range vals {
				if i != 0 {
					var resp common.ParseResult
					err = json.Unmarshal(v.([]byte), &resp)
					if err != nil {
						log.Error("master json.Unmarshal error:%v", err)
						continue
					}

					for _, request := range resp.Requests {
						if !m.filter.TestString(request.Url) {
							reqJson, err := json.Marshal(request)
							if err != nil {
								log.Error("master json.Marshal error:%v", err)
								continue
							}
							_, err = m.conn.Do("lpush", "REQUEST_KEY", reqJson)
							if err != nil {
								log.Error("master lpush REQUEST_KEY error:%v", err)
							}

							m.filter.AddString(request.Url)
							log.Info("resp:%s", i, request.Url)
						}

					}
				}
			}

		}
	}()
	//	}

	m.wg.Wait()
}
