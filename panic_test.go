package panix_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nelz9999/panix"
)

func TestPanic(t *testing.T) {
	// Arrange
	now := time.Now().UnixNano()
	content := fmt.Sprintf("panic-body-%d", now)
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(content)
	})

	stat := int(500 + now%17)

	marks := make([]bool, 3)
	h = panix.New(
		h,
		panix.ResponderFunc(func(parg interface{}, w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(stat)
			fmt.Fprintf(w, "PANICED: %v\n", parg)
		}),
		panix.ObserverFunc(func(parg interface{}, req *http.Request) {
			marks[0] = true
		}),
		panix.ObserverFunc(func(parg interface{}, req *http.Request) {
			marks[1] = true
		}),
		panix.ObserverFunc(func(parg interface{}, req *http.Request) {
			marks[2] = true
		}),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/nada", nil)

	// Act
	h.ServeHTTP(w, r)

	//Assert
	for _, mark := range marks {
		if !mark {
			t.Errorf("expected all trues: %v\n", marks)
		}
	}

	resp := w.Result()
	if resp.StatusCode != stat {
		t.Errorf("expected %d; got %d\n", stat, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), content) {
		t.Errorf("expected %q in %q\n", content, body)
	}
}
