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
	Board		string		`json:"board"`
	Name		string		`json:"name"`
}

type Day struct {
	gorm.Model `json:"-"`
	Date		time.Time	`json:"date"`
	Checked		[]bool		`json:"checked"`
}

type Checked struct {
	gorm.Model `json:"-"`
	Date		time.Time	`json:"date"`
	Column		int			`json:"column"`
}

type Response struct {
	gorm.Model `json:"-"`
	Owner		string		`json:"owner"`
	Token		string		`json:"token"`
	Columns		[]Column	`json:"columns"`
	Days		[]Day		`json:"days"`
}