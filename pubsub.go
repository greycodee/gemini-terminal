package main

import (
	"sync"
)

type item struct {
	topic string
	data  subItem
}

type subItem struct {
	id   int64
	data string
}

type subscriber struct {
	id   int64
	buff chan subItem
}

type Publisher struct {
	m           sync.RWMutex
	subscribers map[string][]*subscriber
	ch          chan item
}

func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(map[string][]*subscriber),
		ch:          make(chan item),
	}
}

func (p *Publisher) Start() {
	for i := range p.ch {

		p.m.RLock()
		if subs, found := p.subscribers[i.topic]; found {
			for _, sub := range subs {
				sub.buff <- i.data
			}
		}
		p.m.RUnlock()
	}
}

func (p *Publisher) Subscribe(topic string, id int64) *subscriber {
	p.m.Lock()
	defer p.m.Unlock()

	sub := &subscriber{
		id:   id,
		buff: make(chan subItem, 1),
	}

	p.subscribers[topic] = append(p.subscribers[topic], sub)

	return sub
}

func (p *Publisher) Publish(topic string, data subItem) {
	p.ch <- item{topic: topic, data: data}
}

// func main() {
// 	p := NewPublisher()

// 	go func() {
// 		chatId := int64(0)
// 		sub1 := p.Subscribe("topic1", chatId)
// 		for item := range sub1.buff {
// 			if item.id == chatId {
// 				fmt.Printf("from: %d, msg:%s \n", item.id, item.data)
// 			}
// 		}
// 	}()

// 	go func() {
// 		chatId := int64(1)
// 		sub1 := p.Subscribe("topic1", chatId)
// 		for item := range sub1.buff {
// 			if item.id == chatId {
// 				fmt.Printf("from: %d, msg:%s\n", item.id, item.data)
// 			}
// 		}
// 	}()

// 	// sub2 := p.Subscribe("topic1", 2)

// 	go func() {
// 		sum := 1
// 		for {
// 			id := int64(sum % 2)
// 			p.Publish("topic1", subItem{
// 				id:   id,
// 				data: "hello",
// 			})
// 			time.Sleep(1 * time.Second)
// 			sum++
// 		}
// 	}()

// 	go p.Start()

// 	select {}

// }
