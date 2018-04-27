package main

import (
	"fmt"
)

// The following constants denote various limitations in the engine
const roomZero = 0
const roomInventory = -1
const verbAuto = 1
const nounAny = 1
const counters = 8
const flags = 32
const alternateRooms = 6

// Contains all the data from a game file
type gameStaticData struct {
	filename          string
	advVariable       map[string]int
	action            [][]int
	verb              []string
	noun              []string
	roomDirection     []map[int]int
	roomDescription   []string
	message           []string
	itemDescription   []string
	itemNoun          []string
	itemStartLocation []int
	actionComment     []string
	treasureItem      []int
}

// Contains all the data for a player session
type gameDynamicData struct {
	currentRoom      int
	itemLocation     []int
	bitFlag          []bool
	alternateCounter []int
	alternateRoom    []int
}

func main() {
	var loadedGame gameStaticData
	loadedGame.filename = "adv01.dat"
	loadData(&loadedGame)

	var gameInstance gameDynamicData

	initGame(loadedGame.advVariable, loadedGame.itemStartLocation, &gameInstance)

	// Print some stuff, to get the program to compile without being complete...
	fmt.Println(gameInstance.currentRoom)
	fmt.Println(gameInstance.itemLocation[0])
	fmt.Println(gameInstance.bitFlag[0])
	fmt.Println(gameInstance.alternateCounter[0])
	fmt.Println(gameInstance.alternateRoom[0])

	fmt.Println(loadedGame.advVariable["wordLength"])
	fmt.Println(loadedGame.action[0][0])
	fmt.Println(loadedGame.verb[0])
	fmt.Println(loadedGame.noun[0])
	fmt.Println(loadedGame.roomDirection[0][0])
	fmt.Println(loadedGame.roomDescription[0])
	fmt.Println(loadedGame.message[0])
	fmt.Println(loadedGame.itemDescription[0])
	fmt.Println(loadedGame.itemNoun[0])
	fmt.Println(loadedGame.itemStartLocation[0])
	fmt.Println(loadedGame.actionComment[0])
	fmt.Println(loadedGame.treasureItem[0])
}

func initGame(advVariable map[string]int, itemStartLocation []int, instance *gameDynamicData) {
	itemLocation := make([]int, len(itemStartLocation))
	copy(itemLocation, itemStartLocation)
	bitFlag := make([]bool, flags)
	alternateCounter := make([]int, counters)
	alternateRoom := make([]int, alternateRooms)

	instance.currentRoom = advVariable["startingRoom"]
	instance.itemLocation = itemLocation
	instance.bitFlag = bitFlag
	instance.alternateCounter = alternateCounter
	instance.alternateRoom = alternateRoom
	return
}
