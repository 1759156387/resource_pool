package resource_pool

import (
	"time"
)

type ResourceHandle interface {
	CreateResouce() (interface{}, error)
	CloseResource(*Resouce)
}
type Pool struct {
	rh             ResourceHandle
	resources      map[*Resouce]bool
	created        int
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
		resources:      make(map[*Resouce]bool),
		created:        0,
	}
	go p.maintain()

	return p
}

func (this *Pool) createResouce() {
	err_times := 0
	f := func() error {
		p, err := this.rh.CreateResouce()
		if err != nil {
			return err
		}
		r := &Resouce{
			R:       p,
			R_Dirty: nil,
		}
		r.touch()
		this.available_chan <- r
		this.resources[r] = true
		this.created++
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
	f := func() { //close "dirty&&expired" resouce
		for k, _ := range this.resources {
			if k.canFree() {
				this.closeResouce(k)
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

}
