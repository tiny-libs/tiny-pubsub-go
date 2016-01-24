package pubsub

import (
)

type channel struct {
	namespace string
	callbacks map[int]func(data interface{})
	lastIndex int
	channels map[string]*channel
}

type subscription struct {
	namespace string
	callback func(data interface{})
	chn *channel
	index int
	ps *Pubsub
}

type Pubsub struct {
	channel
}

func NewPubsub() *Pubsub {
	ps := new(Pubsub)

	ps.channels = make(map[string]*channel)
	return ps
}

func (sub *subscription) off() *Pubsub {
	callbacks := sub.chn.callbacks

	delete(callbacks, sub.index)

	return sub.ps
}

func (ps *Pubsub) on(namespace string, callback func(data interface{})) *subscription {
	chann, ok := ps.channels[namespace]
	if(!ok) {
		chann = &channel{
			namespace : namespace,
			callbacks : make(map[int]func(data interface{})),
		}
		ps.channels[namespace] = chann
	}

	chann.lastIndex++
	chann.callbacks[chann.lastIndex] = callback

	return &subscription{
		namespace,
		callback,
		chann,
		chann.lastIndex,
		ps,
	}
}

func (ps *Pubsub) publish(namespace string, args interface{}) *Pubsub {
	chann := ps.channels[namespace]

	for _, callback := range chann.callbacks {
		callback(args)
	}
	return ps
}