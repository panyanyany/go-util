package rod_util

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/go-rod/rod"
)

type WaitForAction2RequestHandler struct {
	UrlPattern    string
	Hook          func(hijack *rod.Hijack) error
	HookStoppable func(hijack *rod.Hijack) (bool, error)
}

func (r WaitForAction2RequestHandler) AddRoute(router *rod.HijackRouter, routing *RoutingStatus, routerResult *MultiRouterResult)  {
	router.MustAdd(r.UrlPattern, func(ctx *rod.Hijack) {
		seelog.Debugf("进入 hijack(%v)(%v), routing=%v", r.UrlPattern, ctx.Request.URL().String(), routing.IsRouting(r.UrlPattern))
		if !routing.IsRouting(r.UrlPattern) {
			return
		}
		defer func() {
			seelog.Debugf("出hijack")
		}()
		var stop bool
		var err error
		if r.Hook != nil {
			err = r.Hook(ctx)
			stop = true
		} else {
			stop, err = r.HookStoppable(ctx)
		}
		if err != nil {
			//seelog.Errorf("is canceled?: %v, routing = %v", errors.Is(err, context.Canceled), routing)
			// 有些请求先进来，但是后面才处理到，此时我们想要的请求已经拿到了，路由结束了，这个请求就会超时
			if routing.IsRouting(r.UrlPattern) == false && errors.Is(err, context.Canceled) {
				err = nil
				return
			}
			//seelog.Errorf("hijack: %v", err)
			err = fmt.Errorf("41 r.Hook: %w", err)
			//chRouter <- err
			routerResult.AddResult(r.UrlPattern, err)
			return
		}
		if stop {
			routing.SetRouting(r.UrlPattern, false)
			routerResult.AddResult(r.UrlPattern, nil)
		}
	})
}

type WaitForAction2Input struct {
	Timeout         time.Duration
	Browser         *rod.Browser
	Page            *rod.Page
	Action          func() error
	StillWait       func() bool
	RequestHandlers []WaitForAction2RequestHandler
}

func WaitForAction2(r WaitForAction2Input) (err error) {
	seelog.Debugf("13 hijack")
	routing := NewRoutingStatus()
	var router *rod.HijackRouter
	if r.Page != nil {
		router = r.Page.HijackRequests()
	} else {
		router = r.Browser.HijackRequests()
	}
	//chBs := make(chan []byte)
	//chRouter := make(chan error)
	routerResult := NewMultiRouterResult(len(r.RequestHandlers))

	ctxBs, cancel := context.WithTimeout(context.Background(), r.Timeout*2)
	defer cancel()
	go func() {
		if r.StillWait != nil {
			for r.StillWait() {
				time.Sleep(time.Second)
			}
			cancel()
		}
	}()

	seelog.Debugf("32 must add")
	for _, rr := range r.RequestHandlers {
		routing.SetRouting(rr.UrlPattern, true)
		seelog.Debugf("adding route-1: %#v", rr.UrlPattern)
		rr.AddRoute(router, routing, routerResult)
	}
	defer func() {
		routing.SetAllRouting(false)
	}()
	go router.Run()
	defer router.Stop()

	chActionErr := make(chan error)
	go func() {
		err = r.Action()
		if err != nil {
			err = fmt.Errorf("r.Action: %w", err)
			chActionErr <- err
			return
		}
	}()

	seelog.Debugf("56 chBs")
	//var bs []byte
	select {
	case err = <-chActionErr:
		break
	case err = <-routerResult.Wait():
		break
	case <-ctxBs.Done():
		err = ctxBs.Err()
		if err != nil {
			err = fmt.Errorf("等待请求超时: %w", err)
			return
		}

		return
	}
	seelog.Debugf("67 chBs done")
	return
}

type MultiRouterResult struct {
	Results           map[string]error
	ExpectResultCount int
	Lock              sync.Mutex
}

func NewMultiRouterResult(threadCnt int) (r *MultiRouterResult) {
	r = new(MultiRouterResult)
	r.Results = make(map[string]error)
	r.ExpectResultCount = threadCnt
	return
}

func (r *MultiRouterResult) Wait() chan error {
	errCh := make(chan error)
	go func() {
		tick := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-tick.C:
				r.Lock.Lock()
				if len(r.Results) < r.ExpectResultCount {
					r.Lock.Unlock()
					continue
				}
				var finalErr error
				for ptn, err := range r.Results {
					if err != nil {
						err = fmt.Errorf("request handler of (%v): %w", ptn, err)
						finalErr = err
						break
					}
				}
				r.Lock.Unlock()
				errCh <- finalErr
				return
			}
		}
	}()

	return errCh
}

func (r *MultiRouterResult) AddResult(ptn string, err error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	r.Results[ptn] = err
	return
}

type RoutingStatus struct {
	Lock    sync.Mutex
	Routing map[string]bool
}

func NewRoutingStatus() (r *RoutingStatus) {
	r = new(RoutingStatus)
	r.Routing = make(map[string]bool)
	return
}
func (r *RoutingStatus) SetAllRouting(routing bool) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	for ptn := range r.Routing {
		r.Routing[ptn] = routing
	}

	return
}

func (r *RoutingStatus) SetRouting(ptn string, routing bool) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	r.Routing[ptn] = routing

	return
}

func (r *RoutingStatus) IsRouting(ptn string) bool {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	return r.Routing[ptn]
}
