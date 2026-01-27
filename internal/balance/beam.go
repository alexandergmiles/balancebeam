package balance

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Beam is our load balancer. For the time being, it will only accept HTTP traffic
type Beam struct {
	Endpoints []string
	Router    *chi.Router
}

func NewBeam() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	return router
}

func (b *Beam) RegisterRoutes(routes map[string](func(w http.ResponseWriter, r *http.Request) error)) error {
	for route, handler := range routes {
		fmt.Printf("%s %v\n", route, handler)
	}

	return nil
}
