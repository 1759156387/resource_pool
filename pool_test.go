package resource_pool

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	data       = 500
	goroutines = 10
	wg         sync.WaitGroup
)

type PoolTest struct {
}

func (this *PoolTest) CreateResouce() (interface{}, error) {
	return sql.Open("mysql", "root:123@tcp(localhost:3306)/test?charset=utf8")
}
func (this *PoolTest) CloseResource(r *Resouce) {
	d := r.R.(*sql.DB)
	d.Close()
}

func Test_pool(t *testing.T) {
	start_time := time.Now()
	defer func() {
		wg.Wait()
		fmt.Println("used time:", time.Now().Sub(start_time))
	}()
	pt := new(PoolTest)
	p := NewPool(pt, 5, 100)

	f := func(numofgoroutinues int, n int) {
		wg.Add(1)
		defer wg.Done()
		for i := 0; i < n; i++ {
			r := p.Get()
			dbconn := r.R.(*sql.DB)
			_, r.R_Dirty = dbconn.Exec("insert into user values(?,?)", fmt.Sprintf("user_%d_%d", numofgoroutinues, n), n%100)
			p.Put(r)
		}
	}
	for i := 0; i < goroutines; i++ {
		if i == goroutines-1 {
			f(i, data/goroutines)
		} else {
			go f(i, data/goroutines)
		}

	}

}

func Test_getconn(t *testing.T) {
	panic(123)
	start_time := time.Now()
	defer func() {
		wg.Wait()
		fmt.Println("used time:", time.Now().Sub(start_time))
	}()

	f := func(numofgoroutinues int, n int) {
		wg.Add(1)
		defer wg.Done()
		for i := 0; i < n; i++ {
			dbconn, _ := sql.Open("mysql", "root:123@tcp(localhost:3306)/test?charset=utf8")
			dbconn.Exec("insert into user values(?,?)", fmt.Sprintf("user_%d_%d", numofgoroutinues, n), n%100)
		}
	}
	for i := 0; i < goroutines; i++ {
		if i == goroutines-1 {
			f(i, data/goroutines)
		} else {
			go f(i, data/goroutines)
		}

	}

}
