package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeLink(t *testing.T) {
	t.Run("HomePage", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		HomeLink(response, request)
		got := response.Body.String()
		want := "HomePage !!"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}

func TestGetOneUser(t *testing.T) {
	t.Run("Return one user", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		response := httptest.NewRecorder()

		getOneUser(response, request)
		got := response.Body.String()
		want := ""

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

	})
}

func TestGetAllUsers(t *testing.T) {

}
