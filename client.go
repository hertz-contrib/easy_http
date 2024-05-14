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
	"github.com/cloudwego/hertz/pkg/protocol"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app/client"
)

type Client struct {
	QueryParam url.Values
	PathParams map[string]string
	Header     http.Header
	Cookies    []*http.Cookie

	beforeRequest       []RequestMiddleware
	udBeforeRequest     []RequestMiddleware
	afterResponse       []ResponseMiddleware
	afterResponseLock   *sync.RWMutex
	udBeforeRequestLock *sync.RWMutex

	client *client.Client
}

type (
	RequestMiddleware  func(*Client, *Request) error
	ResponseMiddleware func(*Client, *Response) error
)

func createClient(c *client.Client) *Client {
	return &Client{client: c}
}

func (c *Client) SetQueryParam(param, value string) *Client {
	c.QueryParam.Set(param, value)
	return c
}

func (c *Client) SetQueryParams(params map[string]string) *Client {
	for k, v := range params {
		c.QueryParam.Set(k, v)
	}
	return c
}

func (c *Client) SetQueryParamsFromValues(params url.Values) *Client {
	for k, v := range params {
		for _, v1 := range v {
			c.QueryParam.Add(k, v1)
		}
	}
	return c
}

func (c *Client) SetQueryString(query string) *Client {
	// todo: parse query string
	return c
}

func (c *Client) AddQueryParam(param, value string) *Client {
	c.QueryParam.Add(param, value)
	return c
}

func (c *Client) AddQueryParams(params map[string]string) *Client {
	for k, v := range params {
		c.QueryParam.Add(k, v)
	}
	return c
}

func (c *Client) SetPathParam(param, value string) *Client {
	c.PathParams[param] = value
	return c
}

func (c *Client) SetPathParams(params map[string]string) *Client {
	for k, v := range params {
		c.PathParams[k] = v
	}
	return c
}

func (c *Client) SetHeader(header, value string) *Client {
	c.Header.Set(header, value)
	return c
}

func (c *Client) SetHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.Header.Set(k, v)
	}
	return c
}

func (c *Client) SetHeaderMultiValues(headers map[string][]string) *Client {
	for k, header := range headers {
		for _, v := range header {
			c.Header.Set(k, v)
		}
	}
	return c
}

func (c *Client) AddHeader(header, value string) *Client {
	c.Header.Add(header, value)
	return c
}

func (c *Client) AddHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.Header.Add(k, v)
	}
	return c
}

func (c *Client) AddHeaderMultiValues(headers map[string][]string) *Client {
	for k, header := range headers {
		for _, v := range header {
			c.Header.Add(k, v)
		}
	}
	return c
}

func (c *Client) SetContentType(contentType string) *Client {
	c.Header.Set("Content-Type", contentType)
	return c
}

func (c *Client) SetJSONContentType() *Client {
	c.Header.Set("Content-Type", "application/json")
	return c
}

func (c *Client) SetXMLContentType() *Client {
	c.Header.Set("Content-Type", "application/xml")
	return c
}

func (c *Client) SetHTMLContentType() *Client {
	c.Header.Set("Content-Type", "text/html")
	return c
}

func (c *Client) SetFormContentType() *Client {
	c.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c

}

func (c *Client) SetXFormData() *Client {
	c.Header.Set("Content-Type", "multipart/form-data")
	return c
}

func (c *Client) SetCookie(hc *http.Cookie) *Client {
	c.Cookies = append(c.Cookies, hc)
	return c
}

func (c *Client) SetCookies(hcs []*http.Cookie) *Client {
	c.Cookies = append(c.Cookies, hcs...)
	return c
}

func (c *Client) R() *Request {
	r := &Request{
		QueryParam: url.Values{},
		Header:     http.Header{},
		Cookies:    make([]*http.Cookie, 0),
		PathParams: map[string]string{},

		client: c,
	}
	return r
}

func (c *Client) NewRequest() *Request {
	return c.R()
}

func (c *Client) execute(req *Request) (*Response, error) {
	// Lock the user-defined pre-request hooks.
	c.udBeforeRequestLock.RLock()
	defer c.udBeforeRequestLock.RUnlock()

	// Lock the post-request hooks.
	c.afterResponseLock.RLock()
	defer c.afterResponseLock.RUnlock()

	// Apply Request middleware
	var err error

	// user defined on before request methods
	// to modify the *resty.Request object
	for _, f := range c.udBeforeRequest {
		if err = f(c, req); err != nil {
			return nil, err
		}
	}

	for _, f := range c.beforeRequest {
		if err = f(c, req); err != nil {
			return nil, err
		}
	}

	if hostHeader := req.Header.Get("Host"); hostHeader != "" {
		req.RawRequest.SetHost(hostHeader)
	}

	req.Time = time.Now()

	resp := &protocol.Response{}
	err = c.client.Do(context.Background(), req.RawRequest, resp)

	response := &Response{
		Request:     req,
		RawResponse: resp,
	}

	if err != nil {
		response.receiveAt = time.Now()
		return response, err
	}

	response.receiveAt = time.Now()

	// Apply Response middleware
	for _, f := range c.afterResponse {
		if err = f(c, response); err != nil {
			break
		}
	}

	return response, err
}
