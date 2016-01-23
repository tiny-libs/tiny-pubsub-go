package pubsub

import (
)

type channel struct {
	namespace string
	callbacks []func(data interface{})
	channels map[string]*channel
}

type pubsub struct {
	channel
}

func newPubsub() *pubsub {
	ps := new(pubsub)

	ps.channels = make(map[string]*channel)
	return ps
}

func (ps *pubsub) on(namespace string, callback func(data interface{})) *pubsub {
	chann, ok := ps.channels[namespace]
	if(!ok) {
		chann = &channel{ namespace : namespace}
		ps.channels[namespace] = chann
	}

	chann.callbacks = append(chann.callbacks, callback)
	return ps
}

func (ps *pubsub) publish(namespace string, args interface{}) *pubsub {
	chann := ps.channels[namespace]

	for _, callback := range chann.callbacks {
		callback(args)
	}
	return ps
}