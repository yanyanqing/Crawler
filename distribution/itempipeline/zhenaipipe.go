package itempipeline

import (
	"sync"

	"github.com/juju/errors"
	"github.com/liangdas/mqant/log"
)

// ZhenAiPipeline implements Pipeline interface
type ZhenAiPipeline struct {
	itemChan chan interface{}
	wg       sync.WaitGroup
}

// HandleItem implements Pipeline.HandleItem interface
func (zaPipe *ZhenAiPipeline) HandleItem(item interface{}) error {
	log.Error("HandleItem:%v", item)
	session := CloneSession()
	defer session.Close()
	c := session.DB("scrawler").C("user")
	// var user Profile
	// user.ID = item.(string)
	// err := c.Insert(user)
	err := c.Insert(item)
	if err != nil {
		log.Error("err:%v", err)
		return errors.Trace(err)
	}

	var users []UserItem
	err = c.Find(nil).All(&users)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Run implements Pipeline.Run interface
func (zaPipe *ZhenAiPipeline) Run() {
	for {
		select {
		case req := <-zaPipe.itemChan:
			zaPipe.HandleItem(req)
		default:
			continue
		}
	}

}

func NewZhenAiPipeline(itemChan chan interface{}) Pipeline {
	return &ZhenAiPipeline{
		itemChan: itemChan,
		wg:       sync.WaitGroup{},
	}
}
