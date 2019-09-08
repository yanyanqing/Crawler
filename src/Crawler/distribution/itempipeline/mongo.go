package itempipeline

import (
	"time"

	"gopkg.in/mgo.v2"
)

var GlobalMgoSession *mgo.Session

func init() {
	//globalMgoSession, err := mgo.DialWithTimeout("mongodb://admin:123456@localhost:27017/admin", 10*time.Second)
	globalMgoSession, err := mgo.DialWithTimeout("127.0.0.1:27017", 10*time.Second)
	if err != nil {
		panic(err)
	}
	GlobalMgoSession = globalMgoSession
	GlobalMgoSession.SetMode(mgo.Monotonic, true)
	//default is 4096
	GlobalMgoSession.SetPoolLimit(300)
}

func CloneSession() *mgo.Session {
	return GlobalMgoSession.Clone()
}
