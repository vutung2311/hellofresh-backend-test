package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vutung2311-golang-test/pkg/router"
)

func TestRouter_New(t *testing.T) {
	r := router.New(nil)
	if r == nil {
		t.Error("expect returned router instance should not be nil")
	}
}

func TestRouter_AddRoute(t *testing.T) {
	counter1 := 0
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter1++
			next.ServeHTTP(w, r)
		})
	}
	counter2 := 0
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter2++
			next.ServeHTTP(w, r)
		})
	}
	r := router.New([]func(next http.Handler) http.Handler{middleware1, middleware2})

	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`foo`))
	})
	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`bar`))
	})

	r.AddRoute("/abc", handler1)
	r.AddRoute("/def", handler2)

	request1 := httptest.NewRequest("GET", "http://foo.com/abc", nil)
	response1 := httptest.NewRecorder()
	r.ServeHTTP(response1, request1)
	if counter1 != 1 {
		t.Error("expect middleware1 should be called")
	}
	if counter2 != 1 {
		t.Error("expect middleware2 should be called")
	}
	if !strings.Contains(response1.Body.String(), "foo") {
		t.Error("expect handle1 should be called")
	}

	request2 := httptest.NewRequest("GET", "http://bar.com/def", nil)
	response2 := httptest.NewRecorder()
	r.ServeHTTP(response2, request2)
	if counter1 != 2 {
		t.Error("expect middleware1 should be called")
	}
	if counter2 != 2 {
		t.Error("expect middleware2 should be called")
	}
	if !strings.Contains(response2.Body.String(), "bar") {
		t.Error("expect handle2 should be called")
	}
}
