package thread_util

import (
    "github.com/cihub/seelog"
    "time"
)

type Handler func()

type FuncTicker struct {
    Interval         time.Duration
    Running          bool
    Handler          Handler
    stopped          chan bool
    StartImmediately bool
}

func NewFuncTicker(interval time.Duration) (r *FuncTicker) {
    r = new(FuncTicker)
    r.Interval = interval
    r.stopped = make(chan bool)
    return
}

func (r *FuncTicker) Start() {
    r.Running = true
    go r.loop()
}

func (r *FuncTicker) Stop() {
    seelog.Infof("func ticker stopping")
    r.Running = false
    <-r.stopped
    seelog.Infof("func ticker stopped")
}

func (r *FuncTicker) loop() {
    tick := time.NewTicker(r.Interval)
    if r.StartImmediately {
        r.Handler()
    }
    for {
        if !r.Running {
            r.stopped <- true
            break
        }
        select {
        case <-tick.C:
            r.Handler()
        }
    }
}

func SetInterval(h Handler, interval time.Duration, startImmediately bool) *FuncTicker {
    t := NewFuncTicker(interval)
    t.Handler = h
    t.StartImmediately = startImmediately
    t.Start()
    return t
}
