package itempipeline

import "sync"

// MusicPipeline implements Pipeline interface
type MusicPipeline struct {
	itemChan chan interface{}
	wg       sync.WaitGroup
}

// HandleItem implements Pipeline.HandleItem interface
func (music *MusicPipeline) HandleItem(item interface{}) error {
	return nil
}

// Run implements Pipeline.Run interface
func (music *MusicPipeline) Run() {
}

func NewMusicPipeline(itemChan chan interface{}) Pipeline {
	return &MusicPipeline{
		itemChan: itemChan,
		wg:       sync.WaitGroup{},
	}
}
