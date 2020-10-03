package sr

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Service holds the map of minions and provides methods CRUD operations on the map
type Service struct {
	connectionString string
	minions          map[string]Minion
	sync.RWMutex
}

// NewService returns a Service with a connectionString configured and can be a map of minions setup. The minions map can be empty,
// or can contain minions
func NewService(connectionString string, minions map[string]Minion) *Service {
	return &Service{
		connectionString: connectionString,
		minions:          minions,
	}
}

// ListenAndServe registers the routes to the server and starts the server on the host:port configured in Service
func (s *Service) ListenAndServe() error {
	r := mux.NewRouter()

	// Each handler is wrapped in logs() and auth() to log out the method and path and to
	// ensure that a non-empty Authorization header is present
	r.HandleFunc("/minion", logs(auth(s.PostMinion))).Methods("POST")
	r.HandleFunc("/minion", logs(auth(s.GetMinions))).Methods("GET")
	r.HandleFunc("/minion/{name}", logs(auth(s.GetMinion))).Methods("GET")
	r.HandleFunc("/minion/{name}", logs(auth(s.PutMinion))).Methods("PUT")
	r.HandleFunc("/minion/{name}", logs(auth(s.DeleteMinion))).Methods("DELETE")

	log.Printf("Starting server on %s", s.connectionString)
	err := http.ListenAndServe(s.connectionString, r)
	if err != nil {
		return err
	}
	return nil
}

// logs prints the Method and Path to stdout
func logs(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		log.Printf("%s %s", method, path)
		handlerFunc(w, r)
		return
	}
}

// auth checks that a non-empty authorization header has been sent with the request
func auth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Please supply and Authorization token", http.StatusUnauthorized)
			return
		}
		handlerFunc(w, r)
		return
	}
}
