package resource_pool

import (
	"sync"
)

type Resouce struct {
	R         interface{}
	R_Dirty   error //使用资源出现了错误
	dirty     bool
	expire_in int64
	use       uint32
	lock      sync.Mutex
}

func (this *Resouce) touch() { //unix-style modify expire_in
	this.expire_in = current_time + 3
}

func (this *Resouce) tryLock() bool { // lock resource
	this.lock.Lock()
	if this.use != 0 { //lock failed
		this.lock.Unlock()
		return false
	}
	this.use = 1
	this.lock.Unlock()
	return true
}

func (this *Resouce) unlock() { //unlock resource
	this.lock.Lock()
	this.use = 0
	this.lock.Unlock()
}

func (this *Resouce) valid() bool { //if resource valid
	canlock := this.tryLock()
	if !canlock {
		return false
	}
	valid := (!this.expired()) && !this.dirty
	if !valid {
		this.unlock()
	}
	return valid
}

func (this *Resouce) expired() bool {
	this.lock.Lock()
	b := this.expire_in < current_time //
	if b {
		this.dirty = true
	}
	this.lock.Unlock()
	return b
}
func (this *Resouce) canFree() bool { //test this resourc can be free and lock resource
	if this.dirty || this.expired() {
		return this.tryLock()
	}
	return false
}
