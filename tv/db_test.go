package tv

import (
	"testing"
	"time"
)

func TestProgrammeDB_Index(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := &TestIplayerDocument{doc}
	cat := NewCategory("films", td)
	pdb := &ProgrammeDB{[]*Category{cat}, time.Now()}
	pdb.index()
	fp := pdb.Categories[0].Programmes[0]
	if fp.Index != 0 {
		t.Errorf("Expected first programme to have index '0', got: %d", fp.Index)
	}
	for _, i := range pdb.Categories[0].Programmes[1:] {
		if !(i.Index > 0) {
			t.Errorf("Expected for title %q index to be '> 0', got: %d", i.Title, i.Index)
		}
	}

}

func TestProgrammeDB_Save(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := &TestIplayerDocument{doc}
	cat := NewCategory("films", td)
	pdb := &ProgrammeDB{[]*Category{cat}, time.Now()}
	err := pdb.Save("testdb.json")
	if err != nil {
		t.Error("Expected saving db should not return error, got: ", err)
	}
}