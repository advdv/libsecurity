package twitter

type Poller struct{}

func NewPoller() (*Poller, error) {
	return &Poller{}, nil
}
