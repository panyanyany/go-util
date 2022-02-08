package cclient
// concurrency client

import (
	"context"
	"math"
	"sync"
	"time"

	"tugou-hunter/repo/config"

	"github.com/cihub/seelog"
	"github.com/go-redis/redis/v8"
	"github.com/panyanyany/go-web3/jsonrpc"
)

type CClient struct {
	Rdb          *redis.Client
	RedisTimeout time.Duration
	EnableCache  bool
	FindLock     sync.Mutex
	LastRequest  time.Time

	QueryApiList   []string
	PostApiList    []string
	HistoryApiList []string
	UrlLocks       map[string]*sync.Mutex

	ApiLastRequest map[string]time.Time

	Clients   map[string]jsonrpc.IClient
	Endpoints *jsonrpc.Endpoints

	ReqInterval int64
}

func (r *CClient) Close() error {
	// 不用管，Eth 不会调用它
	return nil
}

func New(rdb *redis.Client, redisTimeout time.Duration) (r *CClient) {
	r = new(CClient)
	r.RedisTimeout = redisTimeout
	r.Rdb = rdb
	r.FindLock = sync.Mutex{}
	r.EnableCache = rdb != nil && redisTimeout.Seconds() != 0
	r.ReqInterval = 100
	r.QueryApiList = config.Inst.Endpoints
	r.PostApiList = r.QueryApiList
	r.HistoryApiList = r.QueryApiList

	r.ApiLastRequest = map[string]time.Time{}
	r.Clients = map[string]jsonrpc.IClient{}
	r.UrlLocks = map[string]*sync.Mutex{}

	allUrlList := []string{}
	allUrlListExists := make(map[string]bool)
	for _, url := range r.QueryApiList {
		_, found := allUrlListExists[url]
		if found {
			continue
		}
		allUrlListExists[url] = true
		allUrlList = append(allUrlList, url)
	}
	for _, url := range r.PostApiList {
		_, found := allUrlListExists[url]
		if found {
			continue
		}
		allUrlListExists[url] = true
		allUrlList = append(allUrlList, url)
	}
	for _, url := range r.HistoryApiList {
		_, found := allUrlListExists[url]
		if found {
			continue
		}
		allUrlListExists[url] = true
		allUrlList = append(allUrlList, url)
	}

	for _, url := range allUrlList {
		r.ApiLastRequest[url] = time.Time{}
		api, err := jsonrpc.NewClient(url)
		if err != nil {
			panic(err)
		}
		//api := jsonrpc.Eth{Provider: r}
		if r.Endpoints == nil {
			r.Endpoints = api.Endpoints
		}
		r.Clients[url] = api
		r.UrlLocks[url] = &sync.Mutex{}
	}
	return
}

var ctx = context.Background()

func (r *CClient) GetLock(apiList []string) string {
	r.FindLock.Lock()
	var foundUrl string
	interval := int64(math.MaxInt64)

	for _, url := range apiList {
		if r.ApiLastRequest[url].IsZero() {
			foundUrl = url
			interval = 0
			break
		}
		now := time.Now()
		diff := r.ReqInterval - int64(now.Sub(r.ApiLastRequest[url]).Milliseconds())
		if diff < interval {
			interval = diff
			foundUrl = url
		}
		//seelog.Infof("diff=%v, interval=%v, foundUrl=%v, url=%v, last=%v, now=%v",
		//	diff, interval, foundUrl,
		//	url, r.ApiLastRequest[url], now,
		//)
	}
	r.ApiLastRequest[foundUrl] = time.Now()
	r.FindLock.Unlock()

	r.UrlLocks[foundUrl].Lock()
	if interval > 0 {
		restTime := time.Duration(interval) * time.Millisecond
		seelog.Debugf("sleep: %v", restTime)
		time.Sleep(restTime)
	}
	return foundUrl
}
func (r *CClient) ReleaseLock(url string) {
	//r.LastRequest = time.Now()
	//r.ApiLastRequest[url] = time.Now()
	//r.FindLock.Unlock()
	r.UrlLocks[url].Unlock()
}

func (r *CClient) Call(method string, out interface{}, params ...interface{}) (err error) {
	baseUrl := r.GetLock(r.PostApiList)
	defer r.ReleaseLock(baseUrl)

	api := r.Clients[baseUrl]
	seelog.Debugf("using url: %v, call=%v", baseUrl, method)

	err = api.Call(method, out, params...)
	return
}
