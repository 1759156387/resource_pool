package resource_pool

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

var current_time int64

func init() {
	go func() {
		for {
			current_time = time.Now().Unix()
			time.Sleep(time.Second)
		}
	}()
	go si()
}
func si() {
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGINT) //监听指定信号

	s := <-c //阻塞直至有信号传入
	panic(s)
}
