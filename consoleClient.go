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
		fmt.Printf("DEBUG: truncated words: \"%s\", \"%s\"\n", truncatedWord[0], truncatedWord[1])
		identifiedWords := identifyWords(loadedGame, truncatedWord)
		fmt.Println("DEBUG: found verb(s):", identifiedWords.verb)
		fmt.Println("DEBUG: found noun(s):", identifiedWords.noun)
		fmt.Println("DEBUG: found noun for object(s):", identifiedWords.object)
		fmt.Println("DEBUG: found directions for noun(s):", identifiedWords.direction)
	}
}

func identifyWords(loadedGame *gameStaticData, sentence []string) identifiedWords {
	// Identify possible verbs
	wordLength := loadedGame.advVariable["wordLength"]
	var result identifiedWords
	result.verb = findWordInList(sentence[0], loadedGame.verb, wordLength, 1)
	result.noun = findWordInList(sentence[1], loadedGame.noun, wordLength, 1)
	result.object = findWordInList(sentence[1], loadedGame.itemNoun, wordLength, 0)
	// No noun found for object? Check if there's a synonym noun!
	if len(result.object) <= 0 {
		for _, wordIndex := range result.noun {
			rootNoun := loadedGame.itemNoun[wordIndex]
			foundItemIndex := findWordInList(rootNoun, loadedGame.itemNoun, wordLength, 0)
			if len(foundItemIndex) > 0 {
				result.object = append(result.object, foundItemIndex[0])
			}
		}
	}
	// Resolve the exit indices if word corresponds to a direction noun
	for directionIndex, directionWordIndex := range []int{nounNorth, nounSouth, nounEast, nounWest, nounUp, nounDown} {
		fmt.Println(directionIndex, directionWordIndex)
		for _, resultNoun := range result.noun {
			if resultNoun == directionWordIndex {
				result.direction = append(result.direction, directionIndex)
			}
		}
	}
	return result
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
