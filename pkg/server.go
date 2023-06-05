package pkg

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// App represents the server's internal state.
// It holds configuration about providers and content
type App struct {
	ContentClients map[Provider]Client
	Config         ContentMix
}

func (a App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.String())
	v := req.URL.Query()
	if len(v["count"]) == 0 || len(v["offset"]) == 0 || len(v["count"][0]) == 0 || len(v["offset"][0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	count, err := strconv.Atoi(v["count"][0])
	if err != nil {
		log.Fatal(err)
	}
	offset, err := strconv.Atoi(v["offset"][0])
	if err != nil {
		log.Fatal(err)
	}
	var result []*ContentItem
	for i := 0; len(result) < count; i++ {
		r, err := a.ContentClients[a.Config[i%len(a.Config)].Type].GetContent(req.URL.Hostname(), 1+offset)
		if err != nil {
			if a.Config[i%len(a.Config)].Fallback != nil {
				r, err := a.ContentClients[*a.Config[i%len(a.Config)].Fallback].GetContent(req.URL.Hostname(), 1+offset)
				if err != nil {
					result = append(result, r[len(r)-1])
					if err := json.NewEncoder(w).Encode(result); err != nil {
						log.Print(err)
						w.WriteHeader(http.StatusInternalServerError)
					}
				} else {
					result = append(result, r[len(r)-1])
				}
			} else {
				if err := json.NewEncoder(w).Encode(result); err != nil {
					log.Print(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		} else {
			result = append(result, r[len(r)-1])
		}
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

/* func (a App) ServeHTTPConcurrent(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.String())
	v := req.URL.Query()
	if len(v["count"]) == 0 || len(v["offset"]) == 0 || len(v["count"][0]) == 0 || len(v["offset"][0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	count, err := strconv.Atoi(v["count"][0])
	if err != nil {
		log.Fatal(err)
	}
	offset, err := strconv.Atoi(v["offset"][0])
	if err != nil {
		log.Fatal(err)
	}
	ch, errCh := make(chan *ContentItem), make(chan error)
	for i := 0; len(result) < count; i++ {
		go a.getEntry(i, offset, req, ch, errCh)
	}

	var result []*ContentItem

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
} */

/* func (a App) getEntry(i, offset int, req *http.Request, ch chan *ContentItem, errCh chan error) {
	r, err := a.ContentClients[a.Config[i%len(a.Config)].Type].GetContent(req.URL.Hostname(), 1+offset)
	if err != nil {
		r, err := a.ContentClients[a.Config[i%len(a.Config)].Type].GetContent(req.URL.Hostname(), 1+offset)
		if err != nil {
			if a.Config[i%len(a.Config)].Fallback != nil {
				r, err := a.ContentClients[*a.Config[i%len(a.Config)].Fallback].GetContent(req.URL.Hostname(), 1+offset)
				if err != nil {
					log.Print(err)
					return
				}
				ch <- r[len(r)-1]
			} else {
				return
			}

		}
	} else {
		ch <- r[len(r)-1]
	}
} */
