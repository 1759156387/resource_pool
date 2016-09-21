package resource_pool

import (
	"sync"
	"sync/atomic"
	"time"
)

type Resouce struct {
	R      interface{}
	Dirty  error
	next   *Resouce
	prev   *Resouce
	expire int64
	use    uint32
	lock   sync.Mutex
}

func (this *Resouce) init() {
	this.prev = this
	this.next = this
}

func (this *Resouce) push(r *Resouce) {
	ori := this.prev
	ori.next = r
	r.prev = ori
	r.next = this
	this.prev = r
}

func (this *Resouce) touch() error {
	this.last_access = time.Now().Unix() + 60
}

func (this *Resouce) tryLock() bool {
	this.lock.Lock()
	if atomic.LoadUint32(&this.use) != 0 { //lock failed
		return false
	}
	atomic.StoreUint32(&this.use, 1)
	this.lock.Unlock()
}
func (this *Resouce) unlock() {
	atomic.StoreUint32(&this.use, 0)
}

func (this *Resouce) del(node *Resouce) {
	p := this
	for {
		if p != node && p.next == node {
			break
		}
		p = p.next
		if p == nil {
			break
		}
	}
	if p == nil { // not found,so nothing
		return
	}
	p.next = node.next
}

/*
func (this *Resouce) length() int {
	i := 0
	var p *Resouce
	p = this
	for {
		p = p.next
		if p == nil {
			break
		}
		i++
	}
	return i
}
func (this *Resouce) index(index int) *Resouce {
	if index < 0 {
		return nil
	}
	p := this
	for i := 0; i < index+1; i++ {
		p = p.next
		if p == nil {
			break
		}
	}
	return p
}

func (this *Resouce) add(lis *Resouce) {

}
func (this *Resouce) push() {

}

func (this *Resouce) del(node *Resouce) { // you can`t delete head node
	p := this
	for {
		if p != node && p.next == node {
			break
		}
		p = p.next
		if p == nil {
			break
		}
	}
	if p == nil { // not found,so nothing
		return
	}
	p.next = node.next
}

func (this *Resouce) travel() {
	i := 0
	head := this
	p := this.next
	for p != nil {
		fmt.Printf("i:%d payload:%v\n ", i, p.payload)
		p = p.next
		i++
	}
}
*/
