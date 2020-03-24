package streaming

import (
	"sort"

	"github.com/ahab94/streaming-service/models"
)

type Publisher struct {
	userAdd        chan *User
	userRemove     chan *User
	usersByChannel map[string]map[string]*User
	data           chan map[string][]models.Packet
	done           chan struct{}
}

// NewPublisher instantiate publisher
func NewPublisher() *Publisher {
	return &Publisher{
		usersByChannel: make(map[string]map[string]*User),
		userAdd:        make(chan *User),
		userRemove:     make(chan *User),
		data:           make(chan map[string][]models.Packet, 5),
		done:           make(chan struct{}),
	}
}

// Start - publisher goes in listening state
func (p *Publisher) Start() {
	for {
		select {
		case d := <-p.data:
			p.publish(d)
		case u := <-p.userAdd:
			logger.Infof("adding user %s for channel %s", u.ID, u.channelID)
			p.add(u)
		case u := <-p.userRemove:
			logger.Infof("removing user %s for channel %s session %s", u.ID, u.channelID, u.SessionID)
			p.prune(u)
		case <-p.done:
			logger.Info("stopping publisher...")
			return
		}
	}
}

// Stop halts the progression of publisher and close all users' connections
func (p *Publisher) Stop() {
	for _, users := range p.usersByChannel {
		for _, user := range users {
			user.Stop()
		}
	}
	close(p.done)
}

// DataInput returns receiver channel that can be used to push data
func (p *Publisher) DataInput() chan<- map[string][]models.Packet {
	return p.data
}

// AddUserInput returns receiver channel that can be used to add user in publisher
func (p *Publisher) AddUserInput() chan<- *User {
	return p.userAdd
}

// RemoveUserInput returns receiver channel that can be used to remove user from publisher
func (p *Publisher) RemoveUserInput() chan<- *User {
	return p.userRemove
}

// publish pushes packets towards appropriate users
func (p *Publisher) publish(data map[string][]models.Packet) {
	defer recoverPanic()
	for channel, packets := range data {
		sort.Slice(packets, func(i, j int) bool {
			return packets[i].SequenceNum < packets[j].SequenceNum
		})

		for _, user := range p.usersByChannel[channel] {
			user.input <- packets
		}
	}
}

// add adds a user into publisher's list
func (p *Publisher) add(user *User) {
	defer recoverPanic()
	if _, ok := p.usersByChannel[user.channelID]; !ok {
		p.usersByChannel[user.channelID] = make(map[string]*User)
	}
	p.usersByChannel[user.channelID][user.SessionID] = user
}

// prune removes the user from the list
func (p *Publisher) prune(user *User) {
	defer recoverPanic()
	delete(p.usersByChannel[user.channelID], user.SessionID)
}

func recoverPanic() {
	if r := recover(); r != nil {
		logger.Warnf("recovered from panic %+v", r)
	}
}
