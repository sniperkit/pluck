package colly

// HTMLCallback is a type alias for OnHTML callback functions
type HTMLCallback func(*HTMLElement)

// XMLCallback is a type alias for OnXML callback functions
type XMLCallback func(*XMLElement)

// JSONCallback is a type alias for OnJSON callback functions
type JSONCallback func(*JSONElement)

// TABCallback is a type alias for OnTAB callback functions
type TABCallback func(*TABElement)
