package model

type Metadata struct{
	Title string `json:"title"`
	Name  string `json:"name"`
	Organization string `json:"organization"`
	Document string `json:"hash"`
	ExpirationDate string `json:"expirationDate"` 
}