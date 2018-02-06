package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/goldenbadger/byop-tictactoe-server/internal/app/players"
)

type App struct {
	players *players.Players
	router  *chi.Mux
}

func (a *App) Run() {
	a.initPlayers()
	a.initRouter()

	http.ListenAndServe(":8080", a.router)
}

func (a *App) initPlayers() {
	a.players = players.New()
}

func (a *App) initRouter() {
	a.router = chi.NewRouter()
	a.router.Route("/players", func(r chi.Router) {
		r.Post("/", a.newPlayer)
	})
}

func (a *App) newPlayer(w http.ResponseWriter, r *http.Request) {
	res := struct {
		ID uint
	}{
		a.players.NewPlayer(),
	}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
