package resource_pool

import (
	"sync"
	"sync/atomic"
	"time"
)

type ResourceHandle interface {
	CreateResouce() (interface{}, error)
	CloseResource(*Resouce)
}
type Pool struct {
	rh             ResourceHandle
	resources      *Resouce
	created        int
	lock           sync.Mutex
	min            int
	max            int
	expire_in      time.Duration //expire in second
	available_chan chan *Resouce
	resourcemap    map[*Resouce]bool
}

func NewPool(rh ResourceHandle, min int, max int) *Pool {
	if rh == nil {
		panic("nil create resource function")
	}
	p := &Pool{
		rh:             rh,
		min:            min,
		max:            max,
		available_chan: make(chan *Resouce, max),
		expire_in:      60,
		resources:      &Resouce{},
		created:        0,
	}
	go p.maintain()

	return p
}

func (this *Pool) createResouce() error {
	p, err := this.rh.CreateResouce()
	if err != nil {
		return err
	}
	r := &Resouce{
		R:     p,
		Dirty: nil,
	}
	this.lock.Lock()

	this.resources_len++
	this.lock.Unlock()
	return nil
}

func (this *Pool) maintain() {
	//	index := this.min + 2*(this.max-this.min)/3
	f := func() {

	}
	for {
		time.Sleep(time.Second)
		f()
	}
}

func (this *Pool) Get() *Resouce {
	if this.created < this.max { //策略是创建更多的资源
		go this.createResouce()
	}
	var r *Resouce
	for r = range this.available_chan {
		r.touch()
	}
	return r
}
func (this *Pool) Free(r *Resouce) {
	atomic.StoreUint32(&r.Use, 0)
	this.available_chan <- r
}
func (this *Pool) release() {
	f := func(r *Resouce) {

	}

	p = this.resources
	for p != nil {

	}
}
