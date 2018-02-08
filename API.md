# Tic Tac Toe Server API

A board is represented by a single string 9-characters long, representing each
of the 9 squares on the board. The string may contain the following characters:

 - `-`: The square is empty.
 - `x`: The square contains an x.
 - `o`: The square contains an o.

An example of how the 3x3 grid translates into a 9-character string:

```
--o
xxo --> "--oxxo-xo"
-xo
```

# Calls

## GET /games

Returns an array of games (active or inactive) by ID.

#### Example response

```json
{
    "games": [
        1, 2, 3, 4, 5
    ],
}
```

## POST /games

Creates a new game with two players, specified by player ID.

#### Example request

```
    player_x=<ID of the X player>&player_o=<ID of the O player>
```

#### Example response

```json
{
    "game_id": "<ID of the newly-created game>",
}
```

## GET /games/:id

Returns information on the game with the given ID.

Possible statuses are:

 - `X_MOVE` if the player playing X moves next.
 - `O_MOVE` if the player playing O moves next.
 - `X_WIN` if the game is over and X won.
 - `O_WIN` if the game is over and O won.
 - `DRAW` if the game is over and finished in a draw.

#### Example response

```json
{
    "id": 1,
    "player_x": 1,
    "player_o": 2,
    "board": "---xoo--x",
    "status": "IN_PROGRESS",
}
```

## POST /games/:id

Plays a move. A move is a single number representing the index on the board to
make the move. The board begins at index 0.

#### Example request

```
    player=<ID of player making the move>&move=1
```

## POST /players

Asks the server to create and return a new player ID.

#### Example response

```json
{
    "id": "1337",
}
```
