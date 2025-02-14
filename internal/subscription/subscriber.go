package subscription

import (
	"sync"

	"github.com/erknas/forum/graph/model"
)

type Subscriber interface {
	Subscribe(string) chan *model.Comment
	Unsubscribe(string, chan *model.Comment)
	Publish(string, *model.Comment)
}

type Subscription struct {
	mu    sync.Mutex
	chans map[string]map[chan *model.Comment]struct{}
}

func New() *Subscription {
	return &Subscription{
		chans: make(map[string]map[chan *model.Comment]struct{}),
	}
}

func (s *Subscription) Subscribe(postID string) chan *model.Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan *model.Comment)

	if _, ok := s.chans[postID]; !ok {
		s.chans[postID] = make(map[chan *model.Comment]struct{})
	}

	s.chans[postID][ch] = struct{}{}
	return ch
}

func (s *Subscription) Unsubscribe(postID string, ch chan *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channels, ok := s.chans[postID]; ok {
		if _, ok := channels[ch]; ok {
			close(ch)
			delete(channels, ch)
		}
		if len(channels) == 0 {
			delete(s.chans, postID)
		}
	}
}

func (s *Subscription) Publish(postID string, comment *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channels, ok := s.chans[postID]; ok {
		for ch := range channels {
			select {
			case ch <- comment:
			default:
				close(ch)
				delete(channels, ch)
			}
		}
	}
}
