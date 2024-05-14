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
	"net/http"
	"net/url"
	"time"

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
	Time           time.Time
	Ctx            context.Context
	RequestOptions []config.RequestOption
	Result         interface{}
	Error          interface{}
}
