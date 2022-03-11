package model

import (
	"strings"
	"sync"
)

const (
	DefaultImageData = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHkAAAA9CAYAAACJM8YzAAAAAXNSR0IArs4c6QAAA1d0RVh0bXhmaWxlACUzQ214ZmlsZSUyMGhvc3QlM0QlMjJlbWJlZC5kaWFncmFtcy5uZXQlMjIlMjBtb2RpZmllZCUzRCUyMjIwMjItMDMtMTFUMTUlM0EwNyUzQTU4LjExOVolMjIlMjBhZ2VudCUzRCUyMjUuMCUyMChXaW5kb3dzJTIwTlQlMjAxMC4wJTNCJTIwV2luNjQlM0IlMjB4NjQpJTIwQXBwbGVXZWJLaXQlMkY1MzcuMzYlMjAoS0hUTUwlMkMlMjBsaWtlJTIwR2Vja28pJTIwQ2hyb21lJTJGOTkuMC40ODQ0LjUxJTIwU2FmYXJpJTJGNTM3LjM2JTIwRWRnJTJGOTkuMC4xMTUwLjM2JTIyJTIwZXRhZyUzRCUyMlpYTDNxTnFfRmk2bko1dWxPTnp4JTIyJTIwdmVyc2lvbiUzRCUyMjE3LjEuMiUyMiUyMHR5cGUlM0QlMjJlbWJlZCUyMiUzRSUzQ2RpYWdyYW0lMjBpZCUzRCUyMmVtWVJTRERSa0RQMllkM0VaeTNyJTIyJTIwbmFtZSUzRCUyMlBhZ2UtMSUyMiUzRWpaTEJjb1FnRElhZmhydksxTzFldTkzdVhucnkwRE9WVkpnaWNSQlg3ZE5YUzZneU81M3BpZkFsSkg4U0dEJTJCMTA4V0pUcjJpQk1PS1RFNk1QN09peU11c1hJNlZ6SUU4SGc0Qk5FN0xnTElOVlBvTDZHV2tnNWJRRXd2SUl4cXZ1eFRXYUMzVVBtSENPUnpUc0E4MGFkVk9OSkJFcktDcWhibW5iMXA2UlYwVTVjYXZvQnNWSyUyQmZsTVhqZVJmM1pPQndzMWJOb0lYaGFFZE9RaGw0SmllTU84VFBqSjRmb2c5Vk9KekRyV05PSnZmemhKY205bjJNVFVhc0Q2JTJGJTJCVklUWnhFMlpJY3V5U2prcDdxRHBSciUyRmR4MlQ3alQ4cTNacm5saTNsZmszVGR3SG1ZZG9nMFhBQmI4RzVlUXNoN2ZBZ3Y2T3ZrQmFrYXQwWGtjWVpxdDRTU21LRGRONyUyQlp0NTRYZzlxTzEyM2VQNzdkZiUyQmJuYnclM0QlM0QlM0MlMkZkaWFncmFtJTNFJTNDJTJGbXhmaWxlJTNFngzGMAAAAPVJREFUeF7t04EJwEAMw8Bk/6G/FEqHeJ03kIR3Zs7Y1Qb2jXyOzrdW3t0R+da6H5fIlwd+8UQWOWAggOjJIgcMBBA9WeSAgQCiJ4scMBBA9GSRAwYCiJ4scsBAANGTRQ4YCCB6ssgBAwFETxY5YCCA6MkiBwwEED1Z5ICBAKInixwwEED0ZJEDBgKInixywEAA0ZNFDhgIIHqyyAEDAURPFjlgIIDoySIHDAQQPVnkgIEAoieLHDAQQPRkkQMGAoieLHLAQADRk0UOGAggerLIAQMBRE8WOWAggOjJIgcMBBA9WeSAgQCiJ4scMBBA/J8cYE0jPh5C7gE8XAF0AAAAAElFTkSuQmCC"
)

type Image struct {
	Alt  string `json:"alt"`
	Data string `json:"data"`
	Xml  string `json:"xml"`
}

func NewDefaultImage(key string) Image {
	return Image{
		Alt:  key,
		Data: DefaultImageData,
	}
}

type OssListCache struct {
	list map[string]string
	sync.RWMutex
}

func (o *OssListCache) InitCache(list map[string]string) {
	o.Lock()
	defer o.Unlock()
	o.list = list
	if o.list == nil {
		o.list = map[string]string{}
	}
}

func (o *OssListCache) HasObject(key string) bool {
	o.RLock()
	defer o.RUnlock()
	_, ok := o.list[key]
	return ok
}

func (o *OssListCache) GetOssObjects(key string) (results map[string]string) {
	o.RLock()
	defer o.RUnlock()
	results = map[string]string{}
	for k, v := range o.list {
		if strings.HasPrefix(k, key) {
			results[k] = v
		}
	}
	return results
}
func (o *OssListCache) PutOssObject(key, value string) {
	o.Lock()
	defer o.Unlock()
	o.list[key] = value
}
