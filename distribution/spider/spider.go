package spider

import (
	"Crawler/distribution/common"

	"github.com/garyburd/redigo/redis"
	"github.com/liangdas/mqant/log"
)

// Spider download html and parser it
type Spider interface {
	Download(url string) ([]byte, error)
	Parse([]byte) (common.ParseResult, error)
	Run()
}

func NewSpider(name string, conn redis.Conn, workerNum int) Spider {
	switch name {
	case "zhenai":
		return NewZhenAiSpider(name, conn, workerNum)
	case "music":
		return NewMusicSpider(name, conn, workerNum)
	default:
		log.Error("Unsupport type :%v", name)
		return nil
	}
}
