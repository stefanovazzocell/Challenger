package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/stefanovazzocell/Challenger/provider"
)

type Handler struct {
	lock        *sync.Mutex
	verifiers   map[string][]byte
	static_html []byte
	static_js   []byte
}

func (h *Handler) handleDemo(w http.ResponseWriter, r *http.Request) {
	log.Println("Got / request")
	w.Header().Set("Content-Type", "text/html")
	w.Write(h.static_html)
}
func (h *Handler) handleJs(w http.ResponseWriter, r *http.Request) {
	log.Println("Got /challenge.js request")
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(h.static_js)
}
func (h *Handler) handleChallenge(w http.ResponseWriter, r *http.Request) {
	log.Println("Got /challenge request")
	challenge, verifier := provider.NewChallenge(100, 15000)
	log.Printf("[challenge info] id: %q", challenge.Id)
	log.Printf("[challenge info] query: %v", challenge.Query)
	log.Printf("[challenge info] challenge: %v", challenge.Challenge)
	log.Printf("[challenge info] expected: %v", verifier.Expected)

	h.lock.Lock()
	h.verifiers[verifier.Id] = verifier.Expected
	h.lock.Unlock()

	response, err := json.Marshal(challenge)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func (h *Handler) handleSolve(w http.ResponseWriter, r *http.Request) {
	log.Println("Got /solve request")
	solution := provider.Solution{}
	if err := json.NewDecoder(r.Body).Decode(&solution); err != nil {
		panic(err)
	}
	h.lock.Lock()
	verifierData, found := h.verifiers[solution.Id]
	if !found {
		w.WriteHeader(404)
		return
	}
	delete(h.verifiers, solution.Id)
	h.lock.Unlock()

	verifier := provider.Verifier{
		Id:       solution.Id,
		Expected: verifierData,
	}
	log.Printf("[solution info] response: %v", solution.Response)
	if verifier.Verify(solution) {
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(400)
}

func NewHandler() Handler {
	log.Printf("> Loading client/demo.html...")
	static_html, err := os.ReadFile("./client/demo.html")
	if err != nil {
		panic(err)
	}
	log.Printf("> Loading client/challenge.js...")
	static_js, err := os.ReadFile("./client/challenge.js")
	if err != nil {
		panic(err)
	}
	return Handler{
		lock:        &sync.Mutex{},
		verifiers:   map[string][]byte{},
		static_html: static_html,
		static_js:   static_js,
	}
}
