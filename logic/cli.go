package logic

import (
	"fmt"
	"os"
)

// Print the Board to the terminal
func (b GameData) printBoard() {
	for i := 0; i < len(b.board); i++ {
		if i%10 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%v	", b.board[i])
	}
	fmt.Printf("\n")
}

// Process a Human Move in the terminal
func (g *Game) humanMove(legal map[int]bool) {
	var move int
	fmt.Printf("%v, enter your move. Choose any of the above legal moves: ", g.state.currentTurn.name)
	fmt.Scanln(&move)
	elem, ok := legal[move]
	if !ok {
		fmt.Printf("%v INVALID MOVE!\n", elem)
	} else {
		g.flip(move)
	}
}

// Play Othello in the terminal!
func PlayGame() {
	gameType := setGameType()
	switch gameType {
	case 1:
		HumanGame()
	case 2:
		BotGame("Randy")
	case 3:
		BotGame("Max")
	case 4:
		displayHelpMsg()
	default:
		fmt.Println("EXITING....")
	}
	os.Exit(0)
}

func setGameType() int {
	var input int
	fmt.Printf("Enter 1 for to play against a human\n")
	fmt.Printf("Enter 2 to play against Randy, a bot that selects moves at random\n")
	fmt.Printf("Enter 3 to play against Max, a minimax bot\n")
	fmt.Printf("Enter 4 for the rules of Othello\n")
	fmt.Printf("GAME TYPE: ")
	fmt.Scanln(&input)
	return input
}

// Initialize game against two Humans
func HumanGame() {
	p := setPlayersHuman()
	g := NewGame(p)
	fmt.Printf("PLAYERS: (%s: %s), (%s: %s)\n", p.o.name, p.o.piece, p.x.name, p.x.piece)
	for !g.gameOver() {
		g.state.printBoard()
		legal := g.gameStatus()
		g.humanMove(legal)
	}
	result(g)
}

// Set up the players for a game against two Humans in the terminal.
func setPlayersHuman() *PlayerList {
	var name, name2 string
	fmt.Printf("Enter Player 1's name. They will be the \"O\" pieces: ")
	fmt.Scanln(&name)
	fmt.Printf("Enter Player 2's name. They will be the \"X\" pieces: ")
	fmt.Scanln(&name2)
	p1, p2 := NewPlayer(name, "O"), NewPlayer(name2, "X")
	players := NewPlayerList(p2, p1)
	return players
}

// Set up the players for a game against a Bot and a Human in the terminal.
func BotGame(bot string) {
	p := setPlayersBot(bot)
	g := NewGame(p)
	fmt.Printf("PLAYERS: (%s: %s), (%s: %s)\n", p.o.name, p.o.piece, p.x.name, p.x.piece)
	for !g.gameOver() {
		g.state.printBoard()
		legal := g.gameStatus()
		if len(legal) < 1 {
			g.flip(-1)
		}
		var move int
		if g.state.currentTurn.piece == "O" {
			fmt.Printf("%s is thinking on a move...\n", bot)
			switch bot {
			case "Randy":
				move = RandyMove(legal)
			case "Max":
				// TODO move = MaxMove(g, 3)
				move = MaxMove(g, 3)
			}
			fmt.Printf("%s chose: %d\n", bot, move)
			g.flip(move)
		} else {
			g.humanMove(legal)
		}
	}
	result(g)
}

// Set up the players for a game against a Human and a Bot in the terminal.
func setPlayersBot(bot string) *PlayerList {
	var name string
	fmt.Printf("Initializing %s...complete! %s will play as `O`\n", bot, bot)
	fmt.Printf("Enter your name: ")
	fmt.Scanln(&name)
	fmt.Printf("\nWelcome, %s. Good luck against %s!\n", name, bot)
	p1, p2 := NewPlayer(bot, "O"), NewPlayer(name, "X")
	players := NewPlayerList(p2, p1)
	return players
}

func displayHelpMsg() {
	fmt.Println("TODO: PRINT OTHELLO RULES")
}

// Result screen for CLI
func result(g *Game) {
	g.UpdateScore()
	fmt.Println("------GAME OVER!------")
	g.state.printBoard()
	fmt.Printf("FINAL SCORE:	O: %d	X: %d\n", g.state.score.o, g.state.score.x)
	switch {
	case g.state.score.o > g.state.score.x:
		winner := g.state.players.o
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.name)
	case g.state.score.o < g.state.score.x:
		winner := g.state.players.x
		fmt.Printf("WINNER: %v. CONGRATULATIONS!\n", winner.name)
	default:
		fmt.Printf("RESULT: TIE GAME!\n")
	}
}

// Updates the current status of a Game. Returns a map of the available moves for
// the Player whose turn it is.
func (g *Game) gameStatus() map[int]bool {
	g.UpdateScore()
	legal := g.availableMoves()

	fmt.Printf("CURRENT TURN: %v, %s\n", g.state.currentTurn.name, g.state.currentTurn.piece)
	fmt.Printf("SCORE:	O: %d X: %d\n", g.state.score.o, g.state.score.x)
	fmt.Printf("LEGAL MOVES FOR %s: %v\n", g.state.currentTurn.piece, legal)
	return legal
}
