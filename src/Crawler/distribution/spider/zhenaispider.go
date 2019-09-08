package spider

import (
	"Crawler/distribution/common"
	"encoding/json"

	//"Crawler/distribution/distribute"
	//"Crawler/distribution/distribute"
	itempipe "Crawler/distribution/itempipeline"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
	"github.com/liangdas/mqant/log"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" [^>]*>([^<]+)</a>`
	CityRe     = `<a href="(http://album.zhenai.com/u/[\d]+)" [^>]*>([^<]+)</a>`
	//CityReNext = `<a href="(http://www.zhenai.com/zhenghun/aba/[\d]+)">下一页`
)

var (
	reProfile = regexp.MustCompile(`<script>window.__INITIAL_STATE__=(.+);\(function`)
	reID      = regexp.MustCompile(`ID：(\d+)</div>`)
	reName    = regexp.MustCompile(`<h1 class="nickName" [^>]*>(.+)</h1>`)
)

// ZhenAiSpider implements Spider interface
type ZhenAiSpider struct {
	pipe      itempipe.Pipeline
	pipeChan  chan interface{}
	workerNum int
	conn      redis.Conn
	wg        sync.WaitGroup
	sync.RWMutex
}

func NewZhenAiSpider(name string, conn redis.Conn, workerNum int) Spider {
	pipeChan := make(chan interface{})
	pipe := itempipe.NewPipeline(name, pipeChan)

	return &ZhenAiSpider{
		pipe:      pipe,
		pipeChan:  pipeChan,
		workerNum: workerNum,
		conn:      conn,
		wg:        sync.WaitGroup{},
	}
}

// Run implements Spider.Run interface
func (zaSpider *ZhenAiSpider) Run() {
	//zaSpider.wg.Add(zaSpider.workerNum + 1)
	zaSpider.wg.Add(2)
	go zaSpider.pipe.Run()
	//	for i := 0; i < zaSpider.workerNum; i++ {
	go func() {
		for {
			// 多个 gorutine 同时对一张表执行操作会出现 use of closed network connection
			// https://blog.csdn.net/chenbaoke/article/details/39899177
			//zaSpider.Lock()
			vals, err := redis.Values(zaSpider.conn.Do("brpop", "REQUEST_KEY", 0))
			if err != nil {
				log.Error("slave spider brpop error:%v conn:%v", err, zaSpider.conn)
				time.Sleep(3 * time.Second)
				continue
			}

			for i, v := range vals {
				if i != 0 {
					var req common.Request
					log.Error("v:%v", string(v.([]byte)))
					err = json.Unmarshal(v.([]byte), &req)
					if err != nil {
						log.Error("slave spider json.Unmarshal error:%v", err)
						continue
					}
					//fmt.Printf("i:%d-->v:%v", i, v.(common.Request))
					//zaSpider.pipeChan <- req.Url
					body, err := zaSpider.Download(req.Url)
					if err != nil {
						log.Error("slave spider Download: error Download url %s %v", req.Url, err)
						continue
					}
					resp, err := zaSpider.bodyParse(req.Flag, body)
					if err != nil {
						log.Error("slave spider ParserFunc error:%v", err)
						continue
					}
					respJson, err := json.Marshal(resp)
					if err != nil {
						log.Error("slave spider json Marshal resp error:%v", err)
						continue
					}

					_, err = zaSpider.conn.Do("lpush", "RESPONSE_KEY", respJson)
					if err != nil {
						log.Error("slave spider lpush error:%v", err)
					}
					log.Info("Download %s body:%v resp:%s", req.Url, len(body), string(respJson))
					for item := range resp.Items {
						log.Info("Item:%v", item)
						zaSpider.pipeChan <- item
					}
				}
			}
			//	zaSpider.Unlock()
		}
	}()
	//	}

	zaSpider.wg.Wait()
}

// Download implements Spider.Download interface
func (zaSpider *ZhenAiSpider) Download(url string) ([]byte, error) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")

	resp, _ := http.DefaultClient.Do(request)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("resp:", resp)
		return nil, fmt.Errorf("error:status code:%d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

// Parse implements Spider.Parse interface
func (zaSpider *ZhenAiSpider) Parse(contents []byte) (common.ParseResult, error) {
	re := regexp.MustCompile(cityListRe)
	all := re.FindAllSubmatch(contents, -1)
	i := 0
	result := common.ParseResult{}
	for _, c := range all {
		//result.Items = append(result.Items, string(c[2])) //cityName
		result.Requests = append(result.Requests, &common.Request{
			Url:  string(c[1]),
			Flag: 2,
			//ParserFunc: zaSpider.cityParse,
		})
		if i == 10 {
			break
		}
		i++
	}

	return result, nil
}

func (zaSpider *ZhenAiSpider) bodyParse(flag int, body []byte) (common.ParseResult, error) {
	switch flag {
	case 1:
		return zaSpider.Parse(body)
	case 2:
		return zaSpider.cityParse(body)
	case 3:
		return zaSpider.profileParse(body)
	default:
		return common.ParseResult{}, nil
	}

}
func (zaSpider *ZhenAiSpider) cityParse(contents []byte) (common.ParseResult, error) {
	re := regexp.MustCompile(CityRe)
	all := re.FindAllSubmatch(contents, -1)
	result := common.ParseResult{}
	for _, c := range all {
		//result.Items = append(result.Items, string(c[2])) //username
		result.Requests = append(result.Requests, &common.Request{
			Url:  string(c[1]),
			Flag: 3,
			//ParserFunc: zaSpider.profileParse,
		})
	}

	return result, nil
}

func (zaSpider *ZhenAiSpider) profileParse(contents []byte) (common.ParseResult, error) {
	var ID string
	var name string
	match := reID.FindSubmatch(contents)
	if len(match) >= 2 {
		ID = string(match[1])
	}
	match = reName.FindSubmatch(contents)
	if len(match) >= 2 {
		name = string(match[1])
	}
	match = reProfile.FindSubmatch(contents)
	if len(match) >= 2 {
		json := match[1]
		profile := parseJson(json)
		if profile != nil {
			profile.ID = ID
			profile.Name = name
			zaSpider.pipeChan <- profile
		}
	}
	return common.ParseResult{}, nil
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Error("transfer encoding error:%v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

func parseJson(json []byte) *itempipe.UserItem {
	res, err := simplejson.NewJson(json)
	if err != nil {
		log.Error("parseJson error:%v", err)
		return nil
	}
	infos, err := res.Get("objectInfo").Get("basicInfo").Array()

	profile := &itempipe.UserItem{}
	for k, v := range infos {
		if e, ok := v.(string); ok {
			switch k {
			case 0:
				profile.Marriage = e
			case 1:
				profile.Age = e
			case 2:
				profile.Xingzuo = e
			case 3:
				profile.Height = e
			case 4:
				profile.Weight = e
			case 6:
				profile.Income = e
			case 7:
				profile.Occupation = e
			case 8:
				profile.Education = e
			}
		}

	}

	return profile
}
