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

type PlayerMove struct {
	GameId          string
	PlayerToMove    Player
	ResponseChannel chan (PlayerMoveResponse)
}

type PlayerMoveResponse struct {
	Board GameBoard
	Error error
}

func NewPlayerMoveResponse(board GameBoard, err error) PlayerMoveResponse {
	response := new(PlayerMoveResponse)
	response.Board = board
	response.Error = err
	return *response
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

var GameChannels map[string]chan PlayerMove

func init() {

	GameChannels = make(map[string]chan PlayerMove)

}

// there is one channel per game used to receive and process player moves
func createNewGameChannel() chan PlayerMove {

	playerRequestChannel := make(chan PlayerMove)

	go func() {
		// process concurrent player updates
		for {

			playerMoveRequest := <-playerRequestChannel
			fmt.Println("Request received for game :", playerMoveRequest.GameId)
			fmt.Println("Request received for player :", playerMoveRequest.PlayerToMove.Name)
			ConcurrentMovePlayer(playerMoveRequest)
		}

	}()

	return playerRequestChannel

}

func ConcurrentMovePlayer(playerMove PlayerMove) {

	var response PlayerMoveResponse

	// fetch board
	fmt.Println("1. load game board")
	board, err := LoadGameBoard(playerMove.GameId)

	if err != nil {
		response = NewPlayerMoveResponse(*board, err)
		playerMove.ResponseChannel <- response
		return
	}
	// only allow moves when game playing
	if board.State != PlayingGame {
		response = NewPlayerMoveResponse(*board, errors.New("Not playing, you cannot move now"))
		playerMove.ResponseChannel <- response
		return
	}

	fmt.Println("2. fetch current player from board")
	// check if player belongs to this game
	playerServerState := board.getPlayer(playerMove.PlayerToMove.Id)

	if &playerServerState == nil {
		response = NewPlayerMoveResponse(*board, errors.New("You are not a player in this game."))
		playerMove.ResponseChannel <- response
		return
	}

	// check player is alive
	if playerServerState.State != Alive {
		response = NewPlayerMoveResponse(*board, errors.New("You are not alive so cannot move."))
		playerMove.ResponseChannel <- response
		return
	}

	// check move is valid
	if !isMoveValid(playerServerState.Location, playerMove.PlayerToMove.Location) {
		response = NewPlayerMoveResponse(*board, errors.New("Cheat, invalid move"))
		playerMove.ResponseChannel <- response
		return
	}

	// check for wrap around to other side of board
	if playerMove.PlayerToMove.Location.X < 0 {
		fmt.Println("Player wrap left")
		playerMove.PlayerToMove.Location.X = (BOARD_WIDTH - 1)
	}

	if playerMove.PlayerToMove.Location.X >= BOARD_WIDTH {
		fmt.Println("Player wrap right")
		playerMove.PlayerToMove.Location.X = 0
	}

	fmt.Println("3. check cell moving to")
	cell := board.GetCellAtLocation(playerMove.PlayerToMove.Location)

	switch cell {
	case WALL:
		response = NewPlayerMoveResponse(*board, errors.New("Invalid move, you can't walk through walls"))
		playerMove.ResponseChannel <- response
		return
	case PILL:
		if playerMove.PlayerToMove.Type == GoMan {
			board.eatPillAtLocation(playerMove.PlayerToMove.Location)
			playerServerState.Score += PILL_POINTS
			board.UpdatePillsRemaining()
			if board.PillsRemaining == 0 {
				board.gameWon()
			}
		}
		break
	case POWER_PILL:
		if playerMove.PlayerToMove.Type == GoMan {
			board.eatPowerPillAtLocation(playerMove.PlayerToMove.Location)
			playerServerState.Score += POWER_PILL_POINTS
			board.UpdatePillsRemaining()
			if board.PillsRemaining == 0 {
				board.gameWon()
			}
		}
		break

	}

	fmt.Println("4. move player")
	// update board with player's location
	//fmt.Println("Moving player:", playerServerState.Name)
	playerServerState.Location.X = playerMove.PlayerToMove.Location.X
	playerServerState.Location.Y = playerMove.PlayerToMove.Location.Y

	// check for player collisions
	fmt.Println("5. check for player collisions")
	board.checkPlayerCollisions(playerServerState)

	// save change to board
	/*err = board.SaveGameBoard()
	if err != nil {
		response = NewPlayerMoveResponse(*board, err)
		playerMove.ResponseChannel <- response
		return
	}*/

	fmt.Println("6. return response")

	// no errors
	response = NewPlayerMoveResponse(*board, nil)
	playerMove.ResponseChannel <- response
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

func (board *GameBoard) getPlayer(playerId string) *Player {

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

			go concurrentPlayAsCPU(board.Id, player.Id)
		}
	}
}

func (board *GameBoard) gameWon() {

	board.State = GameWon
}

func concurrentPlayAsCPU(gameId string, playerId string) {

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

		player := board.getPlayer(playerId)

		if &player == nil {
			fmt.Println("Error player not found in game")
			return
		}

		movedPlayer := board.planBestMoveForPlayer(*player)

		playerMoveRequest := new(PlayerMove)
		playerMoveRequest.GameId = gameId
		playerMoveRequest.PlayerToMove = movedPlayer

		playerResponseChannel := make(chan PlayerMoveResponse)

		playerMoveRequest.ResponseChannel = playerResponseChannel

		// send request to game channel
		var gameRequestChannel chan PlayerMove
		gameRequestChannel = GameChannels[gameId]

		if gameRequestChannel == nil {
			fmt.Println("Error no request channel found for game")
			return
		}

		// send
		gameRequestChannel <- *playerMoveRequest

		// receive response
		var playerMoveResponse PlayerMoveResponse

		playerMoveResponse = <-playerResponseChannel

		if playerMoveResponse.Error != nil {
			fmt.Println("Error moving player, carry on", playerMoveResponse.Error)
			fmt.Println("Player in error:", playerMoveRequest.PlayerToMove.Name)
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
