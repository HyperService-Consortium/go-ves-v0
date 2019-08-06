package bitmap

import (
	"fmt"
	"testing"

	redispool "github.com/Myriad-Dreamin/go-ves/lib/database/redis"
)

func TestBitMap(t *testing.T) {
	bm, err := NewBitMap([]byte("orz"), 7, redispool.RedisCacheClient.Pool.Get())
	fmt.Println(err)
	fmt.Println(bm.Get(3))
	fmt.Println(bm.Set(3))
	fmt.Println(bm.Get(3))
	fmt.Println(bm.Count())
	fmt.Println(bm.Reset(3))
	fmt.Println(bm.Count())
	fmt.Println(bm.length)
	fmt.Println(bm.InLength(7))
	fmt.Println(bm.InLength(6))
	fmt.Println(bm.length)
	fmt.Println(bm.Clear())

	fmt.Println(bm.Get(3))
}

func BenchmarkSet(b *testing.B) {
	b.StopTimer()
	bm, _ := NewBitMap([]byte("orz"), 7, redispool.RedisCacheClient.Pool.Get())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bm.Set(3)
	}
}

func BenchmarkGSet(b *testing.B) {
	b.StopTimer()
	bm, _ := NewBitMap([]byte("orz"), 7, redispool.RedisCacheClient.Pool.Get())
	x := make(chan bool, 30000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			bm.Set(3)
			x <- true
		}()
	}
	for i := 0; i < b.N; i++ {
		<-x
	}
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	bm, _ := NewBitMap([]byte("orz"), 7, redispool.RedisCacheClient.Pool.Get())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bm.Get(3)
	}
}
func BenchmarkClear(b *testing.B) {
	b.StopTimer()
	bm, _ := NewBitMap([]byte("orz"), 7, redispool.RedisCacheClient.Pool.Get())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bm.Clear()
	}
}
func BenchmarkInlength(b *testing.B) {
	b.StopTimer()
	bm, _ := NewBitMap([]byte("orz"), 7, redispool.RedisCacheClient.Pool.Get())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bm.InLength(3)
	}
}
