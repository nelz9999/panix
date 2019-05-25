// Package panix provides a modular HTTP panic recovery handler.
package panix

import "net/http"

// Observer is an interface that describes types that will be
// notified when a panic is recovered, possibly for recording
// the panic in a different system.
type Observer interface {
	Observe(parg interface{}, req *http.Request)
}

// ObserverFunc is an adapter for making a function an Observer
type ObserverFunc func(parg interface{}, req *http.Request)

// Observe conforms to the Observer interface
func (fn ObserverFunc) Observe(parg interface{}, req *http.Request) {
	fn(parg, req)
}

// Responder is an interface that describes the type that will be
// relied upon for expressing whatever is the expected client
// output when a panic has been recovered.
type Responder interface {
	Respond(parg interface{}, w http.ResponseWriter, req *http.Request)
}

// ResponderFunc is an adapter for making a function a Responder
type ResponderFunc func(parg interface{}, w http.ResponseWriter, req *http.Request)

// Respond conforms to the Responder interface
func (fn ResponderFunc) Respond(parg interface{}, w http.ResponseWriter, req *http.Request) {
	fn(parg, w, req)
}

// New adds another handler in the chain that will catch and recover from
// any panics that escape from further down the chain. One Responder is
// expected to be provided, along with zero to many Observers.
func New(next http.Handler, resp Responder, obs ...Observer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if parg := recover(); parg != nil {
				for _, o := range obs {
					o.Observe(parg, r)
				}
				resp.Respond(parg, w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
