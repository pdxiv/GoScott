package main

import (
	"fmt"
	"os"
	"strings"
)

// The following constants denote various limitations in the interpreter
const roomZero = 0
const roomInventory = -1
const verbAuto = 1
const nounAny = 1
const counters = 8
const flags = 32
const alternateRooms = 6

// Contains all the data from a game file
type gameStaticData struct {
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
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}
	filename := os.Args[1]
	loadData(filename, &loadedGame)
	var gameInstance gameDynamicData
	initGame(loadedGame.advVariable, loadedGame.itemStartLocation, &gameInstance)
	condition := initConditions()
	command := initCommands()

	// printSomeGameData(&gameInstance, &loadedGame)
	showRoomDescription(&gameInstance, &loadedGame)
	// for {
	getConsoleInput(&loadedGame)
	// }
	fmt.Println("DEBUG: condition:", condition[0](1))
	command[0](1)
}

func showRoomDescription(gameInstance *gameDynamicData, loadedGame *gameStaticData) {
	// gameInstance.currentRoom = 6 // debug
	fmt.Println(loadedGame.roomDescription[gameInstance.currentRoom])
	var objectInRoom []string
	for index, element := range gameInstance.itemLocation {
		if element == gameInstance.currentRoom {
			objectInRoom = append(objectInRoom, loadedGame.itemDescription[index])
		}
	}
	if len(objectInRoom) > 0 {
		fmt.Println("Visible items here:")
		fmt.Println(strings.Join(objectInRoom, ", "))
	}
	return
}

func printSomeGameData(gameInstance *gameDynamicData, loadedGame *gameStaticData) {
	fmt.Println("currentRoom", gameInstance.currentRoom)
	fmt.Println("itemLocation[0]", gameInstance.itemLocation[0])
	fmt.Println("bitFlag[0]", gameInstance.bitFlag[0])
	fmt.Println("alternateCounter[0]", gameInstance.alternateCounter[0])
	fmt.Println("alternateRoom[0]", gameInstance.alternateRoom[0])

	fmt.Println("advVariable[\"wordLength\"]", loadedGame.advVariable["wordLength"])
	fmt.Println("action[0][0]", loadedGame.action[0][0])
	fmt.Println("verb[0]", loadedGame.verb[0])
	fmt.Println("noun[0]", loadedGame.noun[0])
	fmt.Println("roomDirection[0][0]", loadedGame.roomDirection[0][0])
	fmt.Println("roomDescription[0]", loadedGame.roomDescription[0])
	fmt.Println("message[0]", loadedGame.message[0])
	fmt.Println("itemDescription[0]", loadedGame.itemDescription[0])
	fmt.Println("itemNoun[0]", loadedGame.itemNoun[0])
	fmt.Println("itemStartLocation[0]", loadedGame.itemStartLocation[0])
	fmt.Println("actionComment[0]", loadedGame.actionComment[0])
	fmt.Println("treasureItem[0]", loadedGame.treasureItem[0])
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
