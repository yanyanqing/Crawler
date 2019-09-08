package spider

import (
	"Crawler/single/common"
	itempipe "Crawler/single/itempipeline"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/bitly/go-simplejson"
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
	in        chan *common.Request
	out       chan *common.ParseResult
	pipe      itempipe.Pipeline
	pipeChan  chan interface{}
	workerNum int
}

func NewZhenAiSpider(name string, in chan *common.Request, out chan *common.ParseResult) Spider {
	pipeChan := make(chan interface{})
	pipe := itempipe.NewPipeline(name, pipeChan)

	return &ZhenAiSpider{
		in:        in,
		out:       out,
		pipe:      pipe,
		pipeChan:  pipeChan,
		workerNum: 10,
	}
}

// Run implements Spider.Run interface
func (zaSpider *ZhenAiSpider) Run() {
	go zaSpider.pipe.Run()
	for i := 0; i < zaSpider.workerNum; i++ {
		go func() {
			for {
				select {
				case req := <-zaSpider.in:
					//zaSpider.pipeChan <- req.Url
					body, err := zaSpider.Download(req.Url)
					//log.Error("Download %s body:%v", req.Url, len(body))
					if err != nil {
						log.Error("Download: error Download url %s %v", req.Url, err)
						continue
					}
					resp, err := req.ParserFunc(body)
					if err != nil {
						log.Error("ParserFunc error:%v", err)
						continue
					}
					zaSpider.out <- &resp
					for item := range resp.Items {
						//log.Error("Item:%v", item)
						zaSpider.pipeChan <- item
					}
				default:
					continue
				}
			}
		}()
	}
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
			Url:        string(c[1]),
			ParserFunc: zaSpider.cityParser,
		})
		if i == 10 {
			break
		}
		i++
	}

	return result, nil
}

func (zaSpider *ZhenAiSpider) cityParser(contents []byte) (common.ParseResult, error) {
	re := regexp.MustCompile(CityRe)
	all := re.FindAllSubmatch(contents, -1)
	//all := re.FindAll(contents, -1)
	//log.Printf("len:%v, all:%v\n", len(contents), len(all))
	result := common.ParseResult{}
	for _, c := range all {
		//result.Items = append(result.Items, string(c[2])) //username
		result.Requests = append(result.Requests, &common.Request{
			Url:        string(c[1]),
			ParserFunc: zaSpider.parseProfile,
		})
	}

	return result, nil
}

func (zaSpider *ZhenAiSpider) parseProfile(contents []byte) (common.ParseResult, error) {
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
