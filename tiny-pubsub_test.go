package pubsub

import (
	check "gopkg.in/check.v1"
	"testing"
	"log"
	"sync/atomic"
	"sync"
)

var _ = check.Suite(new(Suite))

func Test(t *testing.T) {
	check.TestingT(t)
}

type Suite struct{}

func (s *Suite) TestSub(ch *check.C) {
	pubsub := NewPubsub()
	var val string

	pubsub.On("hello", func(data []interface{}) {
		val = data[0].(string)
	})
	pubsub.Publish("hello", "world")

	ch.Check(val, check.Equals, "world")
	pubsub.Publish("hello", "globe")

	ch.Check(val, check.Equals, "globe")

	pubsub.Publish("hello", "one")
	pubsub.Publish("hello", "two")

	ch.Check(val, check.Equals, "two")
}

func (s *Suite) TestMoreSubs(ch *check.C) {
	pubsub := NewPubsub()
	val1 := 0
	val2 := "test"
	val3 := []int{}

	pubsub.On("hello", func(data []interface{}) {
		val, ok := data[0].(int)
		if(ok) {
			val1 = val
		}
	})
	pubsub.On("hello", func(data []interface{}) {
		val, ok := data[0].(string)
		if(ok) {
			val2 = val
		}
	})
	pubsub.On("hello", func(data []interface{}) {
		val, ok := data[0].(int)
		if(ok) {
			val3 = append(val3, val)
		}
	})
	ch.Check(val1, check.Equals, 0)
	ch.Check(val2, check.Equals, "test")
	ch.Check(len(val3), check.Equals, 0)

	pubsub.Publish("hello", "world")

	ch.Check(val1, check.Equals, 0)
	ch.Check(val2, check.Equals, "world")
	ch.Check(len(val3), check.Equals, 0)

	pubsub.Publish("hello", 5)
	ch.Check(val1, check.Equals, 5)
	ch.Check(val2, check.Equals, "world")
	ch.Check(len(val3), check.Equals, 1)
	ch.Check(val3[0], check.Equals, 5)
}

func (s *Suite) TestChaining(ch *check.C) {
	pubsub := NewPubsub()
	val1 := "test"

	pubsub.On("hello", func(data []interface{}) {
		val, ok := data[0].(string)
		if(ok) {
			val1 = val
		}
	})

	pubsub.Publish("hello", "world").Publish("hello", "globe")

	ch.Check(val1, check.Equals, "globe")
}

func (s *Suite) TestConcurrentCallback(ch *check.C) {
	pubsub := NewPubsub()
	var counter int
	ch1 := make(chan int)

	pubsub.On("hello", func(data []interface{}) {
		go func() {
			counter += data[0].(int)
			ch1 <- counter
		}()
	})

	pubsub.Publish("hello", 1)

	ch.Check(<-ch1, check.Equals, 1)
	ch.Check(counter, check.Equals, 1)

	pubsub.Publish("hello", 2)
	ch.Check(<-ch1, check.Equals, 3)
	ch.Check(counter, check.Equals, 3)
}

func (s *Suite) TestMultipleConcurrentCallbacks(ch *check.C) {
	pubsub := NewPubsub(true)
	var counter uint64
	var wg sync.WaitGroup

	// ch1 := make(chan uint64)
    for i := 1; i <= 100; i++ {
		pubsub.On("hello", func(data []interface{}) {
			defer wg.Done()

			val, ok := data[0].(int)

			if (ok) {
				atomic.AddUint64(&counter, uint64(val))
			}
		})
    }

    for i := 1; i <= 2; i++ {
	    wg.Add(100)
		go pubsub.Publish("hello", i)
    }

	wg.Wait()

	ch.Check(counter, check.Equals, uint64(100 + 200))

	counter = uint64(0)

    for i := 1; i <= 5; i++ {
	    wg.Add(100)
		go pubsub.Publish("hello", i)
    }

	wg.Wait()

	ch.Check(counter, check.Equals, uint64(100 + 200 + 300 + 400 + 500))
}

func (s *Suite) TestUnsubscribe(ch *check.C) {
	pubsub := NewPubsub()
	var counter int

	sub := pubsub.On("hello", func(data []interface{}) {
		counter += data[0].(int)
	})

	pubsub.Publish("hello", 3)

	ch.Check(counter, check.Equals, 3)

	pubsub.Publish("hello", 5)
	ch.Check(counter, check.Equals, 8)

	sub.Off()
	pubsub.Publish("hello", 5)
	ch.Check(counter, check.Equals, 8)
	pubsub.Publish("hello", 5)
	ch.Check(counter, check.Equals, 8)
}

func (s *Suite) TestPubsubNewGoroutine(ch *check.C) {
	pubsub := NewPubsub(true)
	var counter uint64
	ch1 := make(chan uint64)

	sub := pubsub.On("hello", func(data []interface{}) {
		val, ok := data[0].(int)

		if (ok) {
			atomic.AddUint64(&counter, uint64(val))
		}

		ch1 <- counter
	})

	ch.Check(counter, check.Equals, uint64(0))
	pubsub.Publish("hello", 3)

	<-ch1
	ch.Check(counter, check.Equals, uint64(3))

	ch.Check(counter, check.Equals, uint64(3))
	pubsub.Publish("hello", 5)
	<-ch1
	ch.Check(counter, check.Equals, uint64(8))

	sub.Off()
	pubsub.Publish("hello", 5)
	ch.Check(counter, check.Equals, uint64(8))
	pubsub.Publish("hello", 5)
	ch.Check(counter, check.Equals, uint64(8))
}

func (s *Suite) TestPublishToUnsubscribed(ch *check.C) {
	pubsub := NewPubsub(true)
	counter := uint64(0)
	ch1 := make(chan uint64)

	sub := pubsub.On("hello", func(data []interface{}) {
		val, ok := data[0].(int)

		if (ok) {
			atomic.AddUint64(&counter, uint64(val))
		}

		ch1 <- counter
	})

	ch.Check(counter, check.Equals, uint64(0))
	pubsub.Publish("hello", 3)
	<-ch1
	ch.Check(counter, check.Equals, uint64(3))

	pubsub.Publish("hello", 5)
	<-ch1
	ch.Check(counter, check.Equals, uint64(8))

	sub.Off()

	pubsub.Publish("hello", 5)
	ch.Check(counter, check.Equals, uint64(8))
	pubsub.Publish("world", 5)
	ch.Check(counter, check.Equals, uint64(8))
}

func (s *Suite) TestPublishMultipleArgs(ch *check.C) {
	pubsub := NewPubsub()
	counter := 0

	pubsub.On("hello", func(data []interface{}) {
		log.Println(data)

		for _, num := range data {
			log.Println(num)
			counter += num.(int)
		}
	})

	pubsub.Publish("hello", 3)
	ch.Check(counter, check.Equals, 3)

	pubsub.Publish("hello", 1, 2)
	ch.Check(counter, check.Equals, 6)
}

// Race detection test case
func (s *Suite) TestConcurrentSubscribePublish(ch *check.C) {
	pubsub := NewPubsub(true)
	var counter uint64

    for i := 1; i <= 10; i++ {
		pubsub.On("hello", func(data []interface{}) {
			val, ok := data[0].(int)

			if (ok) {
				atomic.AddUint64(&counter, uint64(val))
			}
			go pubsub.On("hello" + string(i), func(data []interface{}) {
			})
			go pubsub.Publish("hello" + string(i), i)
		})
    }

    for i := 1; i <= 5; i++ {
		go pubsub.Publish("hello", i)
    }
}