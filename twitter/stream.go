package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	// "github.com/hashicorp/errwrap"
)

type Stream struct {
	api     *anaconda.TwitterApi
	fstream anaconda.Stream
}

func NewStream() (*Stream, error) {
	api := anaconda.NewTwitterApi("access_token", "access_token_secret")
	api.SetLogger(anaconda.BasicLogger)

	fstream := api.PublicStreamFilter(url.Values{})

	return &Stream{
		api:     api,
		fstream: fstream,
	}, nil
}

func (s *Stream) Start() chan interface{} {
	s.fstream.Start("urlStr", url.Values{}, 0)
	return s.fstream.C
}

func (s *Stream) Quit() chan bool {
	return s.fstream.Quit
}
