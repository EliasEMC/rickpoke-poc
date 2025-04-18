package model

type Rick struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Pokemon struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
