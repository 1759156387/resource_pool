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
	data       = 10000
	goroutines = 10
	wg         sync.WaitGroup
)

func Test_pool(t *testing.T) {
	start_time := time.Now()
	defer func() {
		wg.Wait()
		fmt.Println("used time:", time.Now().Sub(start_time))
	}()
	pt := new(PoolTest)
	p := NewPool(pt, 5, 10)

	f := func(numofgoroutinues int, n int) {
		wg.Add(1)
		defer wg.Done()
		//		gets := 0
		//		getd := 0
		for i := 0; i < n; i++ {
			//			gets++
			//			fmt.Println("gets:", gets)
			r := p.Get()
			//			getd++
			//			fmt.Println("geted", getd)
			dbconn := r.R.(*sql.DB)
			dbconn.Exec("insert into user values(?,?)", fmt.Sprintf("user_%d_%d", numofgoroutinues, n), n%100)
			p.Put(r)
			//			fmt.Println("puted")
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

type PoolTest struct {
}

func (this *PoolTest) CreateResouce() (interface{}, error) {
	return sql.Open("mysql", "root:123@tcp(localhost:3306)/test?charset=utf8")
}
func (this *PoolTest) CloseResource(r *Resouce) {
	d := r.R.(*sql.DB)
	d.Close()
}
