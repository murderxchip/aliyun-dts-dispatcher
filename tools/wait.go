package tools

import (
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"time"
)

func WaitTimeout(waiter string, exit chan interface{}, timeout time.Duration) {
	select {
	case <-time.After(timeout):
		log.Info(define.ServerTag, waiter, " exit timeout, do not wait anymore.")
	case <-exit:
		log.Info(define.ServerTag, waiter, " exit normally, done.")
	}
}
