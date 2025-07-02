package main

import (
	"fmt"
	"strings"
)

var dictionary = []string{
	"apple",
	"banana",
	"cherry",
	"date",
	"elderberry",
	"fig",
	"grape",
	"honeydew",
	"kiwi",
	"lemon",
	"mango",
	"nectarine",
	"orange",
	"papaya",
	"quince",
	"raspberry",
	"strawberry",
	"tangerine",
	"watermelon",
	"yam",
	"zucchini",
}

var ASCIISIZE int = 256

var keyboard = [][]byte{
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm'},
}

type Cell struct {
	Cost              int
	IndexesAddedWords []int
	LastWord          string
}

func main() {
	char2CharDist := make([][]int, ASCIISIZE)
	for row := range keyboard {
		for col, char := range keyboard[row] {
			char2CharDist[char] = bfs(row, col, keyboard)
		}
	}

	k := 4
	minLen := 20
	maxLen := 24

	dp := dpAlgorithm(dictionary, char2CharDist, k, minLen, maxLen, 'g')

	var bestCell *Cell
	for curLen := minLen; curLen <= maxLen; curLen++ {
		if dp[k][curLen].IndexesAddedWords != nil {
			if bestCell == nil || dp[k][curLen].Cost < bestCell.Cost {
				bestCell = &dp[k][curLen]
			}
		}
	}

	if bestCell == nil {
		fmt.Println("нет решений")
		return
	}

	var builder strings.Builder
	for _, indexWord := range bestCell.IndexesAddedWords {
		builder.WriteString(dictionary[indexWord])
	}
	fmt.Printf("надежынй пароль: %s\nколичество передвижений пальцем: %v", builder.String(), bestCell.Cost)
}

// k - number words in password
// minLen - min len of password
// maxLen - max len of password
func dpAlgorithm(dictionary []string, char2CharDist [][]int, k, minLen, maxLen int, initChar byte) [][]Cell {
	dp := make([][]Cell, k+1)
	for i := range dp {
		dp[i] = make([]Cell, maxLen+1)
	}
	dp[0][0] = Cell{
		Cost:              0,
		IndexesAddedWords: []int{},
		LastWord:          string(initChar),
	}

	for indexWord, word := range dictionary {
		wordCost := calculateWordCost(word, char2CharDist)
		wordLen := len(word)

		for curK := k - 1; curK >= 0; curK-- {
			for curLen := 0; curLen <= maxLen-wordLen; curLen++ {
				if dp[curK][curLen].IndexesAddedWords != nil {
					newK := curK + 1
					newLen := curLen + wordLen
					newCost := dp[curK][curLen].Cost + calculateCostPath2Word(dp[curK][curLen].LastWord, word, char2CharDist) + wordCost

					if dp[newK][newLen].IndexesAddedWords == nil {
						dp[newK][newLen] = Cell{
							Cost:              newCost,
							IndexesAddedWords: append([]int{}, dp[curK][curLen].IndexesAddedWords...),
							LastWord:          word,
						}
						dp[newK][newLen].IndexesAddedWords = append(dp[newK][newLen].IndexesAddedWords, indexWord)

					} else if dp[newK][newLen].Cost > newCost {
						dp[newK][newLen].IndexesAddedWords[len(dp[newK][newLen].IndexesAddedWords)-1] = indexWord
						dp[newK][newLen].Cost = newCost
					}
				}
			}
		}
	}

	return dp
}

func bfs(row, col int, keyboard [][]byte) []int {
	dist := make([]int, ASCIISIZE)
	for i := range dist {
		dist[i] = -1
	}
	dist[keyboard[row][col]] = 0

	type Point struct {
		Row int
		Col int
	}
	var q = []Point{
		{Row: row, Col: col},
	}

	dCol := []int{1, 0, -1, 0}
	dRow := []int{0, 1, 0, -1}

	for len(q) > 0 {
		curChar := q[0]
		q = q[1:]

		for i := 0; i < 4; i++ {
			newRow := curChar.Row + dRow[i]
			newCol := curChar.Col + dCol[i]

			if newRow >= 0 && newRow < len(keyboard) && newCol >= 0 && newCol < len(keyboard[newRow]) && dist[keyboard[newRow][newCol]] < 0 {
				dist[keyboard[newRow][newCol]] = dist[keyboard[curChar.Row][curChar.Col]] + 1
				q = append(q, Point{Row: newRow, Col: newCol})
			}
		}
	}

	return dist
}

func calculateWordCost(word string, char2CharDist [][]int) int {
	var cost int
	for i := 1; i < len(word); i++ {
		cost += char2CharDist[word[i-1]][word[i]]
	}
	return cost
}

func calculateCostPath2Word(from, to string, char2CharDist [][]int) int {
	if len(to) == 0 || len(from) == 0 {
		return 0
	}
	return char2CharDist[from[len(from)-1]][to[0]]
}
