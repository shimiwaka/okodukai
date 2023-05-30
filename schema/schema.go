package schema

import (
	"github.com/jinzhu/gorm"
)

type Board struct {
	gorm.Model `json:"-"`
	Owner      string         `json:"owner"`
	Token	   string		  `json:"token"`
}

type Column struct {
	gorm.Model `json:"-"`
	Board		string		`json:"board"`
	Name		string		`json:"name"`
}

type Response struct {
	gorm.Model `json:"-"`
	Owner		string		`json:"owner"`
	Token		string		`json:"token"`
	Columns		[]Column	`json:"columns"`
}