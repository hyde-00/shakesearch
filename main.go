package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

// OtherContent denotes the result of query not part of any title.
var OtherContent = "OTHER CONTENT"

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	searcher.GetWorkTitles()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Searcher defines the struct used for searching.
type Searcher struct {
	CompleteWorks     string
	SuffixArray       *suffixarray.Index
	WorkTitlesNames   []string
	WorkTitlesIndices []int
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		var sanitizedQuery string
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)

		sanitizedQuery = string(searcher.Sanitize(strings.Join(query, " ")))
		results, countResults := searcher.Search(sanitizedQuery)
		finalResults := map[string]interface{}{
			"data":  results,
			"count": countResults,
		}
		err := enc.Encode(finalResults)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

// Load reads the completeworks.txt and stores it in-memory
func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(s.Sanitize(string(dat)))
	return nil
}

// Sanitize removes case in string to enable case insensitive search.
func (s *Searcher) Sanitize(data string) []byte {
	sanitizedData := strings.ToLower(data)
	return []byte(sanitizedData)
}

// GetWorkTitles fetches titles from the content to describe in
// which title the search query was found.
func (s *Searcher) GetWorkTitles() {
	allWorkTitles := []string{}
	foundFlag := false
	idxs := s.SuffixArray.Lookup([]byte("contents"), 1)

	foundString := s.CompleteWorks[idxs[0] : idxs[0]+2500]
	tmp := strings.Split(foundString, "\n")
	for _, val := range tmp {
		foundFlag = false
		trimmedValue := strings.TrimSpace(val)
		if len(trimmedValue) > 1 && IsUpper(val) {
			for _, value := range allWorkTitles {
				if trimmedValue == value {
					foundFlag = true
					break
				}
			}
			if foundFlag == false {
				allWorkTitles = append(allWorkTitles, strings.TrimSpace(val))
			}
		}
	}

	// Create a slice of titles and their respective indexes.
	// This is then compared to the results, which can help determine to which title
	// the result was found in.
	for _, workTitle := range allWorkTitles {
		idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(workTitle)), 1)
		s.WorkTitlesIndices = append(s.WorkTitlesIndices, idxs[0])
		s.WorkTitlesNames = append(s.WorkTitlesNames, workTitle)
	}
}

// Search through the text using SuffixArray
func (s *Searcher) Search(query string) (map[string][]string, int) {
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	var foundString string
	var mappingTitle string
	var results string
	mappedResults := make(map[string][]string)

	for _, idx := range idxs {
		// Get the string till a complete word, so that a word is not broken while splicing.
		if len(s.CompleteWorks) > idx+30 {
			foundString = s.CompleteWorks[idx-30 : idx+30]
			tmpSlice := strings.Split(foundString, " ")
			foundSlice := tmpSlice[1 : len(tmpSlice)-1]
			results = strings.Join(foundSlice, " ")
		} else if idx == 0 {
			foundString = s.CompleteWorks[0 : idx+30]
			tmpSlice := strings.Split(foundString, " ")
			foundSlice := tmpSlice[:len(tmpSlice)-1]
			results = strings.Join(foundSlice, " ")
		} else if len(s.CompleteWorks) <= idx+30 {
			foundString = s.CompleteWorks[idx-30 : idx+len(query)]
			results = s.CompleteWorks[idx-30 : idx+len(query)]
		}

		mappingTitle = s.MapResultToTitle(idx)
		mappedResults[mappingTitle] = append(mappedResults[mappingTitle], results)
	}
	return mappedResults, len(idxs)
}

// MapResultToTitle creates a mapping of the title to the query result.
func (s *Searcher) MapResultToTitle(index int) string {
	for idx, value := range s.WorkTitlesIndices {
		if value >= index && idx > 0 {
			return s.WorkTitlesNames[idx-1]
		}
	}
	return OtherContent
}

// IsUpper helper function to check if string is completely in uppercase or not.
func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
