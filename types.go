package main

import "sync"

// Minimal structs for Ticketmaster response (only fields we need)
type TMResponse struct {
	Embedded *struct {
		Events []Event `json:"events"`
	} `json:"_embedded"`
	Page *Page `json:"page"`
}

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

type Event struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Dates struct {
		Start struct {
			LocalDate string `json:"localDate"`
			DateTime  string `json:"dateTime"`
		} `json:"start"`
	} `json:"dates"`
	// Add more fields as needed
}

// in-memory store (map marketplace -> []Event)
var store = struct {
	sync.RWMutex
	data map[string][]Event
}{
	data: map[string][]Event{},
}
