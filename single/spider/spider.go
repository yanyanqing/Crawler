package spider

import (
	"Crawler/single/common"

	"github.com/liangdas/mqant/log"
)

// Spider download html and parser it
type Spider interface {
	Download(url string) ([]byte, error)
	Parse([]byte) (common.ParseResult, error)
	Run()
}

func NewSpider(name string, in chan *common.Request, out chan *common.ParseResult) Spider {
	switch name {
	case "zhenai":
		return NewZhenAiSpider(name, in, out)
	case "music":
		return NewMusicSpider(name, in, out)
	default:
		log.Error("Unsupport type :%v", name)
		return nil
	}
}
