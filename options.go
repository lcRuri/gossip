package gossip

type Options struct {
	Addr string

	Protocol string

	ListenAddr string

	Amount int

	Cycle int64

	Buffer int

	Size int
}

var DefaultOptions = Options{
	Addr:       "0.0.0.0",
	Protocol:   "UDP",
	ListenAddr: "0.0.0.0",
	Amount:     3,
	Cycle:      6,
	Buffer:     9,
	Size:       16384,
}
