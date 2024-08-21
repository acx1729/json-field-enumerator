package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Item struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Label              string   `json:"label"`
	Tier               string   `json:"tier"`
	Tags               []string `json:"tags"`
	AutoOnboardSupport bool     `json:"auto_onboard_support"`
	ShortDescription   string   `json:"short_description"`
	Description        string   `json:"description"`
	Logo               string   `json:"logo"`
}

func main() {
	// Define command-line arguments
	filename := flag.String("filename", "sample.json", "Path to the JSON file")
	enumField := flag.String("enum-field", "tags", "Field to enumerate (e.g., tier, tags)")
	minCount := flag.Int("min-count", 1, "Minimum count to display")
	sortBy := flag.String("sort-by", "name", "Sort results by 'name' or 'count'")

	flag.Parse()

	// Read the JSON file
	file, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON data
	var items []Item
	err = json.Unmarshal(file, &items)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Create a map to hold the counts
	counts := make(map[string]int)

	// Enumerate based on the specified field
	for _, item := range items {
		switch *enumField {
		case "tier":
			counts[item.Tier]++
		case "tags":
			for _, tag := range item.Tags {
				counts[tag]++
			}
		default:
			fmt.Printf("Unknown enum field: %s\n", *enumField)
			os.Exit(1)
		}
	}

	// Filter results based on the min-count argument
	filteredCounts := make(map[string]int)
	for key, count := range counts {
		if count >= *minCount {
			filteredCounts[key] = count
		}
	}

	// Sort the results
	var sortedKeys []string
	for key := range filteredCounts {
		sortedKeys = append(sortedKeys, key)
	}

	switch *sortBy {
	case "name":
		sort.Strings(sortedKeys)
	case "count":
		sort.Slice(sortedKeys, func(i, j int) bool {
			return filteredCounts[sortedKeys[i]] > filteredCounts[sortedKeys[j]]
		})
	default:
		fmt.Printf("Unknown sort-by option: %s\n", *sortBy)
		os.Exit(1)
	}

	// Print the results
	for _, key := range sortedKeys {
		fmt.Printf("%s: %d\n", key, filteredCounts[key])
	}
}
