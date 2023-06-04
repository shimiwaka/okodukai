package schema

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Board struct {
	gorm.Model `json:"-"`
	Owner      string         `json:"owner"`
	Token	   string		  `json:"token"`
}

type Column struct {
	gorm.Model `json:"-"`
	Board		uint		`json:"board"`
	Name		string		`json:"name"`
	Price		int			`json:"price"`
}

type Day struct {
	gorm.Model `json:"-"`
	Date		time.Time	`json:"date"`
	Checked		[]bool		`json:"checked"`
	Payment		int			`json:"payment"`
}

type Check struct {
	gorm.Model `json:"-"`
	Date		time.Time	`json:"date"`
	Column		uint		`json:"column"`
	Board		uint 		`json:"board"`
}

type Response struct {
	gorm.Model `json:"-"`
	Owner		string		`json:"owner"`
	Token		string		`json:"token"`
	Columns		[]Column	`json:"columns"`
	Days		[]Day		`json:"days"`
}

type Payment struct {
	gorm.Model `json:"-"`
	Board		uint 		`json:"board"`
	Date		time.Time	`json:"date"`
}