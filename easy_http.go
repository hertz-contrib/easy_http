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
	"crypto/tls"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/retry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network"
	"time"
)

type Option struct {
	c *client.Client
	F []config.ClientOption
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) WithHertzRawOption(c *client.Client) *Option {
	o.c = c
	return o
}

func (o *Option) WithDialTimeout(dialTimeout time.Duration) *Option {
	o.F = append(o.F, client.WithDialTimeout(dialTimeout))
	return o
}

func (o *Option) WithMaxConnsPerHost(mc int) *Option {
	o.F = append(o.F, client.WithMaxConnsPerHost(mc))
	return o
}

func (o *Option) WithMaxIdleConnDuration(t time.Duration) *Option {
	o.F = append(o.F, client.WithMaxIdleConnDuration(t))
	return o
}

func (o *Option) WithMaxConnDuration(t time.Duration) *Option {
	o.F = append(o.F, client.WithMaxConnDuration(t))
	return o
}

func (o *Option) WithMaxConnWaitTimeout(t time.Duration) *Option {
	o.F = append(o.F, client.WithMaxConnWaitTimeout(t))
	return o
}

func (o *Option) WithKeepAlive(keepAlive bool) *Option {
	o.F = append(o.F, client.WithKeepAlive(keepAlive))
	return o
}

func (o *Option) WithTLSConfig(tlsConfig *tls.Config) *Option {
	o.F = append(o.F, client.WithTLSConfig(tlsConfig))
	return o
}

func (o *Option) WithDialer(dialer network.Dialer) *Option {
	o.F = append(o.F, client.WithDialer(dialer))
	return o
}

func (o *Option) WithResponseBodyStream(flag bool) *Option {
	o.F = append(o.F, client.WithResponseBodyStream(flag))
	return o
}

func (o *Option) WithDisableHeaderNamesNormalizing(flag bool) *Option {
	o.F = append(o.F, client.WithDisableHeaderNamesNormalizing(flag))
	return o
}

func (o *Option) WithName(name string) *Option {
	o.F = append(o.F, client.WithName(name))
	return o
}

func (o *Option) WithNoDefaultUserAgentHeader(flag bool) *Option {
	o.F = append(o.F, client.WithNoDefaultUserAgentHeader(flag))
	return o
}

func (o *Option) WithDisablePathNormalizing(flag bool) *Option {
	o.F = append(o.F, client.WithDisablePathNormalizing(flag))
	return o
}

func (o *Option) WithRetryConfig(retryConfig retry.Option) *Option {
	o.F = append(o.F, client.WithRetryConfig(retryConfig))
	return o
}

func (o *Option) WithWriteTimeout(t time.Duration) *Option {
	o.F = append(o.F, client.WithWriteTimeout(t))
	return o
}

func (o *Option) WithConnStateObserve(hs config.HostClientStateFunc, interval ...time.Duration) *Option {
	o.F = append(o.F, client.WithConnStateObserve(hs, interval...))
	return o
}

func (o *Option) WithDialFunc(f network.DialFunc) *Option {
	o.F = append(o.F, client.WithDialFunc(f))
	return o
}

func (o *Option) WithHostClientConfigHook(h func(hc interface{}) error) *Option {
	o.F = append(o.F, client.WithHostClientConfigHook(h))
	return o
}

func NewClient(opts *Option) (*Client, error) {
	var hertzOptions []config.ClientOption

	if opts.c != nil {
		return createClient(opts.c), nil
	}

	for _, f := range opts.F {
		hertzOptions = append(hertzOptions, f)
	}

	c, err := client.NewClient(hertzOptions...)
	return createClient(c), err
}

func MustNewClient(opts *Option) *Client {
	var hertzOptions []config.ClientOption

	if opts.c != nil {
		return createClient(opts.c)
	}

	for _, f := range opts.F {
		hertzOptions = append(hertzOptions, f)
	}

	c, err := client.NewClient(hertzOptions...)
	if err != nil {
		panic(err)
	}
	return createClient(c)
}
