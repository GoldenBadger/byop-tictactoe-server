package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/goldenbadger/byop-tictactoe-server/internal/app/players"
	"github.com/goldenbadger/byop-tictactoe-server/internal/pkg/tictactoe"
)

type App struct {
	games   *tictactoe.Games
	players *players.Players
	router  *chi.Mux
}

func (a *App) Run() {
	a.initGames()
	a.initPlayers()
	a.initRouter()

	http.ListenAndServe(":8080", a.router)
}

func (a *App) initPlayers() {
	a.players = players.New()
}

func (a *App) initGames() {
	a.games = tictactoe.NewGames()
}

func (a *App) initRouter() {
	a.router = chi.NewRouter()
	a.router.Route("/players", func(r chi.Router) {
		r.Post("/", a.newPlayer)
	})
	a.router.Route("/games", func(r chi.Router) {
		r.Get("/", a.getGames)
		r.Post("/", a.createGame)
		r.Get("/{id}", a.getGameInfo)
		r.Post("/{id}", a.makeMove)
	})
}

func (a *App) newPlayer(w http.ResponseWriter, r *http.Request) {
	res := map[string]uint{"id": a.players.NewPlayer()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (a *App) getGames(w http.ResponseWriter, r *http.Request) {
	res := map[string][]uint{"games": a.games.GetGameIDs()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (a *App) createGame(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	x := r.PostFormValue("player_x")
	o := r.PostFormValue("player_o")
	if x == "" || o == "" {
		http.Error(w, "One or more players have a missing ID.", http.StatusBadRequest)
		return
	}
	xint, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID for player X.", http.StatusBadRequest)
		return
	}
	oint, err := strconv.ParseUint(o, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID for player O.", http.StatusBadRequest)
		return
	}
	if !a.players.PlayerExists(uint(xint)) || !a.players.PlayerExists(uint(oint)) {
		http.Error(w, "One or more players does not exist.", http.StatusBadRequest)
		return
	}

	res := map[string]uint{"game_id": a.games.NewGame(uint(xint), uint(oint))}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (a *App) getGameInfo(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID for game.", http.StatusBadRequest)
		return
	}
	game, err := a.games.GetGame(uint(id))
	if err != nil {
		http.Error(w, "Game does not exist.", http.StatusNotFound)
		return
	}

	res := map[string]interface{}{
		"id":       uint(id),
		"player_x": game.PlayerX,
		"player_o": game.PlayerO,
		"board":    flattenBoard(game.Board),
		"status":   gameStatusString(game.Status),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (a *App) makeMove(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID for game.", http.StatusBadRequest)
		return
	}
	_, err = a.games.GetGame(uint(id))
	if err != nil {
		http.Error(w, "Game does not exist.", http.StatusNotFound)
		return
	}

	r.ParseForm()
	playerstr := r.PostFormValue("player")
	movestr := r.PostFormValue("move")
	if playerstr == "" {
		http.Error(w, "Missing player ID.", http.StatusBadRequest)
		return
	}
	if movestr == "" {
		http.Error(w, "Missing move.", http.StatusBadRequest)
		return
	}
	player, err := strconv.ParseUint(playerstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}
	move, err := strconv.ParseUint(movestr, 10, 64)
	if err != nil || move < 0 || move > 8 {
		http.Error(w, "Invalid move.", http.StatusBadRequest)
		return
	}
	if !a.players.PlayerExists(uint(player)) {
		http.Error(w, "Player does not exist.", http.StatusBadRequest)
		return
	}

	row, col := unflattenMove(byte(move))
	err = a.games.MakeMove(uint(id), uint(player), row, col)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func flattenBoard(board [3][3]byte) string {
	var ret [9]byte
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			ret[3*row+col] = board[col][row]
		}
	}
	return string(ret[:9])
}

func unflattenMove(move byte) (byte, byte) {
	row := move % 3
	col := move / 3
	return row, col
}

func gameStatusString(status tictactoe.GameStatus) string {
	switch status {
	case tictactoe.XMove:
		return "X_MOVE"
	case tictactoe.OMove:
		return "O_MOVE"
	case tictactoe.XWin:
		return "X_WIN"
	case tictactoe.OWin:
		return "O_WIN"
	case tictactoe.Draw:
		return "DRAW"
	}
	return ""
}
