package spider

import (
	"Crawler/distribution/common"
	itempipe "Crawler/distribution/itempipeline"
	"sync"

	"github.com/garyburd/redigo/redis"
)

// MusicSpider implements Spider interface
type MusicSpider struct {
	pipe      itempipe.Pipeline
	pipeChan  chan interface{}
	workerNum int
	conn      redis.Conn
	wg        sync.WaitGroup
	sync.RWMutex
}

func NewMusicSpider(name string, conn redis.Conn, workerNum int) Spider {
	pipeChan := make(chan interface{})
	pipe := itempipe.NewPipeline(name, pipeChan)

	return &MusicSpider{
		pipe:      pipe,
		pipeChan:  pipeChan,
		workerNum: workerNum,
		conn:      conn,
		wg:        sync.WaitGroup{},
	}
}

// Run implements Spider.Run interface
func (music *MusicSpider) Run() {
}

// Download implements Spider.Download interface
func (music *MusicSpider) Download(url string) ([]byte, error) {
	return []byte{}, nil
}

// Parse implements Spider.Parse interface
func (music *MusicSpider) Parse(contents []byte) (common.ParseResult, error) {
	return common.ParseResult{}, nil
}
