package models

import (
	"fmt"
	"gorm.io/gorm/clause"
	"testing"
)

func TestInsertDoc(t *testing.T) {
	doc := &Document{
		RoadNumValue: 5,
		TownId:       13,
		RoadId:       14,
	}

	doc.Town = &Term{
		Text:   "downtown",
		TermId: 13,
	}
	doc.Road = &Term{
		Text:   "deadroad",
		TermId: 14,
	}
	DB.Create(doc)
}

func TestQueryAll(t *testing.T) {
	var docs []Document
	DB.Preload(clause.Associations).Find(&docs)
	//DB.Preload("Town").Preload("Village").Find(&docs)
	fmt.Println(docs[0])
	fmt.Println(docs[0].Town.Text)
	fmt.Println(docs[0].Village)
	fmt.Println(docs[0].Road.Text)
}
