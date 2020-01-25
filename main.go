package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

func main() {
	var redisURL, httpListenAddr string

	if redisURL = os.Getenv("REDIS_URI"); redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	if httpListenAddr = os.Getenv("LISTEN_ADDR"); httpListenAddr == "" {
		httpListenAddr = ":12345"
	}

	flag.StringVar(&redisURL, "d", redisURL, "redis database URL")
	flag.StringVar(&httpListenAddr, "l", httpListenAddr, "HTTP listen address")
	flag.Parse()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL %s", err)
	}

	client := redis.NewClient(opts)
	if _, err := client.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	h := &handler{&redisStore{redis.NewClient(opts)}}

	log.Printf("Listening on %s...", httpListenAddr)
	log.Fatal(http.ListenAndServe(httpListenAddr, h))
}

type store interface {
	get(string) string
	put(string, string)
}

type redisStore struct {
	r *redis.Client
}

func (s *redisStore) get(k string) string {
	v, err := s.r.Get(k).Result()
	if err != nil && err != redis.Nil {
		log.Fatal(err)
	}
	return v
}

func (s *redisStore) put(k string, v string) {
	if err := s.r.Set(k, v, 0).Err(); err != nil {
		log.Printf("Error setting %q (%s)", k, err)
	}
}

type handler struct {
	store store
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.store.put(r.URL.Path[1:], strings.Split(r.RemoteAddr, ":")[0])
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodGet:
		addr := h.store.get(r.URL.Path[1:])
		if addr == "" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, addr)
	default:
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
	}
}
