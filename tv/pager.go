package tv

type Pager interface {
	loadDocument(chan<- *iplayerDocumentResult)
}

type NextPager interface {
	nextPages() []interface{}
	programPages() []interface{}
}