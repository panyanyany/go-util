package rod_util

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cihub/seelog"
	"github.com/go-rod/rod"
)

type WaitForActionInput struct {
	Timeout       time.Duration
	Browser       *rod.Browser
	Page          *rod.Page
	Action        func() error
	UrlPattern    string
	Hook          func(hijack *rod.Hijack) error
	HookStoppable func(hijack *rod.Hijack) (bool, error)
	StillWait     func() bool
}

func WaitForAction(r WaitForActionInput) (err error) {
	seelog.Debugf("13 hijack")
	var routing bool
	var router *rod.HijackRouter
	if r.Page != nil {
		router = r.Page.HijackRequests()
	} else {
		router = r.Browser.HijackRequests()
	}
	//chBs := make(chan []byte)
	chRouter := make(chan error)

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
	router.MustAdd(r.UrlPattern, func(ctx *rod.Hijack) {
		if !routing {
			return
		}
		seelog.Debugf("进入 hijack, routing=%v", routing)
		defer func() {
			seelog.Debugf("出hijack")
		}()
		var stop bool
		if r.Hook != nil {
			err = r.Hook(ctx)
			stop = true
		} else {
			stop, err = r.HookStoppable(ctx)
		}
		//reqUrl := ctx.Request.URL()
		//seelog.Debugf("hijack url: %v", reqUrl)
		//gTk = reqUrl.Query().Get("g_tk")
		//agencyUid = reqUrl.Query().Get("agency_uid")
		//err = ctx.LoadResponse(http.DefaultClient, true)
		if err != nil {
			//seelog.Errorf("is canceled?: %v, routing = %v", errors.Is(err, context.Canceled), routing)
			// 有些请求先进来，但是后面才处理到，此时我们想要的请求已经拿到了，路由结束了，这个请求就会超时
			if routing == false && errors.Is(err, context.Canceled) {
				err = nil
				return
			}
			//seelog.Errorf("hijack: %v", err)
			err = fmt.Errorf("41 r.Hook: %w", err)
			chRouter <- err
			return
		}
		if stop {
			routing = false
			chRouter <- nil
		}

		//bs := []byte(ctx.Response.Body())
		//ioutil.WriteFile("./storage/get_wechat_list.json", bs, 0644)
		//chBs <- bs
	})
	routing = true
	defer func() {
		routing = false
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
	case err = <-chRouter:
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
