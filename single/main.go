package main

import (
	engine "Crawler/single/engine"
	"fmt"
	"os"
	//parser "DistributedCrawler1/arser"
	// persist "DistributedCrawler1/Persist"
	// scheduler "DistributedCrawler1/Scheduler"
)

// go run main.go name starturl
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("argument is Invalid :%v\n", os.Args)
		return
	}
	e := engine.NewEngine(os.Args[1], os.Args[2])
	e.Run()
}

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
