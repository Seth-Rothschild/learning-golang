package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	indexHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Errorf("Expected %v, got %v", 200, res.StatusCode)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "<h1>TODO List</h1>"
	if !strings.Contains(string(data), expected) {
		t.Errorf("expected %v in %v", expected, string(data))
	}
}

func TestAddHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/add", nil)
	w := httptest.NewRecorder()
	addHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != 303 {
		t.Errorf("Expected %v, got %v", 303, res.StatusCode)
	}
}

func TestToggleHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/add", nil)
	w := httptest.NewRecorder()
	toggleHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != 303 {
		t.Errorf("Expected %v, got %v", 303, res.StatusCode)
	}
}

func BenchmarkIndex(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		indexHandler(w, req)

	}
}

func BenchmarkAdd(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		addHandler(w, req)

	}
}

func BenchmarkToggle(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		addHandler(w, req)

	}
}
