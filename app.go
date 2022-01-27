package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	validator "gopkg.in/validator.v2"
)

// App encaosulates Env,Route and middleware

type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	Config      *Env
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	ShortLink string `json:"shortlink"`
}

func (a *App) Initialize(e *Env) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Config = e
	a.Router = mux.NewRouter()
	a.Middlewares = &Middleware{}
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	m := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)
	// a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	// a.Router.HandleFunc("/api/info", a.getShoirtlinkInfo).Methods("GET")
	// a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortlink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShoirtlinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("parse parameters failed%v", r.Body),
		})
		return
	}
	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("validate parameters failed%v", req),
		})
		return
	}
	defer r.Body.Close()
	s, err := a.Config.S.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusCreated, shortlinkResp{ShortLink: s})
	}
}

func (a *App) getShoirtlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	fmt.Printf("%s\n", s)
	d, err := a.Config.S.ShortlinkInfo(s)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusOK, d)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
	u, err := a.Config.S.Unshorten(vars["shortlink"])
	fmt.Printf("%s\n", u)

	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
