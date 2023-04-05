package main

import (
	"bufio"
	"embed"
	"sort"
	"strings"
	"syscall/js"
)

//go:embed dictionary.txt
var f embed.FS

var dict map[string]bool
var hits map[string]bool

var resultElement js.Value

func consoleLog(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}

func search(this js.Value, args []js.Value) interface{} {
	letters := args[0].String()

	// reset hits
	hits = make(map[string]bool)

	perm([]byte(strings.ToUpper(letters)), 0)

	if len(hits) == 0 {
		resultElement.Set("innerText", "no results found")
	} else {
		hitsList := make([]string, 0, len(hits))
		for hit := range hits {
			hitsList = append(hitsList, hit)
		}
		// sort by longest hits first, then lexicographically
		sort.Slice(hitsList, func(i, j int) bool {
			if len(hitsList[i]) > len(hitsList[j]) {
				return true
			}
			if len(hitsList[i]) < len(hitsList[j]) {
				return false
			}
			return hitsList[i] < hitsList[j]
		})
		resultElement.Set("innerText", strings.Join(hitsList, "\n"))
	}
	return true
}

func perm(str []byte, i int) {
	word := string(str[:i])
	if _, ok := dict[word]; ok {
		hits[word] = true
	}
	if i != len(str) {
		for j := i; j < len(str); j++ {
			str[i], str[j] = str[j], str[i]
			perm(str, i+1)
			str[i], str[j] = str[j], str[i]
		}
	}
}

func loadDictionary() {
	dict = make(map[string]bool)

	file, err := f.Open("dictionary.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a buffered scanner
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		dict[line] = true
	}

	// Check for any errors that occurred while reading the file
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	consoleLog("dictionary entry count:", len(dict))
}

func main() {
	loadDictionary()

	resultElement = js.Global().Get("document").Call("getElementById", "results")

	js.Global().Set("search", js.FuncOf(search))

	select {} // Block the main function from exiting
}
