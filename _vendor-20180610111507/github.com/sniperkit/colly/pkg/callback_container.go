package colly

// htmlCallbackContainer
type htmlCallbackContainer struct {
	Selector string
	Function HTMLCallback
}

// xmlCallbackContainer
type xmlCallbackContainer struct {
	Query    string
	Function XMLCallback
}

// jsonCallbackContainer
type jsonCallbackContainer struct {
	Query    string
	Function JSONCallback
}

// tabCallbackContainer
type tabCallbackContainer struct {
	Query    string
	Function TABCallback
}
