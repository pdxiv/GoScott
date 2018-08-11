package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Placeholder console UI code
func getConsoleInput(loadedGame *gameStaticData) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	textInput, _ := reader.ReadString('\n')
	splitWord, extractStatus := extractWords(textInput)
	if extractStatus {
		truncatedWord := truncateWords(loadedGame, splitWord)
		fmt.Printf("\"%s\", \"%s\"\n", truncatedWord[0], truncatedWord[1])
		identifyWords(loadedGame, truncatedWord)
	}
}

func identifyWords(loadedGame *gameStaticData, sentence []string) {
	// Identify possible verbs
	wordLength := loadedGame.advVariable["wordLength"]
	fmt.Println(findWordInList(sentence[0], loadedGame.verb, wordLength, 1))
	fmt.Println(findWordInList(sentence[1], loadedGame.noun, wordLength, 1))
	fmt.Println(findWordInList(sentence[1], loadedGame.itemNoun, wordLength, 0))
}

func findWordInList(wordToLookFor string, wordList []string, wordLength int, listOffset int) []int {
	baseWord := listOffset
	var result []int
	for i, currentWord := range wordList[listOffset:] {
		wordToEvaluate := currentWord
		if len(currentWord) == 0 {
			continue
		}
		if currentWord[0:1] != "*" {
			baseWord = i + listOffset
		} else {
			wordToEvaluate = currentWord[1:]
		}
		if len(wordToEvaluate) > wordLength {
			wordToEvaluate = wordToEvaluate[0:wordLength]
		}
		if wordToEvaluate == wordToLookFor {
			result = append(result, baseWord)
		}
	}
	return result
}

func truncateWords(loadedGame *gameStaticData, sentence []string) []string {
	wordLength := loadedGame.advVariable["wordLength"]
	var truncatedWord []string
	for _, currentWord := range sentence {
		truncated := currentWord

		if len(currentWord) > wordLength {
			truncated = currentWord[0:wordLength]
		}
		truncatedWord = append(truncatedWord, truncated)
	}
	return truncatedWord
}

func extractWords(textInput string) ([]string, bool) {
	r := regexp.MustCompile(`(?P<Verb>\w+)\s*(?P<Noun>\w*)`)
	match := r.MatchString(textInput)
	var verbText, nounText string
	if match != true {
		return []string{"", ""}, false
	}
	verbText = r.FindStringSubmatch(strings.ToUpper(textInput))[1]
	nounText = r.FindStringSubmatch(strings.ToUpper(textInput))[2]
	return []string{verbText, nounText}, true
}
