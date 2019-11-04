package model

type ImageLinkable struct {
	ImageLink string `json:"imageLink"`
}

type Recipe struct {
	ImageLinkable

	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Headline    string       `json:"headline"`
	Description string       `json:"description"`
	Difficulty  int          `json:"difficulty"`
	PrepTime    string       `json:"prepTime"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Ingredient struct {
	ImageLinkable

	Name string `json:"name"`
}
