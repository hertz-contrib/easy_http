package main

import (
	"fmt"
	"github.com/hertz-contrib/easy_http"
)

func main() {
	c := easy_http.MustNew()

	res, err := c.SetHeader("test", "test").SetQueryParam("test1", "test1").R().Get("http://www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
