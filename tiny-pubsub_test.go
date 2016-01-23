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
	pubsub := newPubsub()
	var val string

	pubsub.on("hello", func(data interface{}) {
		val = data.(string)
	})
	pubsub.publish("hello", "world")

	ch.Check(val, check.Equals, "world")
	pubsub.publish("hello", "globe")

	ch.Check(val, check.Equals, "globe")

	pubsub.publish("hello", "one")
	pubsub.publish("hello", "two")

	ch.Check(val, check.Equals, "two")
}

func (s *Suite) TestMoreSubs(ch *check.C) {
	pubsub := newPubsub()
	val1 := 0
	val2 := "test"
	val3 := []int{}

	pubsub.on("hello", func(data interface{}) {
		val, ok := data.(int)
		if(ok) {
			val1 = val
		}
	})
	pubsub.on("hello", func(data interface{}) {
		val, ok := data.(string)
		if(ok) {
			val2 = val
		}
	})
	pubsub.on("hello", func(data interface{}) {
		val, ok := data.(int)
		if(ok) {
			val3 = append(val3, val)
		}
	})
	ch.Check(val1, check.Equals, 0)
	ch.Check(val2, check.Equals, "test")
	ch.Check(len(val3), check.Equals, 0)

	pubsub.publish("hello", "world")

	ch.Check(val1, check.Equals, 0)
	ch.Check(val2, check.Equals, "world")
	ch.Check(len(val3), check.Equals, 0)

	pubsub.publish("hello", 5)
	ch.Check(val1, check.Equals, 5)
	ch.Check(val2, check.Equals, "world")
	ch.Check(len(val3), check.Equals, 1)
	ch.Check(val3[0], check.Equals, 5)
}

func (s *Suite) TestChaining(ch *check.C) {
	pubsub := newPubsub()
	val1 := "test"

	pubsub.on("hello", func(data interface{}) {
		val, ok := data.(string)
		if(ok) {
			val1 = val
		}
	})

	pubsub.publish("hello", "world").publish("hello", "globe")

	ch.Check(val1, check.Equals, "globe")
}

func (s *Suite) TestConcurrentCallback(ch *check.C) {
	pubsub := newPubsub()
	var counter int
	ch1 := make(chan int)

	pubsub.on("hello", func(data interface{}) {
		go func() {
			counter += data.(int)
			ch1 <- counter
		}()
	})

	pubsub.publish("hello", 1)

	ch.Check(<-ch1, check.Equals, 1)
	ch.Check(counter, check.Equals, 1)

	pubsub.publish("hello", 2)
	ch.Check(<-ch1, check.Equals, 3)
	ch.Check(counter, check.Equals, 3)
}

func (s *Suite) TestUnsubscribe(ch *check.C) {
	pubsub := newPubsub()
	var counter int

	sub := pubsub.on("hello", func(data interface{}) {
		counter += data.(int)
	})

	pubsub.publish("hello", 3)

	ch.Check(counter, check.Equals, 3)

	pubsub.publish("hello", 5)
	ch.Check(counter, check.Equals, 8)

	sub.off()
	pubsub.publish("hello", 5)
	ch.Check(counter, check.Equals, 8)
	pubsub.publish("hello", 5)
	ch.Check(counter, check.Equals, 8)
}