package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"sync"

	"github.com/redis/go-redis/v9"
)

type API struct {
	Cache              *redis.Client
	AvailableReactions []string
	once               sync.Once
	router             *http.ServeMux
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.once.Do(a.defineRoutes)
	a.router.ServeHTTP(w, r)
}

func (a *API) defineRoutes() {
	a.router = http.NewServeMux()

	a.router.HandleFunc("GET /livestreams/{livestream}/reactions", a.listReactions)
	a.router.HandleFunc("POST /livestreams/{livestream}/reactions/{reaction}", a.sendReaction)
}

func (a *API) listReactions(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(a.AvailableReactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *API) sendReaction(w http.ResponseWriter, r *http.Request) {
	var (
		liveStream = r.PathValue("livestream")
		reaction   = r.PathValue("reaction")
	)

	if !slices.Contains(a.AvailableReactions, reaction) {
		http.Error(w, "reaction not found", http.StatusNotFound)
		return
	}

	if err := a.Cache.Incr(r.Context(), "livestreams:"+liveStream+":reactions:"+reaction).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
