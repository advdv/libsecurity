package twitter

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/hashicorp/errwrap"
)

type Stream struct {
	api     *anaconda.TwitterApi
	fstream anaconda.Stream
	user    anaconda.User
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
	}, nil
}

func (s *Stream) Events() chan interface{} {
	return s.fstream.C
}

func (s *Stream) Stop() {
	s.fstream.Interrupt()
}

func (s *Stream) Quit() chan bool {
	return s.fstream.Quit
}
