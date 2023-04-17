package main

type SystemFields struct {
	Creator []string `json:"creator"`
	Title   []string `json:"title"`
}

type Fields struct {
	Fields SystemFields `json:"fields"`
}

type Settings struct {
	Section Fields `json:"section1"`
}
