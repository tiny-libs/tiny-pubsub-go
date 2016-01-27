package pubsub

import (
)

type channel struct {
	namespace string
	callbacks map[int]func(data interface{})
	lastIndex int
	channels map[string]*channel
}

type Subscription struct {
	namespace string
	callback func(data interface{})
	chn *channel
	index int
	ps *Pubsub
}

type Pubsub struct {
	channel
	concurrent bool
}

func NewPubsub(config ...bool) *Pubsub {
	ps := new(Pubsub)

	if(len(config) == 1) {
		ps.concurrent = config[0]
	}

	ps.channels = make(map[string]*channel)


	return ps
}

func (sub *Subscription) Off() *Pubsub {
	callbacks := sub.chn.callbacks

	delete(callbacks, sub.index)

	return sub.ps
}

func (ps *Pubsub) On(namespace string, callback func(data interface{})) *Subscription {
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

	return &Subscription{
		namespace,
		callback,
		chann,
		chann.lastIndex,
		ps,
	}
}

func (ps *Pubsub) Publish(namespace string, args interface{}) *Pubsub {
	chann, ok := ps.channels[namespace]
	if(!ok) {
		return ps
	}

	for _, callback := range chann.callbacks {
		if(ps.concurrent) {
			go callback(args)
		} else {
			callback(args)
		}
	}

	return ps
}