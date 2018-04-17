package tv

import (
	"testing"
)

func TestLoadingDocument(t *testing.T) {
	url := TestHTMLURL("testhtml/food1.html")
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil", idr.idoc)
	}
	url = TestHTMLURL("testhtml/films1.html")
	go url.loadDocument(c)
	idr = <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil: ", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil: ", idr.idoc)
	}
}
func TestIplayerSelectionResults(t *testing.T) {
	url := TestHTMLURL("testhtml/films1.html")
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	idr := <-c
	sel := iplayerSelection{idr.idoc.doc.Find(".list-item-inner")}
	selres := sel.selectionResults()
	if len(selres) != 20 {
		t.Error("Expected length of selectionresults to equal: ", len(selres))
	}
	progpage := selres[0]
	if progpage.prog != nil {
		t.Error("Expected proramme to be nil: ", progpage.prog)
	}
	if progpage.programPage != "testhtml/adam_curtis.html" {
		t.Error("Expected program Page to be 'testhtml/adam_curtis.html' not: ", progpage.programPage)
	}
	if selres[1].prog.Title != "A Hijacking" {
		t.Error("Expected second programme title to be 'A Hijacking', got: ", selres[1].prog.Title)
	}
	if selres[1].programPage != "" {
		t.Error("Expected second programPage to be an empty string, got: ", selres[1].programPage)
	}
}

func TestCollectPages(t *testing.T) {
	url := TestHTMLURL("testhtml/films1.html")
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	docres := <-c
	if docres.Error != nil {
		t.Error("Expected error in documentresult to be nil, got: ", docres.Error)
	}
	tid := TestIplayerDocument{docres.idoc}
	np := tid.nextPages()
	if len(np) != 1 {
		t.Error("Expected length of nextPages to be 1, got: ", len(np))
	}
	cp := collectPages(np)
	if len(cp) != 1 {
		t.Error("Expected length of collectedPages to be 1, got: ", len(cp))
	}
	if cp[0].Error != nil {
		t.Error("Expected error for first doc in collected Pages to be nil, got: ", cp[0].Error)
	}
}

func TestProgramPage(t *testing.T) {
	url := TestHTMLURL("testhtml/classic_mary_berry.html")
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	docres := <-c
	if docres.Error != nil {
		t.Error("Expected document to load without error, git: ", docres.Error)
	}
}

func TestNewMainCategory(t *testing.T) {
	url := TestHTMLURL("testhtml/films1.html")
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	docres := <-c
	if docres.Error != nil {
		t.Error("Expected error in documentresult to be nil, got: ", docres.Error)
	}
	tid := TestIplayerDocument{docres.idoc}
	nmd := newMainCategory(&tid)
	if len(nmd.nextdocs) != 2 {
		t.Error("Expected length of nextdocs to be 2, got: ", len(nmd.nextdocs))
	}
	pp := programPage{nmd.nextdocs[1]}
	progs := pp.programmes()
	if len(progs) == 0 {
		t.Error("Expected length of programmes > 0, got: ", len(progs))
	}
}

