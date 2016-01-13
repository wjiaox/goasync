package async

import (
	"runtime"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	async := NewAsync()
	async.Add("1", Try, 100000)
	async.Add("2", Try, 200000)
	async.Add("3", Try, 300000)
	async.Add("4", Try, 400000)
	async.Add("11", Try, 100000)
	async.Add("22", Try, 200000)
	async.Add("33", Try, 300000)
	async.Add("44", Try, 400000)
	err := async.Add("5", Try, 500000)
	if err != nil {
		t.Log("add error:", err.Error())
	}
	t.Log("counter:", async.Count)

	nt := time.Now()
	v, err := async.Go()
	at := time.Now().Sub(nt)

	if err != nil {
		t.Log("error:", err.Error())
	} else {
		for k, val := range v {
			t.Log("name:", k, "result:", val)
		}
	}
	t.Log("use time:", at.Nanoseconds())

	tnt := time.Now()
	Try(100000)
	Try(200000)
	Try(300000)
	Try(400000)
	Try(500000)
	Try(100000)
	Try(200000)
	Try(300000)
	Try(400000)
	tat := time.Now().Sub(tnt)
	t.Log("common,use time:", tat.Nanoseconds())

}

func Try(a int) (int, bool) {
	var re int
	for i := 0; i < a; i++ {
		re += i
	}
	return re, true
}
