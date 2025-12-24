package exa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/search" {
			t.Errorf("expected /search path, got %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("expected x-api-key test-key, got %s", r.Header.Get("x-api-key"))
		}

		resp := SearchResponse{
			RequestID: "test-req-id",
			Results: []Result{
				{ID: "1", URL: "https://example.com", Title: "Example"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	resp, err := client.Search(context.Background(), "test query", SearchOptions{})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.RequestID != "test-req-id" {
		t.Errorf("expected request ID test-req-id, got %s", resp.RequestID)
	}
	if len(resp.Results) != 1 || resp.Results[0].Title != "Example" {
		t.Errorf("unexpected results: %+v", resp.Results)
	}
}

func TestAnswer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := answerResponse{
			Answer: "The answer is 42",
			Citations: []Result{
				{URL: "https://hitchhikers.guide", Title: "Hitchhiker's Guide"},
			},
			RequestID: "ans-id",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	resp, err := client.Answer(context.Background(), "What is the answer?")
	if err != nil {
		t.Fatalf("Answer failed: %v", err)
	}

	if resp.Answer != "The answer is 42" {
		t.Errorf("unexpected answer: %s", resp.Answer)
	}
	if len(resp.Citations) != 1 || resp.Citations[0].Title != "Hitchhiker's Guide" {
		t.Errorf("unexpected citations: %+v", resp.Citations)
	}
}

func TestContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ContextResponse{
			Response:  "some code",
			RequestID: "ctx-id",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	resp, err := client.Context(context.Background(), "code query", "dynamic")
	if err != nil {
		t.Fatalf("Context failed: %v", err)
	}

	if resp.Response != "some code" {
		t.Errorf("unexpected response: %s", resp.Response)
	}
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "invalid request"}`))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	_, err := client.Search(context.Background(), "query", SearchOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedMsg := "api error: invalid request (status 400)"
	if err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %q", expectedMsg, err.Error())
	}
}
