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
	"github.com/cloudwego/hertz/pkg/app/client/discovery"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type Request struct {
	client         *Client
	URL            string
	Method         string
	QueryParam     url.Values
	FormData       url.Values
	Header         http.Header
	PathParams     map[string]string
	RawRequest     *protocol.Request
	Ctx            context.Context
	RequestOptions []config.RequestOption
	Result         interface{}
	Error          interface{}
	isMultiPart    bool
	multipartFiles []*File
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

type File struct {
	Name      string
	ParamName string
	io.Reader
}

func (r *Request) SetQueryParam(param, value string) *Request {
	r.QueryParam.Set(param, value)
	return r
}
func (r *Request) SetQueryParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}
func (r *Request) SetQueryParamsFromValues(params url.Values) *Request {
	for p, v := range params {
		for _, pv := range v {
			r.QueryParam.Add(p, pv)
		}
	}
	return r
}
func (r *Request) SetQueryString(query string) *Request {
	r.RawRequest.SetQueryString(query)
	return r
}
func (r *Request) AddQueryParam(params, value string) *Request {
	r.QueryParam.Add(params, value)
	return r
}
func (r *Request) AddQueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.AddQueryParam(k, v)
	}
	return r
}
func (r *Request) SetPathParam(param, value string) *Request {
	r.PathParams[param] = value
	return r
}
func (r *Request) SetPathParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetPathParam(p, v)
	}
	return r
}

func (r *Request) SetHeader(header, value string) *Request {
	r.Header.Set(header, value)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

func (r *Request) SetHeaderMultiValues(headers map[string][]string) *Request {
	for key, values := range headers {
		r.SetHeader(key, strings.Join(values, ", "))
	}
	return r
}
func (r *Request) AddHeader(header, value string) *Request {
	r.Header.Add(header, value)
	return r
}
func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.AddHeader(k, v)
	}
	return r
}
func (r *Request) AddHeaderMultiValues(headers map[string][]string) *Request {
	for key, value := range headers {
		r.AddHeader(key, strings.Join(value, ", "))
	}
	return r
}

func (r *Request) SetCookie(hc *http.Cookie) *Request {
	r.RawRequest.SetCookie(hc.Name, hc.Value)
	return r
}
func (r *Request) SetCookies(rs []*http.Cookie) *Request {
	for _, c := range rs {
		r.RawRequest.SetCookie(c.Name, c.Value)
	}
	return r
}

func (r *Request) SetBody(body interface{}) *Request {
	t := reflect.Indirect(reflect.ValueOf(body)).Type().Kind()

	switch t {
	case reflect.String:
		r.RawRequest.SetBodyString(body.(string))
	case reflect.TypeOf([]byte{}).Kind():
		r.RawRequest.SetBody(body.([]byte))
	case reflect.TypeOf(io.Reader(nil)).Kind():
		r.RawRequest.SetBodyStream(body.(io.Reader), -1)
	default:
		panic("unsupported body type")
	}

	return r
}
func (r *Request) SetFormData(data map[string]string) *Request {
	for k, v := range data {
		r.FormData.Set(k, v)
	}
	return r
}
func (r *Request) SetFormDataFromValues(data url.Values) *Request {
	for key, value := range data {
		for _, v := range value {
			r.FormData.Add(key, v)
		}
	}
	return r
}
func (r *Request) SetFiles(files map[string]string) *Request {
	r.isMultiPart = true
	for f, fp := range files {
		r.FormData.Set("@"+f, fp)
	}
	return r
}

func (r *Request) SetFileReader(param, fileName string, reader io.Reader) *Request {
	r.isMultiPart = true
	r.multipartFiles = append(r.multipartFiles, &File{
		Name:      fileName,
		ParamName: param,
		Reader:    reader,
	})
	return r
}
func (r *Request) SetResult(res interface{}) *Request {
	if res != nil {
		vv := reflect.ValueOf(res)
		if vv.Kind() == reflect.Ptr {
			r.Result = res
		} else {
			r.Result = reflect.New(vv.Type()).Interface()
		}
	}
	return r
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.Ctx = ctx
	return r
}

const (
	defaultNetwork = "tcp"
)

type customizedResolver struct {
	Address string
}

var _ discovery.Resolver = (*customizedResolver)(nil)

// NewResolver create a service resolver.
func NewResolver(address string) discovery.Resolver {
	return &customizedResolver{
		Address: address,
	}
}

// Target return a description for the given target that is suitable for being a key for cache.
func (c *customizedResolver) Target(_ context.Context, target *discovery.TargetInfo) (description string) {
	return target.Host
}

// Name returns the name of the resolver.
func (c *customizedResolver) Name() string {
	return "easy_http"
}

// Resolve a service info by desc.
func (c *customizedResolver) Resolve(_ context.Context, desc string) (discovery.Result, error) {
	var eps []discovery.Instance

	tags := map[string]string{}
	eps = append(eps, discovery.NewInstance(
		defaultNetwork,
		c.Address,
		1,
		tags,
	))

	return discovery.Result{
		CacheKey:  desc,
		Instances: eps,
	}, nil
}

func (r *Request) WithDC(dc string) *Request {
	resolver := NewResolver(dc)
	r.client.client.Use(sd.Discovery(resolver))
	return r
}
func (r *Request) WithCluster() *Request {
	return r
}
func (r *Request) WithEnv() *Request {
	return r
}
func (r *Request) WIthCallTimeout(t time.Duration) *Request {
	r.RawRequest.SetOptions(config.WithDialTimeout(t))
	return r
}
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
