package rickandmorty

import "time"

type ResponseCaracters struct {
	Info struct {
		Count int    `json:"count"`
		Pages int    `json:"pages"`
		Next  string `json:"next"`
		Prev  any    `json:"prev"`
	} `json:"info"`
	Results []Caracter `json:"results"`
}

type Caracter struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
	Origin  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"origin"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Image   string    `json:"image"`
	Episode []string  `json:"episode"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}
