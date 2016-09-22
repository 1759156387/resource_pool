package resource_pool

import (
	"os"
	"os/signal"
	"syscall"
)

func init() {

}
func signalHandle() {
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGINT) //监听指定信号

	s := <-c //阻塞直至有信号传入
	panic(s)
}
