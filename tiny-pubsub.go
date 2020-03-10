package pubsub

import (
	"sync"
	"sync/atomic"
	// "log"
)

type channel struct {
	namespace string
	callbacks map[int]func(data []interface{})
	lastIndex uint64
	channels *sync.Map
	mutex *sync.Mutex
}

type Subscription struct {
	namespace string
	callback func(data []interface{})
	chn *channel
	index uint64
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

	ps.channels = &sync.Map{}


	return ps
}

func (sub *Subscription) Off() *Pubsub {
	mutex := sub.chn.mutex

	mutex.Lock()
	callbacks := sub.chn.callbacks

	delete(callbacks, int(sub.index))
	mutex.Unlock()

	return sub.ps
}

func (ps *Pubsub) On(namespace string, callback func(data []interface{})) *Subscription {
	chanLoad, ok := ps.channels.Load(namespace)
	var chann *channel

	if(ok == false) {
		chann = &channel{
			namespace : namespace,
			callbacks : make(map[int]func(data []interface{})),
			mutex : &sync.Mutex{},
		}
		ps.channels.Store(namespace, chann)
	} else {
		chann = chanLoad.(*channel)
	}
	// log.Println("CHANN: ", chann)

	chann.mutex.Lock()

	atomic.AddUint64(&chann.lastIndex, uint64(1))
	chann.callbacks[int(chann.lastIndex)] = callback
	subscription := &Subscription{
		namespace,
		callback,
		chann,
		chann.lastIndex,
		ps,
	}

	chann.mutex.Unlock()

	return subscription
}

func (ps *Pubsub) Publish(namespace string, args ...interface{}) *Pubsub {
	chanLoad, ok := ps.channels.Load(namespace)

	if(!ok) {
		return ps
	}

	chann := chanLoad.(*channel)

	chann.mutex.Lock()
	for _, callback := range chann.callbacks {
		if(ps.concurrent) {
			go callback(args)
		} else {
			callback(args)
		}
	}
	chann.mutex.Unlock()

	return ps
}