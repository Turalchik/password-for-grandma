package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var ASCIISIZE int = 256

var keyboard = [][]byte{
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm'},
}

type Cell struct {
	Cost              int
	IndexesAddedWords map[int]struct{}
}

func loadDictionary(path string, minWordLen, maxWordLen int) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("не удалось открыть словарь %s: %v", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var dict []string
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if len(w) >= minWordLen && len(w) <= maxWordLen && isASCIIAlpha(w) {
			dict = append(dict, strings.ToLower(w))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("ошибка чтения словаря: %v", err)
	}
	return dict
}

func isASCIIAlpha(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || (c > 'Z' && c < 'a') || c > 'z' {
			return false
		}
	}
	return true
}

func main() {

	dictFile := "./words"
	dictionary := loadDictionary(dictFile, 1, 24)
	fmt.Printf("Загружено %d слов\n", len(dictionary))

	start := time.Now()

	char2CharDist := make([][]int, ASCIISIZE)
	for row := range keyboard {
		for col, char := range keyboard[row] {
			char2CharDist[char] = bfs(row, col, keyboard)
		}
	}

	k := 4
	minLen, maxLen := 20, 24
	var initChar byte = 'g'
	dp := dpAlgorithm(dictionary, char2CharDist, k, minLen, maxLen, initChar)

	var best *Cell
	for L := minLen; L <= maxLen; L++ {
		for c := 0; c < ASCIISIZE; c++ {
			cell := dp[k][L][c]
			if cell.IndexesAddedWords != nil && (best == nil || cell.Cost < best.Cost) {
				best = &cell
			}
		}
	}

	if best == nil {
		fmt.Println("нет решений")
		return
	}
	var sb strings.Builder
	for idx, _ := range best.IndexesAddedWords {
		sb.WriteString(dictionary[idx])
	}

	elapsed := time.Since(start)

	fmt.Printf("надежынй пароль: %s\nперемещений: %d\nвремя работы алгоритма: %v\n",
		sb.String(), best.Cost, elapsed)
}

// k - number words in password
// minLen - min len of password
// maxLen - max len of password
func dpAlgorithm(dictionary []string, char2CharDist [][]int, k, minLen, maxLen int, initChar byte) [][][]Cell {
	dp := make([][][]Cell, k+1)
	for layer := range dp {
		dp[layer] = make([][]Cell, maxLen+1)
		for L := range dp[layer] {
			dp[layer][L] = make([]Cell, ASCIISIZE)
		}
	}

	dp[0][0][initChar].Cost = 0
	dp[0][0][initChar].IndexesAddedWords = make(map[int]struct{})

	for layer := 0; layer < k; layer++ {
		for L := 0; L <= maxLen; L++ {
			for last := 0; last < ASCIISIZE; last++ {
				cell := dp[layer][L][last]
				if cell.IndexesAddedWords == nil {
					continue
				}

				for idx, word := range dictionary {
					newL := L + len(word)
					if newL > maxLen {
						continue
					}

					if _, ok := cell.IndexesAddedWords[idx]; ok {
						continue
					}

					cost := cell.Cost + char2CharDist[last][word[0]] + calculateWordCost(word, char2CharDist)
					nextCell := &dp[layer+1][newL][word[len(word)-1]]

					if nextCell.IndexesAddedWords == nil || cost < nextCell.Cost {
						newPath := copyMap(cell.IndexesAddedWords)
						newPath[idx] = struct{}{}
						nextCell.Cost = cost
						nextCell.IndexesAddedWords = newPath
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

func copyMap(orig map[int]struct{}) map[int]struct{} {
	newMap := make(map[int]struct{}, len(orig))
	for k, _ := range orig {
		newMap[k] = struct{}{}
	}
	return newMap
}
