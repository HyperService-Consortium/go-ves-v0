package session

import (
	"math/rand"
	"testing"
	"time"
)

func TestActive(t *testing.T) {
	ses := NewMultiThreadSerialSessionBase()
	var isc = []byte("123456789123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ses.ActivateSession(isc)
	time.Sleep(time.Duration(r.Intn(10)) * time.Millisecond)
	ses.InactivateSession(isc)
}

func BenchmarkSessionSet(b *testing.B) {
	b.StopTimer()
	ses := NewMultiThreadSerialSessionBase()
	var isc = []byte("123456789123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := make(chan bool, 30000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			ses.ActivateSession(isc)
			time.Sleep(time.Duration(r.Intn(10)) * time.Millisecond)
			ses.InactivateSession(isc)
			x <- true
		}()
	}
	for i := 0; i < b.N; i++ {
		<-x
	}
}

func BenchmarkPureSessionSet(b *testing.B) {
	b.StopTimer()
	ses := NewMultiThreadSerialSessionBase()
	var isc = []byte("123456789123456789")
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := make(chan bool, 30000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			ses.ActivateSession(isc)
			// time.Sleep(time.Duration(r.Intn(10)) * time.Millisecond)
			ses.InactivateSession(isc)
			x <- true
		}()
	}
	for i := 0; i < b.N; i++ {
		<-x
	}
}
