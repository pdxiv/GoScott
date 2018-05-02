package main

import (
	"bufio"
	"fmt"
	"os"
)

func getConsoleInput( loadedGame *gameStaticData) {
	wordLength := loadedGame.advVariable["wordLength"]
	fmt.Println("wordLength", wordLength)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func extractWords() {

}