package common

type ParseResult struct {
	Requests []*Request    `json:"requests"`
	Items    []interface{} `json:"items"`
}

type Request struct {
	Url string `json:"url"`
	//ParserFunc func([]byte) (ParseResult, error)   `json:"parserfunc"` //函数转不了json格式...
	Flag int `json:"flag"` // flag->ParserFunc
}
