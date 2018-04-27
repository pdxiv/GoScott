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

// Not used yet. Will make it possible to have multiple games loaded
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

// Not used yet. Makes it possible to run many players at the same time
type gameDynamicData struct {
	currentRoom      int
	itemLocation     []int
	bitFlag          []bool
	alternateCounter []int
	alternateRoom    []int
}

func main() {
	// Load a new game data file
	advVariable,
		action,
		verb,
		noun,
		roomDirection,
		roomDescription,
		message,
		itemDescription,
		itemNoun,
		itemStartLocation,
		actionComment,
		treasureItem := loadData("adv01.dat")

	// Initialize variable game data
	var currentRoom int
	var itemLocation []int
	var bitFlag []bool
	var alternateCounter []int
	var alternateRoom []int

	currentRoom, itemLocation, bitFlag, alternateCounter, alternateRoom = initGame(advVariable, itemStartLocation)

	// Print some stuff, to get the program to compile without being complete...
	fmt.Println(currentRoom)
	fmt.Println(itemLocation[0])
	fmt.Println(bitFlag[0])
	fmt.Println(alternateCounter[0])
	fmt.Println(alternateRoom[0])
	fmt.Println(advVariable["wordLength"])
	fmt.Println(action[0][0])
	fmt.Println(verb[0])
	fmt.Println(noun[0])
	fmt.Println(roomDirection[0][0])
	fmt.Println(roomDescription[0])
	fmt.Println(message[0])
	fmt.Println(itemDescription[0])
	fmt.Println(itemNoun[0])
	fmt.Println(itemStartLocation[0])
	fmt.Println(actionComment[0])
	fmt.Println(treasureItem[0])
}

func initGame(advVariable map[string]int, itemStartLocation []int) (int, []int, []bool, []int, []int) {
	itemLocation := make([]int, len(itemStartLocation))
	copy(itemLocation, itemStartLocation)
	bitFlag := make([]bool, flags)
	alternateCounter := make([]int, counters)
	alternateRoom := make([]int, alternateRooms)
	return advVariable["startingRoom"], itemLocation, bitFlag, alternateCounter, alternateRoom
}
