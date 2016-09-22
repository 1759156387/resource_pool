package resource_pool

import (
	"fmt"
	"time"
)

type ResourceHandle interface {
	CreateResouce() (interface{}, error)
	CloseResource(*Resouce)
}
type Pool struct {
	rh             ResourceHandle
	resources      map[*Resouce]bool
	min            int
	max            int
	available_chan chan *Resouce
}

func NewPool(rh ResourceHandle, min int, max int) *Pool {
	if rh == nil {
		panic("nil create resource function")
	}
	p := &Pool{
		rh:             rh,
		min:            min,
		max:            max,
		available_chan: make(chan *Resouce, max*2),
		resources:      make(map[*Resouce]bool),
	}
	for i := 0; i < min; i++ {
		p.createResouce()
	}
	go p.maintain()
	return p
}

func (this *Pool) createResouce() {
	fmt.Println("create resouce")
	err_times := 0
	f := func() error {
		p, err := this.rh.CreateResouce()
		if err != nil {
			return err
		}
		r := &Resouce{
			R:         p,
			R_Dirty:   nil,
			expire_in: time.Now().Unix() + 60,
		}
		r.touch()
		this.available_chan <- r
		this.resources[r] = true
		return nil
	}
	for {
		e := f()
		if e != nil {
			err_times++
		} else {
			return
		}
		if err_times == 5 {
			panic("create resource returned 5 times error")
		}
		time.Sleep(time.Second)
	}

}

func (this *Pool) maintain() {
	alive_resouce := 0
	f := func() { //close "dirty&&expired" resouce
		for k, _ := range this.resources {
			if k.canFree() { //canfree need to lock resource,so need to unlock resource
				this.closeResouce(k)
				k.unlock()
				delete(this.resources, k)
			}
			if !k.dirty {
				alive_resouce++
			}

		}
		if needs := this.min - alive_resouce; needs > 0 {
			for i := 0; i < needs; i++ {
				this.createResouce()
			}
		}

	}
	for {
		time.Sleep(time.Second)
		f()
	}
}

func (this *Pool) Get() *Resouce {
	var r *Resouce
	for r = range this.available_chan {
		if r != nil && r.valid() {
			break
		}
	}
	return r
}
func (this *Pool) Put(r *Resouce) {
	r.unlock()
	if r.R_Dirty != nil {
		r.dirty = true
	} else {
		this.available_chan <- r
	}
}
func (this *Pool) closeResouce(r *Resouce) {
	fmt.Println("close resouce")
	this.rh.CloseResource(r)
}
