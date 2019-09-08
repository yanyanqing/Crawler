package common

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

type Request struct {
	Url        string
	ParserFunc func([]byte) (ParseResult, error)
}
