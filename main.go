package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	reallyBigNumber     int = 32767
	actionCommandOffset int = 6
	actionEntries       int = 8
	commandCodeDivisor  int = 150
	commandsInAction    int = 4
	conditionDivisor    int = 20
	conditions          int = 5
	counterTimeLimit    int = 8
	directionNouns      int = 6

	falseValue             int = 0
	flagLampEmpty          int = 16
	flagNight              int = 15
	lightSourceID          int = 9
	lightWarningThreshold  int = 25
	message1End            int = 51
	message2Start          int = 102
	parConditionCode       int = 0
	percentUnits           int = 100
	roomInventory          int = -1
	roomStore              int = 0
	verbCarry              int = 10
	verbDrop               int = 18
	verbGo                 int = 1
	alternateRoomRegisters int = 6
	alternateCounters      int = 9
	statusFlags            int = 32
	minimumCounterValue    int = -1
	valuesIn16Bits         int = 65536
	prngPrm                int = 75
	prngPrime              int = 65537
)

var directionNounText = []string{"NORTH", "SOUTH", "EAST", "WEST", "UP", "DOWN"}

type GameState struct {
	gameFile                string
	keyboardInput2          string
	carriedObjects          int
	commandOrDisplayMessage int
	commandParameter        int
	commandParameterIndex   int
	contFlag                bool
	counterRegister         int
	currentRoom             int
	globalNoun              string
	maxObjectsCarried       int
	numberOfActions         int
	numberOfMessages        int
	numberOfObjects         int
	numberOfRooms           int
	numberOfTreasures       int
	numberOfWords           int
	startingRoom            int
	storedTreasures         int
	timeLimit               int
	treasureRoomID          int
	wordLength              int
	adventureVersion        int
	adventureNumber         int

	alternateCounter []int
	alternateRoom    []int

	objectDescription      []string
	message                []string
	extractedInputWords    []string
	listOfVerbsAndNouns    [][]string
	roomDescription        []string
	actionData             [][]int
	actionDescription      []string
	objectOriginalLocation []int
	objectLocation         []int
	foundWord              []int
	roomExit               [][]int
	statusFlag             []bool

	commandInHandle  *os.File
	commandOutHandle *os.File
	flagDebug        bool

	prngState int
}

func NewGameState() *GameState {
	gs := &GameState{
		alternateCounter: make([]int, alternateCounters),
		alternateRoom:    make([]int, alternateRoomRegisters),
		foundWord:        make([]int, 2),
		statusFlag:       make([]bool, statusFlags),
		prngState:        int(time.Now().Unix()) % 65536,
	}
	return gs
}

func (gs *GameState) getPRN() int {
	gs.prngState = (prngPrm * (gs.prngState + 1)) % prngPrime
	return gs.prngState % percentUnits
}

func (gs *GameState) getCommandInput() string {
	var inputData string

	if gs.commandInHandle == os.Stdin && gs.commandInHandle.Fd() == 0 {
		gs.commandInHandle = os.Stdin
	}

	scanner := bufio.NewScanner(gs.commandInHandle)
	if scanner.Scan() {
		inputData = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if gs.commandOutHandle != nil {
		fmt.Fprintln(gs.commandOutHandle, inputData)
	}

	return inputData
}

func commandlineHelp() {
	fmt.Print(`Usage: perlscott.pl [OPTION]... game_data_file
Scott Adams adventure game interpreter

-i, --input    Command input file
-o, --output   Command output file
-d, --debug    Show game debugging info
-h, --help     Display this help and exit
`)
	os.Exit(0)
}

func commandlineOptions() (*os.File, *os.File, bool) {
	var inHandle *os.File
	var outHandle *os.File
	var input string
	var output string
	var flagHelp bool
	var flagDebug bool

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-i", "--input":
			i++
			input = os.Args[i]
		case "-o", "--output":
			i++
			output = os.Args[i]
		case "-d", "--debug":
			flagDebug = true
		case "-h", "--help":
			flagHelp = true
		}
	}

	if flagHelp {
		commandlineHelp()
	}

	if input == "" {
		inHandle = os.Stdin
	} else {
		var err error
		inHandle, err = os.Open(input)
		if err != nil {
			panic(fmt.Sprintf("file \"%s\" not found", input))
		}
	}

	if output != "" {
		var err error
		outHandle, err = os.Create(output)
		if err != nil {
			panic(err)
		}
	}

	return inHandle, outHandle, flagDebug
}

func (gs *GameState) stripNounFromObjectDescription(objectNumber int) string {
	strippedText := gs.objectDescription[objectNumber]
	if strings.Contains(strippedText, "/") {
		strippedText = strings.Split(strippedText, "/")[0]
	}
	return strippedText
}

func (gs *GameState) checkAndChangeLightSourceStatus() bool {
	if gs.objectLocation[lightSourceID] == roomInventory {
		gs.alternateCounter[counterTimeLimit]--
		if gs.alternateCounter[counterTimeLimit] < 0 {
			fmt.Print("Light has run out\n")
			gs.objectLocation[lightSourceID] = 0
		} else if gs.alternateCounter[counterTimeLimit] < lightWarningThreshold {
			fmt.Printf("Light runs out in %v turns!\n", gs.alternateCounter[counterTimeLimit])
		}
	}
	return true
}

func (gs *GameState) showIntro() {
	cls()
	introMessage := `
                 *** Welcome ***

 Unless told differently you must find *TREASURES* 
and-return-them-to-their-proper--place!

I'm your puppet. Give me english commands that
consist of a noun and verb. Some examples...

To find out what you're carrying you might say: TAKE INVENTORY 
to go into a hole you might say: GO HOLE 
to save current game: SAVE GAME

You will at times need special items to do things: But I'm 
sure you'll be a good adventurer and figure these things out.

     Happy adventuring... Hit enter to start
`
	fmt.Println(introMessage)

	_ = gs.getCommandInput()
	cls()
}

func (gs *GameState) showRoomDescription() {
	if gs.statusFlag[flagNight] {
		if gs.objectLocation[lightSourceID] != roomInventory && gs.objectLocation[lightSourceID] != gs.currentRoom {
			fmt.Print("I can't see: Its too dark.\n")
			return
		}
	}

	if strings.HasPrefix(gs.roomDescription[gs.currentRoom], "*") {
		fmt.Printf("%s\n", gs.roomDescription[gs.currentRoom][1:])
	} else {
		fmt.Printf("I'm in a %s", gs.roomDescription[gs.currentRoom])
	}

	objectsFound := false
	for i, location := range gs.objectLocation {
		if location == gs.currentRoom {
			if !objectsFound {
				fmt.Print(". Visible items here: \n")
				objectsFound = true
			}
			fmt.Printf("%s. ", gs.stripNounFromObjectDescription(i))
		}
	}
	fmt.Println()

	exitFound := false
	for i, exit := range gs.roomExit[gs.currentRoom] {
		if exit != 0 {
			if !exitFound {
				fmt.Print("Obvious exits: ")
				exitFound = true
			}
			fmt.Printf("%s ", directionNounText[i])
		}
	}
	fmt.Println()
}

func (gs *GameState) handleGoVerb() {
	roomDark := gs.statusFlag[flagNight]
	if roomDark {
		roomDark = roomDark && (gs.objectLocation[lightSourceID] != gs.currentRoom)
		roomDark = roomDark && (gs.objectLocation[lightSourceID] != -1)

		if roomDark {
			fmt.Print("Dangerous to move in the dark!\n")
		}
	}
	if gs.foundWord[1] < 1 {
		fmt.Print("Give me a direction too.\n")
		return
	}
	directionDestination := gs.roomExit[gs.currentRoom][gs.foundWord[1]-1]
	if directionDestination < 1 {
		if roomDark {
			fmt.Print("I fell down and broke my neck.\n")
			directionDestination = gs.numberOfRooms
			gs.statusFlag[flagNight] = false
		} else {
			fmt.Print("I can't go in that direction\n")
			return
		}
	}
	gs.currentRoom = directionDestination
	gs.showRoomDescription()
}

func (gs *GameState) getCommandParameter(currentAction int) {
	conditionCode := 1
	for conditionCode != parConditionCode {
		conditionLine := gs.actionData[currentAction][gs.commandParameterIndex]
		gs.commandParameter = conditionLine / conditionDivisor
		conditionCode = conditionLine - gs.commandParameter*conditionDivisor
		gs.commandParameterIndex++
	}
}

func (gs *GameState) decodeCommandFromData(commandNumber int, actionID int) int {
	var commandCode int
	mergedCommandIndex := commandNumber/2 + actionCommandOffset

	if (commandNumber % 2) == 1 {
		commandCode = gs.actionData[actionID][mergedCommandIndex] % commandCodeDivisor
	} else {
		commandCode = gs.actionData[actionID][mergedCommandIndex] / commandCodeDivisor
	}
	return commandCode
}

func (gs *GameState) loadGameDataFile() {
	content, err := ioutil.ReadFile(gs.gameFile)
	if err != nil {
		panic(err)
	}
	fileContent := string(content)
	fileContent = unixNewlinesToSystemNewlines(fileContent)

	next := fileContent

	_, next = extractIntNumber(next) // game_bytes. ignoring.
	gs.numberOfObjects, next = extractIntNumber(next)
	gs.numberOfActions, next = extractIntNumber(next)
	gs.numberOfWords, next = extractIntNumber(next)
	gs.numberOfRooms, next = extractIntNumber(next)
	gs.maxObjectsCarried, next = extractIntNumber(next)
	if gs.maxObjectsCarried < 0 {
		gs.maxObjectsCarried = reallyBigNumber
	}
	gs.startingRoom, next = extractIntNumber(next)
	gs.numberOfTreasures, next = extractIntNumber(next)
	gs.wordLength, next = extractIntNumber(next)
	gs.timeLimit, next = extractIntNumber(next)
	gs.numberOfMessages, next = extractIntNumber(next)
	gs.treasureRoomID, next = extractIntNumber(next)

	gs.actionData = make([][]int, gs.numberOfActions+1)
	for i := range gs.actionData {
		gs.actionData[i] = make([]int, actionEntries)
	}
	for actionID := 0; actionID <= gs.numberOfActions; actionID++ {
		for actionIDEntry := 0; actionIDEntry < actionEntries; actionIDEntry++ {
			gs.actionData[actionID][actionIDEntry], next = extractIntNumber(next)
		}
	}

	gs.listOfVerbsAndNouns = make([][]string, (gs.numberOfWords+1)*2)
	for i := range gs.listOfVerbsAndNouns {
		gs.listOfVerbsAndNouns[i] = make([]string, 2)
	}
	for word := 0; word < (gs.numberOfWords+1)*2; word++ {
		var input string
		input, next = extractQuotedString(next)
		gs.listOfVerbsAndNouns[word/2][word%2] = input
	}

	gs.roomDescription = make([]string, gs.numberOfRooms+1)
	gs.roomExit = make([][]int, gs.numberOfRooms+1)
	for i := range gs.roomExit {
		gs.roomExit[i] = make([]int, 6)
	}
	for room := 0; room <= gs.numberOfRooms; room++ {
		gs.roomExit[room][0], next = extractIntNumber(next)
		gs.roomExit[room][1], next = extractIntNumber(next)
		gs.roomExit[room][2], next = extractIntNumber(next)
		gs.roomExit[room][3], next = extractIntNumber(next)
		gs.roomExit[room][4], next = extractIntNumber(next)
		gs.roomExit[room][5], next = extractIntNumber(next)
		gs.roomDescription[room], next = extractQuotedString(next)
	}

	gs.message = make([]string, gs.numberOfMessages+1)
	for currentMessage := 0; currentMessage <= gs.numberOfMessages; currentMessage++ {
		gs.message[currentMessage], next = extractQuotedString(next)
	}

	gs.objectDescription = make([]string, gs.numberOfObjects+1)
	gs.objectLocation = make([]int, gs.numberOfObjects+1)
	gs.objectOriginalLocation = make([]int, gs.numberOfObjects+1)
	for object := 0; object <= gs.numberOfObjects; object++ {
		gs.objectDescription[object], gs.objectLocation[object], next = extractStringAndNumber(next)
		gs.objectOriginalLocation[object] = gs.objectLocation[object]
	}

	gs.actionDescription = make([]string, gs.numberOfActions+1)
	for actionCounter := 0; actionCounter <= gs.numberOfActions; actionCounter++ {
		gs.actionDescription[actionCounter], next = extractQuotedString(next)
	}

	gs.adventureVersion, next = extractIntNumber(next)
	gs.adventureNumber, _ = extractIntNumber(next)

	for i := range gs.objectDescription {
		gs.objectDescription[i] = strings.ReplaceAll(gs.objectDescription[i], "`", "\"")
	}
	for i := range gs.message {
		gs.message[i] = strings.ReplaceAll(gs.message[i], "`", "\"")
	}
	for i := range gs.roomDescription {
		gs.roomDescription[i] = strings.ReplaceAll(gs.roomDescription[i], "`", "\"")
	}
}

func cls() {
	fmt.Print("\033[2J")
	fmt.Print("\033[0;0H")
}

func (gs *GameState) extractWords() {
	gs.extractedInputWords = nil
	gs.keyboardInput2 = strings.TrimSpace(gs.keyboardInput2)
	gs.extractedInputWords = strings.Fields(gs.keyboardInput2)
	if len(gs.extractedInputWords) == 0 {
		gs.extractedInputWords = append(gs.extractedInputWords, "")
	}

	gs.resolveGoShortcut()

	if len(gs.extractedInputWords) < 2 {
		gs.extractedInputWords = append(gs.extractedInputWords, "")
	}
	gs.globalNoun = gs.extractedInputWords[1]

	for verbOrNoun := 0; verbOrNoun <= 1; verbOrNoun++ {
		gs.foundWord[verbOrNoun] = 0
		var nonSynonym int
		for wordID, word := range gs.listOfVerbsAndNouns {
			if !strings.HasPrefix(word[verbOrNoun], "*") {
				nonSynonym = wordID
			}
			tempWord := word[verbOrNoun]
			tempWord = strings.TrimPrefix(tempWord, "*")
			tempWord = tempWord[:min(len(tempWord), gs.wordLength)]
			if tempWord == strings.ToUpper(gs.extractedInputWords[verbOrNoun][:min(len(gs.extractedInputWords[verbOrNoun]), gs.wordLength)]) {
				gs.foundWord[verbOrNoun] = nonSynonym
				break
			}
		}
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (gs *GameState) saveGame() {
	fmt.Print("Name of save file:\n")
	saveFilename := gs.getCommandInput()
	saveFilename = strings.TrimSpace(saveFilename)
	var saveData []string

	saveData = append(saveData, strconv.Itoa(gs.adventureVersion))
	saveData = append(saveData, strconv.Itoa(gs.adventureNumber))
	saveData = append(saveData, strconv.Itoa(gs.currentRoom))
	for _, room := range gs.alternateRoom {
		saveData = append(saveData, strconv.Itoa(room))
	}
	saveData = append(saveData, strconv.Itoa(gs.counterRegister))
	for _, counter := range gs.alternateCounter {
		saveData = append(saveData, strconv.Itoa(counter))
	}
	for _, location := range gs.objectLocation {
		saveData = append(saveData, strconv.Itoa(location))
	}
	for _, flag := range gs.statusFlag {
		var boolValue int
		if flag {
			boolValue = 1
		} else {
			boolValue = 0
		}
		saveData = append(saveData, strconv.Itoa(boolValue))
	}

	err := ioutil.WriteFile(saveFilename, []byte(strings.Join(saveData, "\n")), 0644)
	if err != nil {
		panic(err)
	}

}

func (gs *GameState) loadGame() bool {
	fmt.Print("Name of save file:\n")
	saveFilename := gs.getCommandInput()
	saveFilename = strings.TrimSpace(saveFilename)
	if _, err := os.Stat(saveFilename); os.IsNotExist(err) {
		fmt.Printf("Couldn't load \"%s\". Doesn't exist!\n", saveFilename)
		return false
	}
	var saveData []string

	content, err := ioutil.ReadFile(saveFilename)
	if err != nil {
		panic(err)
	}
	saveData = strings.Split(string(content), "\n")

	saveAdventureVersion, _ := strconv.Atoi(saveData[0])
	if saveAdventureVersion != gs.adventureVersion {
		fmt.Print("Invalid savegame version\n")
		return false
	}
	saveAdventureNumber, _ := strconv.Atoi(saveData[1])
	if saveAdventureNumber != gs.adventureNumber {
		fmt.Print("Invalid savegame adventure number\n")
		return false
	}

	gs.currentRoom, _ = strconv.Atoi(saveData[2])
	for i := 0; i < alternateRoomRegisters; i++ {
		gs.alternateRoom[i], _ = strconv.Atoi(saveData[3+i])
	}
	gs.counterRegister, _ = strconv.Atoi(saveData[3+alternateRoomRegisters])
	for i := 0; i < alternateCounters; i++ {
		gs.alternateCounter[i], _ = strconv.Atoi(saveData[4+alternateRoomRegisters+i])
	}
	for i := 0; i < len(gs.objectLocation); i++ {
		gs.objectLocation[i], _ = strconv.Atoi(saveData[4+alternateRoomRegisters+alternateCounters+i])
	}
	for i := 0; i < len(gs.statusFlag); i++ {
		tempFlagValue, _ := strconv.Atoi(saveData[4+alternateRoomRegisters+alternateCounters+len(gs.objectLocation)+i])
		if tempFlagValue == 0 {
			gs.statusFlag[i] = false
		} else {
			gs.statusFlag[i] = true
		}
	}

	return true

}

func (gs *GameState) runActions(inputVerb int, inputNoun int) {
	if inputVerb == verbGo && inputNoun <= directionNouns {
		gs.handleGoVerb()
		return
	}

	var foundWord int = 0

	gs.contFlag = false
	currentAction := 0
	wordActionDone := false
	for range gs.actionData {
		actionVerb := gs.getActionVerb(currentAction)
		actionNoun := gs.getActionNoun(currentAction)

		if gs.contFlag && actionVerb == 0 && actionNoun == 0 {
			gs.printDebug(fmt.Sprintf("Action %d. verb %v, noun %v (CONT %v), \"%s\"", currentAction, actionVerb, actionNoun, gs.contFlag, gs.actionDescription[currentAction]), 31)
			if gs.evaluateConditions(currentAction) {
				gs.executeCommands(currentAction)
			}
		} else {
			gs.contFlag = false
		}

		if inputVerb == 0 {
			if actionVerb == 0 && actionNoun > 0 {
				gs.printDebug(fmt.Sprintf("Action %d. verb %v, noun %v (CONT %v), \"%s\"", currentAction, actionVerb, actionNoun, gs.contFlag, gs.actionDescription[currentAction]), 31)
				gs.contFlag = false
				if gs.getPRN() < actionNoun {
					if gs.evaluateConditions(currentAction) {
						gs.executeCommands(currentAction)
					}
				}
			}
		}

		if inputVerb > 0 {
			if actionVerb == inputVerb {
				if !wordActionDone {
					gs.printDebug(fmt.Sprintf("Action %d. verb %v (%s), noun %v (%s) (CONT %v), \"%s\"", currentAction, actionVerb, gs.listOfVerbsAndNouns[actionVerb][0], actionNoun, gs.listOfVerbsAndNouns[actionNoun][1], gs.contFlag, gs.actionDescription[currentAction]), 31)
					gs.contFlag = false
					if actionNoun == 0 {
						foundWord = 1
						if gs.evaluateConditions(currentAction) {
							gs.executeCommands(currentAction)
							wordActionDone = true
							if !gs.contFlag {
								return
							}
						}
					} else if actionNoun == inputNoun {
						foundWord = 1
						if gs.evaluateConditions(currentAction) {
							gs.executeCommands(currentAction)
							wordActionDone = true
							if !gs.contFlag {
								return
							}
						}
					}
				}
			}
		}

		currentAction++
	}

	if inputVerb == 0 {
		return
	}

	if !wordActionDone {
		if gs.handleCarryAndDropVerb(inputVerb, inputNoun) {
			return
		}
	}

	if wordActionDone {
		return
	}

	if foundWord > 0 {
		fmt.Print("I can't do that yet\n")
	} else {
		fmt.Print("I don't understand your command\n")
	}

}

func (gs *GameState) printDebug(message string, color int) {
	if gs.flagDebug {
		fmt.Printf("\033[%dmDEBUG: %s\033[0m\n", color, message)
	}
}

func (gs *GameState) nounIsInObject() bool {
	truncatedNoun := gs.globalNoun[:gs.wordLength]
	for _, description := range gs.objectDescription {
		if strings.Contains(description, "/") {
			objectNoun := strings.ToLower(strings.Split(description, "/")[1])
			if objectNoun == truncatedNoun {
				return true
			}
		}
	}
	return false
}

func (gs *GameState) handleCarryAndDropVerb(inputVerb int, inputNoun int) bool {
	if inputVerb != verbCarry && inputVerb != verbDrop {
		return false
	}

	if inputNoun == 0 && !gs.nounIsInObject() {
		fmt.Print("What?\n")
		return true
	}

	if inputVerb == verbCarry {
		gs.carriedObjects = 0

		for _, location := range gs.objectLocation {
			if location == roomInventory {
				gs.carriedObjects++
			}
		}
		if gs.carriedObjects >= gs.maxObjectsCarried {
			if gs.maxObjectsCarried >= 0 {
				fmt.Print("I've too much too carry. try -take inventory-\n")
				return true
			}
		} else {
			if gs.getOrDropNoun(inputNoun, gs.currentRoom, roomInventory) {
				return true
			} else {
				fmt.Print("I don't see it here\n")
				return true
			}
		}
	} else {
		if gs.getOrDropNoun(inputNoun, roomInventory, gs.currentRoom) {
			return true
		} else {
			fmt.Print("I'm not carrying it\n")
			return true
		}
	}
	return false

}

func (gs *GameState) getOrDropNoun(inputNoun int, roomSource int, roomDestination int) bool {
	var objectsInRoom []int
	objectCounter := 0

	for _, location := range gs.objectLocation {
		if location == roomSource {
			objectsInRoom = append(objectsInRoom, objectCounter)
		}
		objectCounter++
	}

	for _, roomObject := range objectsInRoom {
		if strings.Contains(gs.objectDescription[roomObject], "/") {
			noun := strings.Split(gs.objectDescription[roomObject], "/")[1]
			if gs.listOfVerbsAndNouns[inputNoun][1] == noun || noun == strings.ToUpper(gs.globalNoun[:gs.wordLength]) {
				gs.objectLocation[roomObject] = roomDestination
				fmt.Print("OK\n")
				return true
			}
		}
	}
	return false

}

func (gs *GameState) getActionVerb(actionID int) int {
	return gs.actionData[actionID][0] / commandCodeDivisor
}

func (gs *GameState) getActionNoun(actionID int) int {
	return gs.actionData[actionID][0] % commandCodeDivisor
}

func (gs *GameState) executeCommands(actionID int) {
	gs.commandParameterIndex = 1
	for command := 0; command < commandsInAction; command++ {
		continueExecutingCommands := true
		gs.commandOrDisplayMessage = gs.decodeCommandFromData(command, actionID)

		if gs.commandOrDisplayMessage >= message2Start {
			gs.printDebug(fmt.Sprintf("Command print message %v", gs.commandOrDisplayMessage), 32)
			fmt.Printf("%s\n", gs.message[gs.commandOrDisplayMessage-message1End+1])
		} else if gs.commandOrDisplayMessage == 0 {
			// Do nothing
		} else if gs.commandOrDisplayMessage <= message1End {
			gs.printDebug(fmt.Sprintf("Command print message %v", gs.commandOrDisplayMessage), 32)
			fmt.Printf("%s\n", gs.message[gs.commandOrDisplayMessage])
		} else {
			commandCode := gs.commandOrDisplayMessage - message1End - 1

			gs.printDebug(fmt.Sprintf("Command code %v %s", commandCode, commandName[commandCode]), 32)
			commandFunction[commandCode](gs, actionID, &continueExecutingCommands)
		}

		if !continueExecutingCommands {
			break
		}
	}

}

func (gs *GameState) evaluateConditions(actionID int) bool {
	evaluationStatus := true
	for condition := 1; condition <= conditions; condition++ {
		conditionCode := gs.getConditionCode(actionID, condition)
		conditionParameter := gs.getConditionParameter(actionID, condition)
		gs.printDebug(fmt.Sprintf("Condition %v %s with parameter %v", conditionCode, conditionName[conditionCode], conditionParameter), 33)
		if !conditionFunction[conditionCode](gs, conditionParameter) {
			evaluationStatus = false
			break
		}
	}
	return evaluationStatus
}

func (gs *GameState) getConditionCode(actionID int, condition int) int {
	conditionRaw := gs.actionData[actionID][condition]
	conditionCode := conditionRaw % conditionDivisor
	return conditionCode
}

func (gs *GameState) getConditionParameter(actionID int, condition int) int {
	conditionRaw := gs.actionData[actionID][condition]
	conditionParameter := conditionRaw / conditionDivisor
	return conditionParameter
}

func (gs *GameState) resolveGoShortcut() {
	enteredInputVerb := strings.ToLower(gs.extractedInputWords[0])

	// Don't attempt to resolve go shortcuts if input is empty
	if len(enteredInputVerb) < 1 {
		return
	}

	// Don't make shortcut if input verb matches legitimate word action
	viablePhrases := gs.getViableWordActions()
	for viableVerb := range viablePhrases {
		possibleVerbText := strings.ToLower(gs.listOfVerbsAndNouns[viableVerb][0])
		shortenedVerb := enteredInputVerb[:min(len(enteredInputVerb), len(possibleVerbText))]
		if shortenedVerb == possibleVerbText {
			return
		}
	}

	for direction := 1; direction <= directionNouns; direction++ {
		directionNounText := strings.ToLower(gs.listOfVerbsAndNouns[direction][1])
		shortenedDirection := enteredInputVerb[:min(len(enteredInputVerb), len(directionNounText))]
		if shortenedDirection == directionNounText {
			gs.extractedInputWords[0] = strings.ToLower(gs.listOfVerbsAndNouns[verbGo][0])
			gs.extractedInputWords = append(gs.extractedInputWords, directionNounText)
			return
		}
	}

}

func (gs *GameState) getViableWordActions() map[int]map[int]string {
	viablePhrases := make(map[int]map[int]string)

	for currentAction := 0; currentAction < len(gs.actionData); currentAction++ {
		actionVerb := gs.getActionVerb(currentAction)
		actionNoun := gs.getActionNoun(currentAction)
		if actionVerb > 0 {
			if gs.evaluateConditions(currentAction) {
				if _, ok := viablePhrases[actionVerb]; !ok {
					viablePhrases[actionVerb] = make(map[int]string)
				}
				viablePhrases[actionVerb][actionNoun] = ""
			}
		}
	}
	return viablePhrases
}

// Condition functions
func Par(gs *GameState, parameter int) bool {
	return true
}
func HAS(gs *GameState, parameter int) bool {
	result := false
	if gs.objectLocation[parameter] == roomInventory {
		result = true
	}
	return result
}
func INW(gs *GameState, parameter int) bool {
	return gs.objectLocation[parameter] == gs.currentRoom
}
func AVL(gs *GameState, parameter int) bool {
	result := gs.objectLocation[parameter] == roomInventory
	result = result || (gs.objectLocation[parameter] == gs.currentRoom)
	return result
}
func IN(gs *GameState, parameter int) bool {
	return gs.currentRoom == parameter
}
func MinusINW(gs *GameState, parameter int) bool {
	return gs.objectLocation[parameter] != gs.currentRoom
}
func MinusHAVE(gs *GameState, parameter int) bool {
	return gs.objectLocation[parameter] != roomInventory
}
func MinusIN(gs *GameState, parameter int) bool {
	return gs.currentRoom != parameter
}
func BIT(gs *GameState, parameter int) bool {
	return gs.statusFlag[parameter]
}
func MinusBIT(gs *GameState, parameter int) bool {
	return !gs.statusFlag[parameter]
}
func ANY(gs *GameState, parameter int) bool {
	result := false
	for _, location := range gs.objectLocation {
		if location == roomInventory {
			result = true
		}
	}
	return result
}
func MinusANY(gs *GameState, parameter int) bool {
	result := false
	for _, location := range gs.objectLocation {
		if location == roomInventory {
			result = true
		}
	}
	return !result
}
func MinusAVL(gs *GameState, parameter int) bool {
	result := gs.objectLocation[parameter] == roomInventory
	result = result || (gs.objectLocation[parameter] == gs.currentRoom)
	return !result
}
func MinusRM0(gs *GameState, parameter int) bool {
	return gs.objectLocation[parameter] != roomStore
}
func RM0(gs *GameState, parameter int) bool {
	return gs.objectLocation[parameter] == roomStore
}
func CTLessEq(gs *GameState, parameter int) bool {
	return gs.counterRegister <= parameter
}
func CTGreater(gs *GameState, parameter int) bool {
	return gs.counterRegister > parameter
}
func ORIG(gs *GameState, parameter int) bool {
	return gs.objectOriginalLocation[parameter] == gs.objectLocation[parameter]
}
func MinusORIG(gs *GameState, parameter int) bool {
	return gs.objectOriginalLocation[parameter] != gs.objectLocation[parameter]
}
func CTEqual(gs *GameState, parameter int) bool {
	return gs.counterRegister == parameter
}

// Command functions
func GETx(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.carriedObjects = 0

	for _, location := range gs.objectLocation {
		if location == roomInventory {
			gs.carriedObjects++
		}
	}
	if gs.carriedObjects >= gs.maxObjectsCarried {
		fmt.Print("I've too much too carry. try -take inventory-\n")

		*continueExecutingCommands = false
	}

	gs.getCommandParameter(actionID)
	gs.objectLocation[gs.commandParameter] = roomInventory

}

func DROPx(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.objectLocation[gs.commandParameter] = gs.currentRoom
}

func GOTOy(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.currentRoom = gs.commandParameter
}

func xToRM0(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.objectLocation[gs.commandParameter] = 0
}

func NIGHT(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.statusFlag[flagNight] = true
}

func DAY(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.statusFlag[flagNight] = false
}

func SETz(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.statusFlag[gs.commandParameter] = true
}

func xToRM0b(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.objectLocation[gs.commandParameter] = 0
}

func CLRz(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.statusFlag[gs.commandParameter] = false
}

func DEAD(gs *GameState, actionID int, continueExecutingCommands *bool) {
	fmt.Print("I'm dead...\n")
	gs.currentRoom = gs.numberOfRooms
	gs.statusFlag[flagNight] = false
	gs.showRoomDescription()
}

func xToy(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	temporary1 := gs.commandParameter
	gs.getCommandParameter(actionID)
	gs.objectLocation[temporary1] = gs.commandParameter
}

func FINI(gs *GameState, actionID int, continueExecutingCommands *bool) {
	os.Exit(0)
}

func DspRM(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.showRoomDescription()
}

func SCORE(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.storedTreasures = 0
	for i, location := range gs.objectLocation {
		if location == gs.treasureRoomID {
			if strings.HasPrefix(gs.objectDescription[i], "") {
				gs.storedTreasures++
			}
		}
	}

	fmt.Printf("I've stored %v treasures. ON A SCALE OF 0 TO %v THAT RATES A %v\n", gs.storedTreasures, percentUnits, gs.storedTreasures/gs.numberOfTreasures*percentUnits)
	if gs.storedTreasures == gs.numberOfTreasures {
		fmt.Print("Well done.\n")
		os.Exit(0)
	}

}

func INV(gs *GameState, actionID int, continueExecutingCommands *bool) {
	fmt.Print("I'm carrying:\n")
	carryingNothingText := "Nothing"
	var objectText string
	for i, location := range gs.objectLocation {
		if location != roomInventory {
			continue
		} else {
			objectText = gs.stripNounFromObjectDescription(i)
		}
		fmt.Printf("%s. ", objectText)
		carryingNothingText = ""
	}
	fmt.Printf("%s\n\n", carryingNothingText)
}

func SET0(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.commandParameter = 0
	gs.statusFlag[gs.commandParameter] = true
}

func CLR0(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.commandParameter = 0
	gs.statusFlag[gs.commandParameter] = false
}

func FILL(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.alternateCounter[counterTimeLimit] = gs.timeLimit
	gs.objectLocation[lightSourceID] = roomInventory
	gs.statusFlag[flagLampEmpty] = false
}

func CLS(gs *GameState, actionID int, continueExecutingCommands *bool) {
	cls()
}

func SAVE(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.saveGame()
}

func EXxx(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	temporary1 := gs.commandParameter
	gs.getCommandParameter(actionID)
	temporary2 := gs.objectLocation[gs.commandParameter]
	gs.objectLocation[gs.commandParameter] = gs.objectLocation[temporary1]
	gs.objectLocation[temporary1] = temporary2
}

func CONT(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.contFlag = true
}

func AGETx(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.carriedObjects = 0
	gs.getCommandParameter(actionID)
	gs.objectLocation[gs.commandParameter] = roomInventory
}

func BYxx(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	firstObject := gs.commandParameter
	gs.getCommandParameter(actionID)
	secondObject := gs.commandParameter
	gs.objectLocation[firstObject] = gs.objectLocation[secondObject]
}

func DspRM2(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.showRoomDescription()
}

func CTMinus1(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.counterRegister--
}

func DspCT(gs *GameState, actionID int, continueExecutingCommands *bool) {
	fmt.Printf("%v", gs.counterRegister)
}

func CTFromn(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.counterRegister = gs.commandParameter
}

func EXRM0(gs *GameState, actionID int, continueExecutingCommands *bool) {
	temp := gs.currentRoom
	gs.currentRoom = gs.alternateRoom[0]
	gs.alternateRoom[0] = temp
}

func EXmCT(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	temp := gs.counterRegister
	gs.counterRegister = gs.alternateCounter[gs.commandParameter]
	gs.alternateCounter[gs.commandParameter] = temp
}

func CTPlusn(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.counterRegister += gs.commandParameter
}

func CTMinusn(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	gs.counterRegister -= gs.commandParameter

	if gs.counterRegister < minimumCounterValue {
		gs.counterRegister = minimumCounterValue
	}

}

func SAYw(gs *GameState, actionID int, continueExecutingCommands *bool) {
	fmt.Println(gs.globalNoun)
}

func SAYwCR(gs *GameState, actionID int, continueExecutingCommands *bool) {
	fmt.Printf("%s\n", gs.globalNoun)
}

func SAYCR(gs *GameState, actionID int, continueExecutingCommands *bool) {
	fmt.Print("\n")
}

func EXcCR(gs *GameState, actionID int, continueExecutingCommands *bool) {
	gs.getCommandParameter(actionID)
	temp := gs.currentRoom
	gs.currentRoom = gs.alternateRoom[gs.commandParameter]
	gs.alternateRoom[gs.commandParameter] = temp
}

func DELAY(gs *GameState, actionID int, continueExecutingCommands *bool) {
	time.Sleep(1 * time.Second)
}

// Utility functions

func unixNewlinesToSystemNewlines(input string) string {
	return strings.ReplaceAll(input, "\n", "\r\n")
}

func extractIntNumber(input string) (int, string) {
	re := regexp.MustCompile(`(?s)^\s*(-?\d+)(.*)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) == 0 {
		panic(fmt.Sprintf("Could not extract number from '%s'", input))
	}
	number, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(fmt.Sprintf("Error converting '%s' to integer: %v", matches[1], err))
	}
	return number, matches[2]
}

func extractQuotedString(input string) (string, string) {
	re := regexp.MustCompile(`(?s)^\s*"([^"]*)"(.*)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) == 0 {
		panic(fmt.Sprintf("Could not extract quoted string from '%s'", input))
	}
	return matches[1], matches[2]
}

func extractStringAndNumber(input string) (string, int, string) {
	str, rest := extractQuotedString(input)
	number, rest := extractIntNumber(rest)
	return str, number, rest
}

func intSliceToStringSlice(ints []int) []string {
	var strings []string
	for _, i := range ints {
		strings = append(strings, strconv.Itoa(i))
	}
	return strings
}

func boolSliceToStringSlice(bools []bool) []string {
	var strings []string
	for _, b := range bools {
		strings = append(strings, strconv.FormatBool(b))
	}
	return strings
}

var conditionFunction = []func(gs *GameState, parameter int) bool{
	Par, HAS, INW, AVL, IN, MinusINW, MinusHAVE, MinusIN, BIT, MinusBIT,
	ANY, MinusANY, MinusAVL, MinusRM0, RM0, CTLessEq, CTGreater, ORIG, MinusORIG, CTEqual,
}
var commandFunction = []func(gs *GameState, actionID int, continueExecutingCommands *bool){
	GETx, DROPx, GOTOy, xToRM0, NIGHT, DAY,
	SETz, xToRM0b, CLRz, DEAD, xToy, FINI,
	DspRM, SCORE, INV, SET0, CLR0, FILL,
	CLS, SAVE, EXxx, CONT, AGETx, BYxx,
	DspRM2, CTMinus1, DspCT, CTFromn, EXRM0, EXmCT,
	CTPlusn, CTMinusn, SAYw, SAYwCR, SAYCR, EXcCR,
	DELAY,
}
var commandName = []string{
	"GETx", "DROPx", "GOTOy", "x->RM0", "NIGHT", "DAY",
	"SETz", "x->RM0", "CLRz", "DEAD", "x->y", "FINI",
	"DspRM", "SCORE", "INV", "SET0", "CLR0", "FILL",
	"CLS", "SAVE", "EXx,x", "CONT", "AGETx", "BYx<-x",
	"DspRM", "CT-1", "DspCT", "CT<-n", "EXRM0", "EXm,CT",
	"CT+n", "CT-n", "SAYw", "SAYwCR", "SAYCR", "EXc,CR",
	"DELAY",
}

var conditionName = []string{
	"Par", "HAS", "IN/W", "AVL", "IN", "-IN/W", "-HAVE", "-IN",
	"BIT", "-BIT", "ANY", "-ANY", "-AVL", "-RM0", "RM0", "CT<=",
	"CT>", "ORIG", "-ORIG", "CT=",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	gs := NewGameState()
	gs.contFlag = false

	gs.commandInHandle, gs.commandOutHandle, gs.flagDebug = commandlineOptions()

	if len(os.Args) > 1 {
		gs.gameFile = os.Args[1]
		gs.loadGameDataFile()
	} else {
		commandlineHelp()
	}

	gs.currentRoom = gs.startingRoom
	for i := 0; i < alternateRoomRegisters; i++ {
		gs.alternateRoom[i] = 0
	}
	for i := 0; i < alternateCounters; i++ {
		gs.alternateCounter[i] = 0
	}
	gs.counterRegister = 0
	gs.statusFlag[flagNight] = false
	gs.alternateCounter[counterTimeLimit] = gs.timeLimit

	gs.showIntro()
	gs.showRoomDescription()

	gs.foundWord[0] = 0
	gs.runActions(gs.foundWord[0], 0)

	for {
		gs.printDebug(strings.Join(boolSliceToStringSlice(gs.statusFlag), " "), 37)
		gs.printDebug(strings.Join(intSliceToStringSlice(gs.alternateCounter), " "), 37)

		fmt.Print("Tell me what to do\n")

		gs.keyboardInput2 = gs.getCommandInput()
		gs.keyboardInput2 = strings.TrimSuffix(gs.keyboardInput2, "\n")
		fmt.Println()

		if strings.TrimSpace(strings.ToUpper(gs.keyboardInput2)) == "LOAD GAME" {
			if gs.loadGame() {
				gs.showRoomDescription()
			}
		} else {
			gs.extractWords()

			undefinedWordsFound := (gs.foundWord[0] < 1) || (len(gs.extractedInputWords[1]) > 0 && gs.foundWord[1] < 1)

			if gs.foundWord[0] == verbCarry || gs.foundWord[0] == verbDrop {
				undefinedWordsFound = false
			}

			if undefinedWordsFound {
				fmt.Print("You use word(s) i don't know\n")
			} else {
				gs.runActions(gs.foundWord[0], gs.foundWord[1])
				gs.checkAndChangeLightSourceStatus()
				gs.foundWord[0] = 0
				gs.runActions(gs.foundWord[0], gs.foundWord[1])
			}
		}
	}
}
