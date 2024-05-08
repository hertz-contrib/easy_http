package eazy_http

import (
	"context"
	"net/http"
	"net/url"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type Request struct {
	client         *Client
	Url            string
	Method         string
	QueryParam     url.Values
	Header         http.Header
	Cookies        []*http.Cookie
	PathParams     map[string]string
	FormParams     map[string]string
	FileParams     map[string]string
	BodyParams     interface{}
	RawRequest     *protocol.Request
	Ctx            context.Context
	RequestOptions []config.RequestOption
	Result         interface{}
	Error          interface{}
}
