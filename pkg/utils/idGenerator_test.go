package utils

import "testing"

func TestIdGenerator_GetId(t *testing.T) {
	generator := NewIdGenerator()

	res := make(chan int64)
	for i := 0; i < 1000000; i++ {
		go func() {
			id := generator.GetId()
			res <- id
		}()
	}

	mp := make(map[int64]int)
	for i := 0; i < 1000000; i++ {
		id := <-res
		if _, ok := mp[id]; ok {
			t.Errorf("find repeat id: %d\n", id)
		}
		mp[id] = i
	}

	defer close(res)
}
