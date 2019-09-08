package itempipeline

import "github.com/liangdas/mqant/log"

// Pipeline handles item and save it to mongo
type Pipeline interface {
	HandleItem(item interface{}) error
	Run()
}

func NewPipeline(name string, pipeChan chan interface{}) Pipeline {
	switch name {
	case "zhenai":
		return NewZhenAiPipeline(pipeChan)
	case "music":
		return NewMusicPipeline(pipeChan)
	default:
		log.Error("Unsupport type :%v", name)
		return nil
	}
}
