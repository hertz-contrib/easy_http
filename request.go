/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package easy_http

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"net/http"
	"net/url"
)

type Request struct {
	client         *Client
	URL            string
	Method         string
	QueryParam     url.Values
	Header         http.Header
	Cookies        []*http.Cookie
	PathParams     map[string]string
	FormData       map[string]string
	FileData       map[string]string
	BodyParams     interface{}
	RawRequest     *protocol.Request
	Ctx            context.Context
	RequestOptions []config.RequestOption
	Result         interface{}
	Error          interface{}
}

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

func (r *Request) Get(url string) (*Response, error) {
	return r.Execute(MethodGet, url)
}

func (r *Request) Head(url string) (*Response, error) {
	return r.Execute(MethodHead, url)
}

func (r *Request) Post(url string) (*Response, error) {
	return r.Execute(MethodPost, url)
}

func (r *Request) Put(url string) (*Response, error) {
	return r.Execute(MethodPut, url)
}

func (r *Request) Delete(url string) (*Response, error) {
	return r.Execute(MethodDelete, url)
}

func (r *Request) Options(url string) (*Response, error) {
	return r.Execute(MethodOptions, url)
}

func (r *Request) Patch(url string) (*Response, error) {
	return r.Execute(MethodPatch, url)
}

func (r *Request) Send() (*Response, error) {
	return r.Execute(r.Method, r.URL)
}

func (r *Request) Execute(method, url string) (*Response, error) {
	r.Method = method

	r.RawRequest.SetRequestURI(url)
	res := &Response{
		Request:     r,
		RawResponse: &protocol.Response{},
	}

	var err error
	res, err = r.client.execute(r)
	return res, err
}
