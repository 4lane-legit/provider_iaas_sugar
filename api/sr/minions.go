package sr

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// Minion represents a single Minion
type Minion struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// GetMinions returns all of the Minions that exist in the server
func (s *Service) GetMinions(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	s.shuffleTags()
	err := json.NewEncoder(w).Encode(s.minions)
	if err != nil {
		log.Println(err)
	}
}

// PostMinion handles adding a new Minion
func (s *Service) PostMinion(w http.ResponseWriter, r *http.Request) {
	var m Minion
	if r.Body == nil {
		http.Error(w, "invalid request: ", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	space := regexp.MustCompile(`\s+`)
	if space.Match([]byte(m.Name)) {
		http.Error(w, "malformed request with whitespaces", 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.minionExists(m.Name) {
		http.Error(w, fmt.Sprintf("duplicate alert!!! %s ", m.Name), http.StatusBadRequest)
		return
	}

	s.minions[m.Name] = m
	log.Printf("added minion: %s", m.Name)
	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

// PutMinion handles updating an Minion with a specific name
func (s *Service) PutMinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	minionName := vars["name"]
	if minionName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var minion Minion
	if r.Body == nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&minion)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	if !s.minionExists(minionName) {
		log.Printf("minion %s does not exist", minionName)
		http.Error(w, fmt.Sprintf("minion %v does not exist", minionName), http.StatusBadRequest)
		return
	}

	s.minions[minionName] = minion
	log.Printf("updated minion: %s", minion.Name)
	err = json.NewEncoder(w).Encode(minion)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

// DeleteMinion handles removing an Minion with a specific name
func (s *Service) DeleteMinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	minionName := vars["name"]
	if minionName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	s.Lock()
	defer s.Unlock()

	if !s.minionExists(minionName) {
		http.Error(w, fmt.Sprintf("minion %s does not exists", minionName), http.StatusNotFound)
		return
	}

	delete(s.minions, minionName)

	_, err := fmt.Fprintf(w, "Deleted minion with name %s", minionName)
	if err != nil {
		log.Println(err)
	}
}

// GetMinion handles retrieving an Minion with a specific name
func (s *Service) GetMinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	minionName := vars["name"]
	if minionName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	s.RLock()
	defer s.RUnlock()
	s.shuffleTags()
	if !s.minionExists(minionName) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(s.minions[minionName])
	if err != nil {
		log.Println(err)
		return
	}
}

// minionExists checks if an minion exists in or not. Does not lock access to the minionService, expects this to
// be done by the calling method
func (s *Service) minionExists(minionName string) bool {
	if _, ok := s.minions[minionName]; ok {
		return true
	}
	return false
}

// suffleMinionTags shuffles the order of the tags within each minion in the minionService.Does not lock access
// to the minionService, expects this to be done by the calling method
func (s *Service) shuffleTags() {
	for _, minion := range s.minions {
		for i := range minion.Tags {
			j := rand.Intn(i + 1)
			minion.Tags[i], minion.Tags[j] = minion.Tags[j], minion.Tags[i]
		}
	}
}
