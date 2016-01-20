package pubsub

type channel struct {
	namespace string
	subscriptions []chan interface{}
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

func (ps *pubsub) on(namespace string) chan interface{} {
	currSub := make(chan interface{})
	chann := &channel{ namespace : namespace}

	chann.subscriptions = append(chann.subscriptions, currSub)
	ps.channels[namespace] = chann

	return currSub
}

func (ps *pubsub) publish(namespace string, args string) *pubsub {
	chann := ps.channels[namespace]

	for _, sub := range chann.subscriptions {
		go func(args string) {
			sub <- args
		}(args)
	}
	return ps
}