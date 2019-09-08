package main

import (
	"Crawler/distribution/common"
	"Crawler/distribution/distribute"
	"encoding/json"

	"fmt"
	"os"

	"github.com/liangdas/mqant/log"
)

// 1) Master (e.g., go run main.go master "http://www.zhenai.com/zhenghun")
// 2) Slave (e.g., go run main.go slave "zhenai")
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("argument is Invalid :%v\n", os.Args)
		return
	}
	switch os.Args[1] {
	case "master":
		startReq, err := json.Marshal(common.Request{
			Url:  os.Args[2],
			Flag: 1,
		})
		if err != nil {
			log.Error("err:%v", err)
		}
		distribute.NewMaster().Run(startReq)
	case "slave":
		distribute.NewSlave(os.Args[2]).Run()
	}
}

// func main() {
// 	//InitRedis("127.0.0.1:6379")
// 	wg := sync.WaitGroup{}
// 	wg.Add(2)
// 	go Custom()
// 	time.Sleep(5 * time.Second)
// 	go Produce([]byte("hello"))
// 	wg.Wait()
// }
// func Custom() {
// 	c, err := redis.Dial("tcp", "127.0.0.1:6379")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	defer c.Close()

// 	for {
// 		vals, err := redis.Values(c.Do("brpop", "REQUEST_KEY", 0))
// 		if err != nil {
// 			fmt.Errorf("brpop error:%v", err)
// 			time.Sleep(3 * time.Second)
// 			continue
// 		}

// 		for i, v := range vals {
// 			//if i != 0 {
// 			fmt.Printf("i:%d-->v:%s", i, string(v.([]byte)))
// 			//}
// 		}
// 	}
// }

// func Produce(szBytes []byte) (err error) {
// 	pConn := distribute.GetRedisConn()
// 	if pConn.Err() != nil {
// 		fmt.Println(pConn.Err().Error())
// 		return
// 	}
// 	defer pConn.Close()

// 	if _, err = pConn.Do("lpush", "REQUEST_KEY", szBytes); err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	return
// }

// func main() {
// 	url := "http://www.zhenai.com/zhenghun"
// 	e := engine.NewEngine(url)
// 	e.Run()
// }

// func main() {
// url := "http://album.zhenai.com/u/109484507"
// resp, err := fetch.Fetch(url)
// if err != nil {
// 	fmt.Println(err)
// }
// err = ioutil.WriteFile("test.txt", resp, 0644)
// if err != nil {
// 	log.Println(err)
// }
// 	payload, err := ioutil.ReadFile("test.txt")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	parser.ParseProfile(payload)
// }
