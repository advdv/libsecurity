package twitter

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/hashicorp/errwrap"
)

type Event struct {
	Tweet anaconda.Tweet
}

type Stream struct {
	api     *anaconda.TwitterApi
	fstream anaconda.Stream
	user    anaconda.User
	events  chan Event
}

func NewStream(name string) (*Stream, error) {
	anaconda.SetConsumerKey(consumer_key)
	anaconda.SetConsumerSecret(consumer_secret)

	api := anaconda.NewTwitterApi(access_token, access_token_secret)
	api.SetLogger(anaconda.BasicLogger)

	api.Log.Debugf("Fetching user %s...", name)
	u, err := api.GetUsersShow(name, nil)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Failed to getUserShow() for %s: {{err}}", name), err)
	}
	api.Log.Debugf("Found user with id %d", u.Id)

	vals := url.Values{}
	vals.Set("follow", strconv.FormatInt(u.Id, 10))

	return &Stream{
		api:     api,
		fstream: api.PublicStreamFilter(vals),
		user:    u,
		events:  make(chan Event),
	}, nil
}

//@todo do something intelligent with tweet
func (s *Stream) handleTweet(t anaconda.Tweet) error {
	ev := Event{
		Tweet: t,
	}

	s.events <- ev

	return nil
}

func (s *Stream) Start() chan Event {
	go func() {
		for {
			select {
			case <-s.fstream.Quit:
				close(s.events)
				return
			case msg := <-s.fstream.C:
				switch v := msg.(type) {
				case anaconda.Tweet:
					err := s.handleTweet(v)
					if err != nil {
						s.api.Log.Errorf("error while handling tweet: %s", err)
					}
				}
			}
		}
	}()

	return s.events
}

func (s *Stream) Events() chan interface{} {
	return s.fstream.C
}

func (s *Stream) Stop() {
	s.fstream.Interrupt()
}
