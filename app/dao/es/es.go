package es

import (
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/tal-tech/xtools/confutil"
	"net/http"
	"sync"
	"time"
)

const (
	defaultHTTPTimeout = 10 * time.Second
)

var (
	client        *elastic.Client
	once          sync.Once
)

func EsClient() *elastic.Client {
	if client == nil {
		if err := InitEngine(); err != nil {
		}
	}
	return client
}

//InitEngine 初始化
func InitEngine() (err error) {
	once.Do(func() {
		//load config
		cfg := confutil.GetConfStringMap("ElasticSearch")

		URLs := confutil.GetConfs("ElasticSearch", "url")
		if len(URLs) == 0 {
			err = errors.New("config ElasticSearch.url does not exists")
			return
		}
		options := []elastic.ClientOptionFunc{
			elastic.SetURL(URLs...),
			elastic.SetSniff(cast.ToBool(cfg["sniff"])),
			elastic.SetHealthcheck(true),
			elastic.SetHealthcheckInterval(10 * time.Second),
		}
		//http超时
		timeout := cast.ToDuration(cfg["timeout"])
		if timeout == 0 {
			timeout = defaultHTTPTimeout
		}
		httpclient := &http.Client{
			Timeout: timeout,
		}
		options = append(options, elastic.SetHttpClient(httpclient))
		//创建
		client, err = elastic.NewClient(options...)
		if err != nil {
			return
		}
	})
	return
}