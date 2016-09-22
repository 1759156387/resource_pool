package resource_pool

import (
	//	"fmt"
	"sync"
	"time"
)

type ResourceHandle interface {
	CreateResouce() (interface{}, error)
	CloseResource(*Resouce)
}
type Pool struct {
	rh              ResourceHandle
	resources       map[*Resouce]bool
	min             int
	max             int
	available_chan  chan *Resouce
	not_dirty       int
	stop            bool
	lock            sync.Mutex
	univerisal_time int64
}

func NewPool(rh ResourceHandle, min int, max int) *Pool {
	if rh == nil {
		panic("nil create resource function")
	}
	p := &Pool{
		rh:              rh,
		min:             min,
		max:             max,
		available_chan:  make(chan *Resouce, max*2),
		resources:       make(map[*Resouce]bool),
		univerisal_time: 0,
	}
	for i := 0; i < min; i++ {
		p.createResouce()
	}
	go p.tick()
	return p
}

func (this *Pool) createResouce() error {
	//fmt.Println("create resouce")
	p, err := this.rh.CreateResouce()
	if err != nil {
		return err
	}
	r := &Resouce{
		R:               p,
		R_Dirty:         nil,
		univerisal_time: &this.univerisal_time,
	}
	r.touch()
	if len(this.available_chan) < this.max {
		this.available_chan <- r
	} else {
		go func() {
			this.available_chan <- r
		}()
	}
	this.lock.Lock()
	this.resources[r] = true
	this.lock.Unlock()
	return nil

}

func (this *Pool) tick() {
	f := func() { //close "dirty&&expired" Resouce
		alive_resouce := 0
		for k, _ := range this.resources {
			if k.canFree() { //canfree need to lock resource,so need to unlock resource
				this.closeResouce(k)
				k.unlock()
				this.lock.Lock()
				delete(this.resources, k)
				this.lock.Unlock()
			}
			if !k.dirty {
				alive_resouce++
			}
		}
		this.not_dirty = alive_resouce
		if needs := this.min - alive_resouce; needs > 0 {
			for i := 0; i < needs; i++ {
				this.createResouce()
			}
		}
	}
	for {
		f()
		this.univerisal_time++
		time.Sleep(time.Second)
	}
}

func (this *Pool) Get() *Resouce {
	if this.not_dirty < this.min {
		go this.createResouce()
	}
	var r *Resouce
	for r = range this.available_chan {
		if r != nil && r.valid() {
			r.touch()
			break
		}
	}
	return r
}
func (this *Pool) Put(r *Resouce) {
	r.unlock()
	if r.R_Dirty != nil {
		r.dirty = true
		return
	}
	if len(this.available_chan) < this.max {
		this.available_chan <- r
		return
	}
	go func() {
		this.available_chan <- r
	}()
}
func (this *Pool) closeResouce(r *Resouce) {
	//fmt.Println("close resouce")
	this.rh.CloseResource(r)
}
