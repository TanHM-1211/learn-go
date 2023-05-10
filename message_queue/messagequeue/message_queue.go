package messagequeue

import "errors"

var (
	errTopicNotFound = errors.New("topic not found")
)

type MessageQueue struct {
	mq map[string]chan *Node
}

func newMessageQueue() *MessageQueue {
	return &MessageQueue{make(map[string]chan *Node)}
}

func (mq *MessageQueue) AddTopic(topic string, capacity int) bool {
	if _, ok := mq.mq[topic]; ok {
		return false
	}
	mq.mq[topic] = make(chan *Node, capacity)
	return true
}

func (mq *MessageQueue) HasTopic(topic string) bool {
	if _, ok := mq.mq[topic]; ok {
		return true
	}
	return false
}

func (mq *MessageQueue) Get(topic string) (*Message, error) {
	if !mq.HasTopic(topic) {
		return nil, errTopicNotFound
	}
	node := <-mq.mq[topic]
	return node.data, nil
}

func (mq *MessageQueue) Put(topic string, message *Message) error {
	if !mq.HasTopic(topic) {
		return errTopicNotFound
	}
	node := &Node{data: message}
	mq.mq[topic] <- node
	return nil
}
