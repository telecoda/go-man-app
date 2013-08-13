package models

import (
	"errors"
	"fmt"
	"github.com/telecoda/go-man-app/utils"
	"math"
	"math/rand"
	"strconv"
	"time"
)

/* this file contains the player specific function */

type PlayerType string

const (
	GoMan   PlayerType = "goman"
	GoGhost            = "goghost"
)

type PlayerState string

const (
	Alive    PlayerState = "alive"
	Dying                = "dying"
	Dead                 = "dead"
	Spawning             = "spawning"
)

type Player struct {
	Location      Point
	Id            string
	Type          PlayerType
	State         PlayerState
	Name          string
	cpuControlled bool
	Score         int
	Lives         int
}

const PLAYER_START_X = 13
const PLAYER_START_Y = 14

const GHOST_START_X = 13
const GHOST_START_Y = 10

//const MAX_GOMAN_PLAYERS int = 1
//const MAX_GOMAN_GHOSTS int = 4

const DEATH_WAIT_SECONDS = 1
const KILLED_GHOST_POINTS = 100
const KILLED_GOMAN_POINTS = 100

func (board *GameBoard) MovePlayer(player Player) error {

	// only allow moves when game playing
	if board.State != PlayingGame {
		return errors.New("Not playing, you cannot move now")
	}

	// check if player belongs to this game
	playerServerState := board.getPlayerFromServer(player.Id)

	if &playerServerState == nil {
		return errors.New("You are not a player in this game.")
	}

	// check player is alive
	if playerServerState.State != Alive {
		return errors.New("You are not alive so cannot move.")
	}

	// check move is valid
	if !isMoveValid(playerServerState.Location, player.Location) {
		return errors.New("Cheat, invalid move")
	}

	// check for wrap around to other side of board
	if player.Location.X < 0 {
		fmt.Println("Player wrap left")
		player.Location.X = (BOARD_WIDTH - 1)
	}

	if player.Location.X >= BOARD_WIDTH {
		fmt.Println("Player wrap right")
		player.Location.X = 0
	}

	cell := board.GetCellAtLocation(player.Location)

	switch cell {
	case WALL:
		return errors.New("Invalid move, you can't walk through walls")
	case PILL:
		if player.Type == GoMan {
			board.eatPillAtLocation(player.Location)
			playerServerState.Score += PILL_POINTS
			board.UpdatePillsRemaining()
			if board.PillsRemaining == 0 {
				board.gameWon()
			}
		}
		break
	case POWER_PILL:
		if player.Type == GoMan {
			board.eatPowerPillAtLocation(player.Location)
			playerServerState.Score += POWER_PILL_POINTS
			board.UpdatePillsRemaining()
			if board.PillsRemaining == 0 {
				board.gameWon()
			}
		}
		break

	}

	// update board with player's location
	playerServerState.Location.X = player.Location.X
	playerServerState.Location.Y = player.Location.Y

	// check for player collisions
	board.checkPlayerCollisions(playerServerState)

	return nil
}

func (board *GameBoard) checkPlayerCollisions(currentPlayer *Player) {

	for id, _ := range board.Players {
		player := board.Players[id]
		// if not same player
		// and not same type of player
		// and they are both alive
		if player.Id != currentPlayer.Id &&
			player.Type != currentPlayer.Type &&
			player.State == Alive &&
			currentPlayer.State == Alive {
			// check co-ords
			if player.Location.X == currentPlayer.Location.X &&
				player.Location.Y == currentPlayer.Location.Y {
				board.playersCollided(currentPlayer, player)
			}
		}
	}
	return

}

func (board *GameBoard) playersCollided(player1 *Player, player2 *Player) {

	fmt.Println("Player:", player1.Name, " hit Player:", player2.Name)

	if board.PowerPillsActive > 0 {
		fmt.Println("powerpill active")
		// if powerpill active ghost will die
		if player1.Type == GoMan {
			go board.ghostDeath(player2)

			player1.Score += KILLED_GHOST_POINTS

		} else {
			go board.ghostDeath(player1)

			player2.Score += KILLED_GHOST_POINTS

		}

	} else {
		fmt.Println("powerpill not active")
		// if powerpill not active goman will die
		if player1.Type == GoGhost {
			go board.gomanDeath(player2)
			fmt.Println("goman death player2")

			player1.Score += KILLED_GOMAN_POINTS

		} else {
			go board.gomanDeath(player1)
			fmt.Println("goman death player1")

			player2.Score += KILLED_GOMAN_POINTS
		}

	}

}

func (board *GameBoard) gomanDeath(goman *Player) {
	fmt.Println(goman.Name, " dying")
	goman.State = Dying

	time.Sleep(time.Duration(DEATH_WAIT_SECONDS) * time.Second)

	if goman.Lives > 0 {
		// lose a life
		goman.Lives--
		board.respawnGoMan(goman)

	} else {
		// dead for good
		goman.State = Dead

		liveGoMen := board.countLiveGoMen()

		if liveGoMen == 0 {
			// game over
			board.State = GameOver
		}
	}

}

func (board *GameBoard) ghostDeath(goGhost *Player) {
	fmt.Println(goGhost.Name, " dying")
	goGhost.State = Dying

	time.Sleep(time.Duration(DEATH_WAIT_SECONDS) * time.Second)

	// ghosts have unlimited lives, just spawn them again
	// and again, until all gomen are dead
	board.respawnGhost(goGhost)
}

func (board *GameBoard) respawnGoMan(goMan *Player) {

	fmt.Println(goMan.Name, " back to life")
	goMan.State = Alive
	goMan.Location = Point{PLAYER_START_X, PLAYER_START_Y}

}

func (board *GameBoard) respawnGhost(goGhost *Player) {

	fmt.Println(goGhost.Name, " back to life")
	goGhost.State = Alive
	goGhost.Location = Point{GHOST_START_X, GHOST_START_Y}

}

func (board *GameBoard) getPlayerFromServer(playerId string) *Player {

	// using a range to iterate through an array of objects
	// works on a copy of the object not a reference
	// to the object...!!

	return board.Players[playerId]

}

func isMoveValid(existingLocation Point, newLocation Point) bool {

	// player can only move in one direction at a time
	// player can only move one cell at a time

	distX := math.Abs(float64(existingLocation.X - newLocation.X))
	distY := math.Abs(float64(existingLocation.Y - newLocation.Y))

	// moved more than one cell
	if distX > 1 || distY > 1 {
		return false
	}

	// moved more than one direction
	if distX > 0 && distY > 0 {
		return false
	}

	// valid move
	return true
}

func (board *GameBoard) AddPlayer(newPlayer *Player) (*Player, error) {

	if newPlayer.Type != GoMan && newPlayer.Type != GoGhost {
		return nil, errors.New("Invalid player type")
	}

	if newPlayer.Type == GoGhost && board.getRemainingGoGhosts() <= 0 {
		return nil, errors.New("Cannot add anymore ghosts to game")
	}

	if newPlayer.Type == GoMan && board.getRemainingGoMen() <= 0 {
		return nil, errors.New("Cannot add anymore go-men to game")
	}

	if newPlayer.Type == GoMan {
		newPlayer.Location = Point{PLAYER_START_X, PLAYER_START_Y}

	} else {
		// its a ghost
		newPlayer.Location = Point{GHOST_START_X, GHOST_START_Y}
	}
	newPlayer.Id, _ = utils.GenUUID()
	newPlayer.State = Alive
	newPlayer.Score = 0
	newPlayer.Lives = 3
	newPlayer.State = Alive

	//board.Players = append(board.Players, *newPlayer)
	board.Players[newPlayer.Id] = newPlayer

	// if game is New this is first playing being added
	if board.State == NewGame {

		board.State = WaitingForPlayers
		// spawn process to wait for players
		go waitForPlayers(board.Id)

	}

	if board.getRemainingGoMen() == 0 &&
		board.getRemainingGoGhosts() == 0 {
		board.startGame()
	}

	return newPlayer, nil
}

func (board *GameBoard) getRemainingGoMen() int {
	goMenCount := board.countGoMen()
	return board.MaxGoMenAllowed - goMenCount
}

func (board *GameBoard) getRemainingGoGhosts() int {
	goGhostCount := board.countGhosts()
	return board.MaxGoGhostsAllowed - goGhostCount
}

func waitForPlayers(gameId string) {

	board, err := LoadGameBoard(gameId)

	if err != nil {
		// silently failed game may have already ended
		fmt.Println("Error loading board:", err)
		return
	}

	if board == nil {
		// silently failed game may have already ended
		fmt.Println("Error: board is empty:")
		return
	}

	// process sleeps until its time to wake up
	fmt.Println("New game: I am going to sleep whilst I wait for players on game:", gameId)
	fmt.Println("I will be asleep for ", board.WaitForPlayersSeconds, " seconds")

	time.Sleep(time.Duration(board.WaitForPlayersSeconds) * time.Second)

	fmt.Println("Yawn, I have awoken!")

	board.PopulateRemainingPlayers()

}

func (board *GameBoard) PopulateRemainingPlayers() {
	// this method is called after wait time has completed
	// its responsibility is to add any CPU controlled players
	// if necessary & then start the game

	if board.State != WaitingForPlayers {
		fmt.Println("Game no longer waiting for players, aborting")
		return
	}

	totalGoMen := board.countGoMen()

	totalGhosts := board.countGhosts()

	missingGoMen := board.MaxGoMenAllowed - totalGoMen

	for i := 0; i < missingGoMen; i++ {
		board.addCPUGoMan(i)
	}

	missingGhosts := board.MaxGoGhostsAllowed - totalGhosts

	for i := 0; i < missingGhosts; i++ {
		board.addCPUGhost(i)
	}

	board.SaveGameBoard()
}

func (board *GameBoard) addCPUGoMan(cpuId int) {

	fmt.Println("Adding goMan:", strconv.Itoa(cpuId))
	newCPUPlayer := new(Player)
	newCPUPlayer.Name = "CPUMAN-" + strconv.Itoa(cpuId)
	newCPUPlayer.Type = GoMan
	newCPUPlayer.Location.X = PLAYER_START_X
	newCPUPlayer.Location.Y = PLAYER_START_Y
	newCPUPlayer.cpuControlled = true
	board.AddPlayer(newCPUPlayer)

	fmt.Println("Added new GoMan:", newCPUPlayer.Name)
}

func (board *GameBoard) addCPUGhost(cpuId int) {

	fmt.Println("Adding goGhost:", strconv.Itoa(cpuId))
	newCPUPlayer := new(Player)
	newCPUPlayer.Name = "CPUGHOST-" + strconv.Itoa(cpuId)
	newCPUPlayer.Type = GoGhost
	newCPUPlayer.Location.X = GHOST_START_X
	newCPUPlayer.Location.Y = GHOST_START_Y
	newCPUPlayer.cpuControlled = true
	board.AddPlayer(newCPUPlayer)

	fmt.Println("Added new GoGhost:", newCPUPlayer.Name)
}

func (board *GameBoard) startGame() {

	board.State = PlayingGame
	board.PowerPillsActive = 0

	// submit go processes for each CPU controlled player
	for _, player := range board.Players {
		if player.cpuControlled {

			go playAsCPU(board.Id, player.Id)
		}
	}
}

func (board *GameBoard) gameWon() {

	board.State = GameWon
}

func playAsCPU(gameId string, playerId string) {

	/* this function will repeat until the current game ends */

	var gamePlaying = true

	for gamePlaying {

		// wait for 1/60 of a second
		//timer := time.NewTimer(time.Second / 60)

		// slow down enemy to 1/4 a move
		timer := time.NewTimer(time.Second / 4)
		<-timer.C

		// get current board state
		board, err := LoadGameBoard(gameId)

		if err != nil {
			fmt.Println("Error retrieving game, aborting.", err)
			return
		}

		if board.State == GameWon {
			// stop playing game is won
			return
		}

		if board.State == GameOver {
			// stop playing game is over
			return
		}

		player := board.getPlayerFromServer(playerId)

		if &player == nil {
			fmt.Println("Error player not found in game")
			return
		}

		movedPlayer := board.planBestMoveForPlayer(*player)

		err = board.MovePlayer(movedPlayer)

		if err != nil {
			//fmt.Println("Error moving player, carry on", err)
		}

	}

}

func (board *GameBoard) planBestMoveForPlayer(player Player) Player {

	if player.Type == GoMan {
		// plan best goman move
		if board.PowerPillsActive > 0 {
			// chase ghost
			player = board.goManChasesGhosts(player)
		} else {
			// avoid ghosts
			player = board.goManAvoidsGhosts(player)
		}
	} else {
		// plan best ghost mode
		if board.PowerPillsActive > 0 {
			// avoid goman
			player = board.ghostAvoidsGoman(player)
		} else {
			// chase goman
			player = board.ghostChasesGoman(player)
		}
	}
	return player
}

func (board *GameBoard) goManChasesGhosts(player Player) Player {

	return board.randomMovement(player)

}

func (board *GameBoard) goManAvoidsGhosts(player Player) Player {

	return board.randomMovement(player)

}

func (board *GameBoard) ghostChasesGoman(player Player) Player {

	return board.randomMovement(player)

}

func (board *GameBoard) ghostAvoidsGoman(player Player) Player {

	return board.randomMovement(player)

}

func (board *GameBoard) randomMovement(player Player) Player {
	// do some random moves for now
	whichWay := rand.Int() % 4

	switch whichWay {
	// left
	case 0:
		player.Location.X--
		break
	// right
	case 1:
		player.Location.X++
		break
	// up
	case 2:
		player.Location.Y--
		break
	// down
	case 3:
		player.Location.Y++
		break
	}

	return player
}

func (board *GameBoard) countGhosts() int {
	totalGhosts := 0
	for _, player := range board.Players {
		if player.Type == GoGhost {
			totalGhosts++
		}
	}

	return totalGhosts
}

func (board *GameBoard) countGoMen() int {
	totalGoMen := 0
	for _, player := range board.Players {
		if player.Type == GoMan {
			totalGoMen++
		}
	}

	return totalGoMen
}

func (board *GameBoard) countLiveGoMen() int {
	totalGoMen := 0
	for _, player := range board.Players {
		if player.Type == GoMan && player.State != Dead {
			totalGoMen++
		}
	}

	return totalGoMen
}
