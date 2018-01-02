package main

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// These constants should be replaced by internal actions for GO and GET/DROP
// Some new internal conditions and commands are required for this
const verbGo = 2
const verbGet = 11
const verbDrop = 19
const nounNorth = 2
const nounSouth = 3
const nounEast = 4
const nounWest = 5
const nounUp = 6
const nounDown = 7
const flagDark = 16
const flagLampEmpty = 17
const itemLight = 10

func loadData(filename string) (
	map[string]int, // advVariable
	[][]int, // action
	[]string, // verb
	[]string, // noun
	[]map[int]int, // roomDirection
	[]string, // roomDescription
	[]string, // message
	[]string, // itemDescription
	[]string, // itemNoun
	[]int, // itemStartLocation
	[]string, // actionComment
	[]int, // treasureItem

) {
	fieldIndex := 0
	var advVariable map[string]int
	advVariable = make(map[string]int)

	// Create the following lists of items from the data file:
	/* action, verb, noun, roomDirection, roomDescription, message,
	   itemDescription, itemNoun, itemStartLocation, actionComment,
	   treasureItem */

	var advData []string
	advData = getDataArray(filename)

	// Get header variables
	advVariable["sizeOfText"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfItems"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfActions"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfActions"]++
	advVariable["numberOfWords"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfWords"]++
	advVariable["numberOfRooms"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfRooms"]++
	advVariable["maxItemsCarried"] = getNumber(advData, &fieldIndex)
	advVariable["startingRoom"] = getNumber(advData, &fieldIndex)
	advVariable["totalTreasures"] = getNumber(advData, &fieldIndex)
	advVariable["wordLength"] = getNumber(advData, &fieldIndex)
	advVariable["timeLimit"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfMessages"] = getNumber(advData, &fieldIndex)
	advVariable["numberOfMessages"]++
	advVariable["treasureRoom"] = getNumber(advData, &fieldIndex)

	// Get actions
	action := make([][]int, advVariable["numberOfActions"])
	for i := 0; i < advVariable["numberOfActions"]; i++ {
		action[i] = getAction(advData, &fieldIndex)
	}

	// Get words
	var verb, noun []string
	for i := 0; i < advVariable["numberOfWords"]; i++ {
		verb = append(verb, getText(advData, &fieldIndex))
		noun = append(noun, getText(advData, &fieldIndex))
	}

	// Get rooms
	var roomDescription []string
	var roomDirection []map[int]int
	for i := 0; i < advVariable["numberOfRooms"]; i++ {
		description, exit := getRoom(advData, &fieldIndex)
		roomDescription = append(roomDescription, description)
		roomDirection = append(roomDirection, exit)
	}

	// Get messages
	var message []string
	for i := 0; i < advVariable["numberOfMessages"]; i++ {
		tempMessage := getText(advData, &fieldIndex)
		message = append(message, tempMessage)
	}

	// Get items
	var itemDescription []string
	var itemNoun []string
	var itemStartLocation []int
	var treasureItem []int
	for i := 0; i < advVariable["numberOfItems"]+1; i++ {
		description, noun, startLocation := getItem(advData, &fieldIndex)
		itemDescription = append(itemDescription, description)
		itemNoun = append(itemNoun, noun)
		itemStartLocation = append(itemStartLocation, startLocation)
		if isTreasure(description) {
			treasureItem = append(treasureItem, i)
		}
	}

	// Get action comments
	var actionComment []string
	for i := 0; i < advVariable["numberOfActions"]; i++ {
		comment := getText(advData, &fieldIndex)
		actionComment = append(actionComment, comment)
	}

	// Get footer variables
	advVariable["engineVersion"] = getNumber(advData, &fieldIndex)
	advVariable["adventureNumber"] = getNumber(advData, &fieldIndex)
	advVariable["gameChecksum"] = getNumber(advData, &fieldIndex)
	return advVariable, action, verb, noun, roomDirection, roomDescription, message,
		itemDescription, itemNoun, itemStartLocation, actionComment,
		treasureItem
}

func getRoom(advField []string, fieldIndex *int) (string, map[int]int) {
	roomDirection := make(map[int]int)
	var exit int
	for i := 0; i < 6; i++ {
		exit = getNumber(advField, fieldIndex)
		if exit != 0 {
			roomDirection[i+1] = exit
		}
	}
	description := getText(advField, fieldIndex)
	return description, roomDirection
}

func getAction(advField []string, fieldIndex *int) []int {
	var actionPart []int
	var actionEntry []int
	for i := 0; i < 8; i++ {
		actionPart = append(actionPart, getNumber(advField, fieldIndex))
	}

	// Extract verb and noun for an action
	var verb, noun int
	verb = actionPart[0] / 150
	noun = actionPart[0] % 150
	actionEntry = append(actionEntry, verb)
	actionEntry = append(actionEntry, noun)

	// Extract condition code and condition data for an action
	var code, data int
	for i := 1; i < 6; i++ {
		code = actionPart[i] % 20
		data = actionPart[i] / 20
		actionEntry = append(actionEntry, code)
		actionEntry = append(actionEntry, data)
	}

	// Extract commands
	{
		var first_command, second_command int
		for i := 6; i < 8; i++ {
			first_command = actionPart[i] % 150
			second_command = actionPart[i] / 150
			actionEntry = append(actionEntry, first_command)
			actionEntry = append(actionEntry, second_command)
		}
	}
	return actionEntry
}

func getText(advField []string, fieldIndex *int) string {
	textLine := ""
	quotes := regexp.MustCompile(`"[^"]*`)
	padding := regexp.MustCompile(`^\s*"(\S*(\s+\S+)*\s*)"\s*$`)
	foundQuotes := 0
	for foundQuotes < 2 {
		textLine += "\n" + advField[*fieldIndex]
		foundQuotes = len(quotes.FindAllString(textLine, -1))
		*fieldIndex++
	}
	textLine = padding.ReplaceAllString(textLine, "$1")
	return textLine
}

func getItem(advField []string, fieldIndex *int) (string, string, int) {
	textLine := ""
	quotes := regexp.MustCompile(`"[^"]*`)
	number := regexp.MustCompile(`-?\d+\s*$`)
	padding := regexp.MustCompile(`^\s*"(\S*(\s+\S+)*\s*)"\s*$`)
	fields := regexp.MustCompile(`^\s*"(\S*(?:\s+\S+?)*?\s*?)(?:\/(.*)\/)?"\s*(-?\d+)\s*$`)
	foundQuotes := 0
	foundNumber := 0
	textLine = advField[*fieldIndex]
	foundQuotes = len(quotes.FindAllString(textLine, -1))
	foundNumber = len(number.FindAllString(textLine, -1))
	*fieldIndex++
	for !((foundQuotes == 2) && (foundNumber == 1)) {
		textLine += "\n" + advField[*fieldIndex]
		foundQuotes = len(quotes.FindAllString(textLine, -1))
		foundNumber = len(number.FindAllString(textLine, -1))
		*fieldIndex++
	}
	textLine = padding.ReplaceAllString(textLine, "$1")
	textField := fields.FindStringSubmatch(textLine)[1:]
	description := textField[0]
	noun := textField[1]
	textRoomNumber := textField[2]
	roomNumber, _ := strconv.Atoi(textRoomNumber)
	return description, noun, roomNumber
}

func isTreasure(description string) bool {
	treasure := regexp.MustCompile(`^\*`)
	return treasure.MatchString(description)
}

func getNumber(advField []string, fieldIndex *int) int {
	var decodedNumber int
	r, _ := regexp.Compile(`([0-9]+)`)
	var cleanedText []string
	cleanedText = r.FindStringSubmatch(advField[*fieldIndex])
	*fieldIndex++
	decodedNumber, _ = strconv.Atoi(cleanedText[0])
	return decodedNumber
}

func getWord(advField []string, fieldIndex *int) [2]string {
	var word [2]string
	word[0] = getText(advField, fieldIndex)
	word[1] = getText(advField, fieldIndex)
	return word
}

func getDataArray(filename string) []string {
	rawByteData, err := ioutil.ReadFile(filename)
	check(err)
	var advData []string
	advData = strings.Split(string(rawByteData), "\n")
	return advData
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
