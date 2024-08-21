package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// Connector struct to hold each entry from the JSON
type Connector struct {
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

// TagInfo struct to hold tag counts and associated labels
type TagInfo struct {
	Count  int
	Labels []string
}

func main() {
	// Optional command-line arguments
	minCount := flag.Int("min-count", 0, "Minimum count to filter tags")
	sortBy := flag.String("sort-by", "name", "Sort by 'name' or 'count'")

	flag.Parse()

	// Read the JSON file
	file, err := ioutil.ReadFile("sample.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	// Unmarshal the JSON into a slice of Connectors
	var connectors []Connector
	err = json.Unmarshal(file, &connectors)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	// Map to hold tag counts and associated labels
	tagMap := make(map[string]*TagInfo)

	// Process the connectors and populate the tagMap
	for _, connector := range connectors {
		tags := connector.Tags
		if len(tags) == 0 {
			tags = []string{"ALL OTHERS"}
		}
		for _, tag := range tags {
			if _, exists := tagMap[tag]; exists {
				tagMap[tag].Count++
				tagMap[tag].Labels = append(tagMap[tag].Labels, connector.Label)
			} else {
				tagMap[tag] = &TagInfo{Count: 1, Labels: []string{connector.Label}}
			}
		}
	}

	// Slice to hold sorted tags
	type TagEntry struct {
		Tag  string
		Info *TagInfo
	}

	var sortedTags []TagEntry
	for tag, info := range tagMap {
		if info.Count >= *minCount {
			sortedTags = append(sortedTags, TagEntry{Tag: tag, Info: info})
		}
	}

	// Sorting based on the "sort-by" argument
	if *sortBy == "count" {
		sort.Slice(sortedTags, func(i, j int) bool {
			return sortedTags[i].Info.Count > sortedTags[j].Info.Count
		})
	} else {
		sort.Slice(sortedTags, func(i, j int) bool {
			return sortedTags[i].Tag < sortedTags[j].Tag
		})
	}

	// Print the sorted tags and their counts
	for _, entry := range sortedTags {
		fmt.Printf("Tag: %s, Count: %d, Labels: %s\n", entry.Tag, entry.Info.Count, join(entry.Info.Labels, ", "))
	}
}

// Helper function to join string slices with a separator
func join(elems []string, sep string) string {
	result := ""
	for i, elem := range elems {
		if i > 0 {
			result += sep
		}
		result += elem
	}
	return result
}
