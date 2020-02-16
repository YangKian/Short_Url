package utils

import (
	"testing"
	"time"
)

func TestNewIdGenerator_generateNewNode(t *testing.T) {
	_, err := NewIdGenerator(time.Now().UnixNano() / 1e6, 0, 12, 10)
	if err != nil {
		t.Fatal("create new node err: ", err)
	}

	_, err = NewIdGenerator(time.Now().UnixNano() / 1e6, 5000, 12, 10)
	if err == nil {
		t.Fatal("create new node err: ", err)
	}
}

func TestIdGenerator_GetId(t *testing.T) {
	generator, err := NewIdGenerator(time.Now().UnixNano() / 1e6, 100, 12, 10)
	if err != nil {
		t.Fatal(err)
	}

	res := make(chan int64)
	for i := 0; i < 1000000; i++ {
		go func() {
			id, err := generator.GetId()
			if err != nil {
				t.Error(err)
			}
			res <- id
		}()
	}

	mp := make(map[int64]int, 1000000)
	for i := 0; i < 1000000; i++ {
		id := <-res
		if _, ok := mp[id]; ok {
			t.Errorf("find repeat id: %d\n", id)
		}
		mp[id] = i
	}

	defer close(res)
}

func TestMultIdGenerator_GetId(t *testing.T) {
	start := time.Now().UnixNano() / 1e6
	node1, _ := NewIdGenerator(start, 1000, 12, 10)
	node2, _ := NewIdGenerator(start, 30, 12, 10)
	node3, _ := NewIdGenerator(start, 500, 12, 10)
	node4, _ := NewIdGenerator(start, 750, 12, 10)

	res := make(chan int64)
	go func() {
		for i := 0; i < 1000000; i++ {
			id, err := node1.GetId()
			if err != nil {
				t.Error(err)
			}
			res <- id
		}
	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			id, err := node3.GetId()
			if err != nil {
				t.Error(err)
			}
			res <- id
		}
	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			id, err := node2.GetId()
			if err != nil {
				t.Error(err)
			}
			res <- id
		}
	}()

	go func() {
		for i := 0; i < 1000000; i++ {
			id, err := node4.GetId()
			if err != nil {
				t.Error(err)
			}
			res <- id
		}
	}()

	mp := make(map[int64]int, 4000000)
	for i := 0; i < 4000000; i++ {
		id := <-res
		if _, ok := mp[id]; ok {
			t.Errorf("find repeat id: %d\n", id)
		}
		mp[id] = i
	}

	defer close(res)
}

func BenchmarkIdGenerator(b *testing.B) {
	generator, err := NewIdGenerator(time.Now().UnixNano() / 1e6, 100, 12, 10)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = generator.GetId()
	}
}