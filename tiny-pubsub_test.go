package pubsub

import (
	check "gopkg.in/check.v1"
	"testing"
)

var _ = check.Suite(new(Suite))

func Test(t *testing.T) {
	check.TestingT(t)
}

type Suite struct{}

func (s *Suite) TestSub(ch *check.C) {
	pubsub := NewPubsub()
	var val string

	pubsub.On("hello", func(data interface{}) {
		val = data.(string)
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

	pubsub.On("hello", func(data interface{}) {
		val, ok := data.(int)
		if(ok) {
			val1 = val
		}
	})
	pubsub.On("hello", func(data interface{}) {
		val, ok := data.(string)
		if(ok) {
			val2 = val
		}
	})
	pubsub.On("hello", func(data interface{}) {
		val, ok := data.(int)
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

	pubsub.On("hello", func(data interface{}) {
		val, ok := data.(string)
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

	pubsub.On("hello", func(data interface{}) {
		go func() {
			counter += data.(int)
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

func (s *Suite) TestUnsubscribe(ch *check.C) {
	pubsub := NewPubsub()
	var counter int

	sub := pubsub.On("hello", func(data interface{}) {
		counter += data.(int)
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