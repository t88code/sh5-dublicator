package saver

import (
	"net/http"
)

type Saver interface {
	OnQuery(req *http.Request, body []byte, reqName string)
	OnReply(res *http.Response, body []byte, reqName string)
}
