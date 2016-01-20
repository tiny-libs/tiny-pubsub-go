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

	sub1 := pubsub.on("hello")
	pubsub.publish("hello", "world")

	ch.Check(<-sub1, check.Equals, "world")
	pubsub.publish("hello", "globe")

	ch.Check(<-sub1, check.Equals, "globe")

	pubsub.publish("hello", "one")
	pubsub.publish("hello", "two")

	ch.Check(<-sub1, check.Equals, "one")
	ch.Check(<-sub1, check.Equals, "two")
}