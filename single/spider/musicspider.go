package spider

import (
	"Crawler/single/common"
	itempipe "Crawler/single/itempipeline"
)

// MusicSpider implements Spider interface
type MusicSpider struct {
	in        chan *common.Request
	out       chan *common.ParseResult
	pipe      itempipe.Pipeline
	pipeChan  chan interface{}
	workerNum int
}

func NewMusicSpider(name string, in chan *common.Request, out chan *common.ParseResult) Spider {
	pipeChan := make(chan interface{})
	pipe := itempipe.NewPipeline(name, pipeChan)

	return &MusicSpider{
		in:        in,
		out:       out,
		pipe:      pipe,
		pipeChan:  pipeChan,
		workerNum: 10,
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
