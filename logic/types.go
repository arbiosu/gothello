package logic

// Represents a Player in a Game of Othello
type Player struct {
	name  string
	piece string
}

// The two players facing off in a Game. Access by piece representation
type PlayerList struct {
	x *Player
	o *Player
}

// Represents a the game data collected in a Game of Othello
type GameData struct {
	board       [100]string
	currentTurn *Player
	score       *score
	players     *PlayerList
	directions  [8]int
}

// Represents a Game of Othello
type Game struct {
	state   *GameData
	history *history
}

// Captures Board state
type history struct {
	moves []GameData
}

type score struct {
	x int
	o int
}
