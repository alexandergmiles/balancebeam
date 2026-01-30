package balance

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

// Beam is our load balancer. For the time being, it will only accept HTTP traffic
type Beam struct {
	Balancer  Forwarder
	Endpoints []string
	Router    *chi.Mux
}

func NewBeam(logger slog.Logger) *Beam {
	forwarder := NewRobinCoordinator()
	beam := Beam{
		Balancer: forwarder,
		Router:   chi.NewRouter(),
	}

	beam.Router.Use(slogchi.New(&logger))

	return &beam
}

func (b *Beam) Register(route string, destination string) error {
	b.Balancer.Register(route, destination)
	b.Router.Get(route, b.Balance)
	return nil
}

func (b *Beam) MapRoute(to string, backend string) {
}

func (b *Beam) Balance(w http.ResponseWriter, r *http.Request) {
	balancedTarget, err := b.Balancer.GetBalancedAddress(r.URL.Path)
	if err != nil {
		w.Write([]byte("unable to get balanced address"))
		return
	}

	res, err := b.Get(balancedTarget, r)
	fmt.Println(res)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to complete request: %s\n", err)))
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.Write([]byte("unable to decode body"))
		return
	}
	w.Write(body)
}

func (b *Beam) Get(endpoint string, from *http.Request) (*http.Response, error) {
	log.Printf("Getting: %s, %s, %s\n", endpoint, from.URL.Host, from.URL.Path)
	newRequest := http.Request{}
	newRequest.URL = &url.URL{
		Host:   endpoint,
		Scheme: from.URL.Scheme,
		Path:   from.URL.Path,
		User:   from.URL.User,
	}

	newRequest.URL.Scheme = from.URL.Scheme
	if from.URL.Scheme == "" {
		newRequest.URL.Scheme = "http"
	}

	newRequest.Header = from.Header

	client := http.Client{}

	res, err := client.Do(&newRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %s", err)
	}

	return res, nil
}
