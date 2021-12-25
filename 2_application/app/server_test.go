package main


import (
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"
)

func TestForm (t *testing.T) {
    req, _ := http.NewRequest("GET", "/form", strings.NewReader("[]"))
    w := httptest.NewRecorder()
    formHandler(w, req)

    if w.Code != http.StatusMethodNotAllowed {
        t.Errorf("got HTTP status code %d, expected %d", w.Code, http.StatusMethodNotAllowed)
    }

    form := url.Values{}
    form.Add("name", "foo")
    form.Add("address", "bar")
    req, _ = http.NewRequest("POST", "/form", strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    w = httptest.NewRecorder()
    formHandler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("got HTTP status code %d, expected %d", w.Code, http.StatusOK)
    }
    if w.Body.String() != "Name = foo\nAddress = bar\n" {
        t.Errorf("got response %s, expected %s", w.Body.String(), "Name = foo\nAddress = bar\n")
    }
}

func TestHello (t *testing.T) {
    req, _ := http.NewRequest("POST", "/hello", strings.NewReader("[]"))
    w := httptest.NewRecorder()
    helloHandler(w, req)

    if w.Code != http.StatusMethodNotAllowed {
        t.Errorf("got HTTP status code %d, expected %d", w.Code, http.StatusMethodNotAllowed)
    }

    req, _ = http.NewRequest("GET", "/hello", strings.NewReader("[]"))
    w = httptest.NewRecorder()
    helloHandler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("got HTTP status code %d, expected %d", w.Code, http.StatusOK)
    }
    if w.Body.String() != "Hello!" {
        t.Errorf("got response %s, expected %s", w.Body.String(), "Hello!")
    }
 }
