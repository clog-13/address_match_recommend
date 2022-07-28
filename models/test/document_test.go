package test

import (
	. "address_match_recommend/models"
	"fmt"
	"gorm.io/gorm/clause"
	"testing"
)

func TestInsertDoc(t *testing.T) {
	did := uint(24)
	doc := &Document{
		Terms:        make([]*Term, 0),
		RoadNumValue: 12,
	}

	doc.TermsId = did
	doc.Terms = append(doc.Terms, &Term{
		Text:       "terms_term",
		Types:      12,
		Idf:        12,
		DocumentID: did,
		Ref: &Term{
			Text:  "ref",
			Types: 13,
			Idf:   13,
		},
	})

	doc.TownId = did
	doc.Town = &Term{
		Text:       "downtown",
		Types:      10,
		Idf:        10,
		DocumentID: did,
		Ref: &Term{
			Text:  "downtown_ref",
			Types: 13,
			Idf:   13,
		},
	}

	DB.Create(doc)
}

func TestQueryAll(t *testing.T) {
	var docs []Document
	DB.Preload(clause.Associations).Find(&docs)
	fmt.Println(docs[0].RoadNumValue)
	fmt.Println(docs[0].Town.Text)
	fmt.Println(docs[0].Town.Ref.Text)
	fmt.Println(docs[0].Terms[0].Text)
	fmt.Println(docs[0].Terms[0].Ref.Text)
}
