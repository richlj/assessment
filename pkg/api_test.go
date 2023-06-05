package pkg_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/richlj/sliide/go-challenge-master/pkg"
)

var (
	SimpleContentRequest = httptest.NewRequest("GET", "/?offset=0&count=5", nil)
	OffsetContentRequest = httptest.NewRequest("GET", "/?offset=5&count=5", nil)
)

func runRequest(t *testing.T, srv http.Handler, r *http.Request) (content []*pkg.ContentItem) {
	response := httptest.NewRecorder()
	srv.ServeHTTP(response, r)

	if response.Code != 200 {
		t.Fatalf("Response code is %d, want 200", response.Code)
		return
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&content); err != nil {
		t.Fatalf("couldn't decode Response json: %v", err)
	}

	return content
}

// app gets initialised with configuration.
// as an example we've added 3 providers and a default configuration
var app = pkg.App{
	ContentClients: map[pkg.Provider]pkg.Client{
		pkg.Provider1: pkg.SampleContentProvider{Source: pkg.Provider1},
		pkg.Provider2: pkg.SampleContentProvider{Source: pkg.Provider2},
		pkg.Provider3: pkg.SampleContentProvider{Source: pkg.Provider3},
	},
	Config: pkg.DefaultConfig,
}

func TestResponseCount(t *testing.T) {

	content := runRequest(t, app, SimpleContentRequest)

	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}

}

func TestResponseOrder(t *testing.T) {
	content := runRequest(t, app, SimpleContentRequest)
	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}
	for i, item := range content {
		if pkg.Provider(item.Source) != pkg.DefaultConfig[i%len(pkg.DefaultConfig)].Type {
			t.Errorf(
				"Position %d: Got Provider %v instead of Provider %v",
				i, item.Source, pkg.DefaultConfig[i].Type,
			)
		}
	}
}

func TestOffsetResponseOrder(t *testing.T) {
	content := runRequest(t, app, OffsetContentRequest)
	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}
	for j, item := range content {
		i := j + 5
		if i > len(pkg.DefaultConfig)-1 {
			break
		}
		if pkg.Provider(item.Source) != pkg.DefaultConfig[i%len(pkg.DefaultConfig)].Type {
			t.Errorf(
				"Position %d: Got Provider %v instead of Provider %v",
				i, item.Source, pkg.DefaultConfig[i].Type,
			)
		}
	}
}

/* func TestRespectsFallback(t *testing.T) {
	content := runRequest(t, app, SimpleContentRequest)
	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}
	for j, item := range content {
		i := j + 5
		if i > len(pkg.DefaultConfig)-1 {
			break
		}
		if pkg.Provider(item.Source) != pkg.DefaultConfig[i%len(pkg.DefaultConfig)].Type {
			t.Errorf(
				"Position %d: Got Provider %v instead of Provider %v",
				i, item.Source, pkg.DefaultConfig[i].Type,
			)
		}
	}
} */
