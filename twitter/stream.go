package twitter

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/hashicorp/errwrap"
)

type EventType int

// matches "CVE-2013-11111 in fed34532hash432"
var EventNewVulnerabilityExp = regexp.MustCompile(`(CVE-\d+-\d+)\s+in\s+([^\s-]+)`)

// matches "fix $(selector) with feda3fe222"
var EventFixVulnerabilityExp = regexp.MustCompile(`fix\s+([^-\s]*)\s*with\s+([^-\s]+)`)

var (
	EventNewVulnerability = EventType(1)
	EventFixVulnerability = EventType(2)
)

type Event struct {
	Selector string
	CVE      string
	Image    string
	Type     EventType
	Tweet    anaconda.Tweet
}

type Stream struct {
	api     *anaconda.TwitterApi
	fstream anaconda.Stream
	user    anaconda.User
	events  chan Event

	replies []anaconda.Tweet
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
		replies: []anaconda.Tweet{},
	}, nil
}

//@todo do something intelligent with tweet
func (s *Stream) handleTweet(t anaconda.Tweet) error {
	ev := Event{
		Tweet: t,
	}

	if EventNewVulnerabilityExp.MatchString(t.Text) {
		m := EventNewVulnerabilityExp.FindStringSubmatch(t.Text)
		if len(m) < 3 {
			s.api.Log.Errorf("Only %d regexp matches in text '%s'", len(m), t.Text)
			return nil
		}

		ev.CVE = m[1]
		ev.Image = m[2]
		ev.Type = EventNewVulnerability
	} else if EventFixVulnerabilityExp.MatchString(t.Text) {
		m := EventFixVulnerabilityExp.FindStringSubmatch(t.Text)
		if len(m) < 3 {
			s.api.Log.Errorf("Only %d regexp matches in text '%s'", len(m), t.Text)
			return nil
		}

		ev.Selector = m[1]
		ev.Image = m[2]
		ev.Type = EventFixVulnerability
	} else {
		s.api.Log.Debugf("Tweet %d didn't match any special action", t.Id)
		return nil
	}

	s.events <- ev
	return nil
}

func (s *Stream) ReplyVulnerable(t anaconda.Tweet, hostname, imgid, cid string) error {
	vals := url.Values{}
	vals.Set("in_reply_to_status_id", strconv.FormatInt(t.Id, 10))

	_, err := s.api.PostTweet(fmt.Sprintf("@%s host %s is vulnerable, container: %s", t.User.ScreenName, hostname, cid), vals)
	if err != nil {
		return err
	}

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
