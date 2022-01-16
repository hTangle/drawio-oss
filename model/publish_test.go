package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPublisher(t *testing.T) {
	blogs := &BlogList{
		Blogs:   nil,
		BlogMap: nil,
		Path:    "./tem.json",
	}
	blogs.InitBlog()
	arr := map[string]string{
		"111": "abc111de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"114": "abc114de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"115": "abc115de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"113": "abc113de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"112": "abc112de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"116": "abc116de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"121": "abc116de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"117": "abc116de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"118": "abc116de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"119": "abc116de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
		"120": "abc116de![3358356a-63f3-4462-84ba-2ab33c0eab0f.png](https://image.ahsup.top/image/7402bdf88d504f40bb81f9f7097f06be.png)",
	}
	for key, value := range arr {
		blogs.AddAPublisher("", key, value)
	}
	data, err := json.Marshal(blogs.GetBlogsPure(0, 4))
	if err == nil {
		fmt.Printf("%s\n", string(data))
	}
	data, err = json.Marshal(blogs.GetBlogsPure(0, 6))
	if err == nil {
		fmt.Printf("%s\n", string(data))
	}
	data, err = json.Marshal(blogs.GetBlogsPure(3, 6))
	if err == nil {
		fmt.Printf("%s\n", string(data))
	}
	data, err = json.Marshal(blogs.GetBlogsPure(0, 11))
	if err == nil {
		fmt.Printf("%s\n", string(data))
	}
	data, err = json.Marshal(blogs.GetBlogsPure(7, 11))
	if err == nil {
		fmt.Printf("%s\n", string(data))
	}
}
