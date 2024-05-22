package easy_http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetQueryParam(t *testing.T) {
	c := MustNew()
	c.SetQueryParam("test1", "test1")
	c.SetQueryParams(map[string]string{"test2": "test2", "test3": "test3"})
	c.SetQueryParamsFromValues(map[string][]string{"test4": {"test41", "test42"}})
	c.SetQueryString("test5=test5")

	assert.Equal(t, "test1", c.QueryParam.Get("test1"))
	assert.Equal(t, "test2", c.QueryParam.Get("test2"))
	assert.Equal(t, "test3", c.QueryParam.Get("test3"))
	assert.Equal(t, []string{"test41", "test42"}, c.QueryParam["test4"])
	assert.Equal(t, "test5", c.QueryParam.Get("test5"))
}
