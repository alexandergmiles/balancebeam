package balance

type Forwarder interface {
	Register(target string, destination string) error
	GetBalancedAddress(from string) (string, error)
}
