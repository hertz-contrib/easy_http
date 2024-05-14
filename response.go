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
	"github.com/cloudwego/hertz/pkg/protocol"
	"strings"
	"time"
)

type Response struct {
	Request     *Request // 上面的 Request 结构体
	RawResponse *protocol.Response

	receiveAt time.Time

	bodyByte []byte
	size     int64
}

func (r *Response) Body() []byte {
	if r.RawResponse == nil {
		return []byte{}
	}
	return r.bodyByte
}

func (r *Response) BodyString() string {
	if r.RawResponse == nil {
		return ""
	}
	return strings.TrimSpace(string(r.bodyByte))
}

func (r *Response) StatusCode() int {
	if r.RawResponse == nil {
		return 0
	}
	return r.RawResponse.StatusCode()
}

func (r *Response) Result() interface{} {
	return r.Request.Result
}

func (r *Response) Error() interface{} {
	return r.Request.Error
}

// todo
//func (r *Response) Header() http.Header {
//	if r.RawResponse == nil {
//		return http.Header{}
//	}
//	return r.RawResponse.Header.GetHeaders()
//}
//
//func (r *Response) Cookies() []*http.Cookie {
//	if r.RawResponse == nil {
//		return make([]*http.Cookie, 0)
//	}
//	return r.RawResponse.Header.GetCookies()
//}
//func (r *Response) ToRawHTTPResponse() string {
//	return r.RawResponse.String()
//}

func (r *Response) IsSuccess() bool {
	return r.StatusCode() > 199 && r.StatusCode() < 300
}

func (r *Response) IsError() bool {
	return r.StatusCode() > 399
}
