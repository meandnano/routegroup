package routegroup_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-pkgz/routegroup"
)

// testMiddleware is simple middleware for testing purposes.
func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Test-Middleware", "true")
		next.ServeHTTP(w, r)
	})
}

func TestGroupMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	group := routegroup.New(mux)

	// apply middleware to the group
	group.Use(testMiddleware)

	// add a test handler
	group.Handle("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Test-Middleware"))
}

func TestGroupHandle(t *testing.T) {
	mux := http.NewServeMux()
	group := routegroup.New(mux)

	// add a test handler
	group.Handle("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGroupSet(t *testing.T) {
	mux := http.NewServeMux()
	group := routegroup.New(mux)

	// configure the group using Set
	group.Set(func(g *routegroup.Bundle) {
		g.Use(testMiddleware)
		g.Handle("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Test-Middleware"))
}

func TestGroupWithMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	group := routegroup.New(mux)

	// original group middleware
	group.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Original-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	})

	// new group with additional middleware
	newGroup := group.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-New-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	})

	// add a test handler to the new group
	newGroup.Handle("/with-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/with-test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Original-Middleware"))
	assert.Equal(t, "true", recorder.Header().Get("X-New-Middleware"))
}

func TestMount(t *testing.T) {
	mux := http.NewServeMux()
	basePath := "/api"
	group := routegroup.Mount(mux, basePath)

	group.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Mounted-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	})

	// add a test handler
	group.Handle("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the mounted handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, basePath+"/test", nil)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Mounted-Middleware"))
}