package colly

// RequestCallback is a type alias for OnRequest callback functions
type RequestCallback func(*Request)

// ResponseCallback is a type alias for OnResponse callback functions
type ResponseCallback func(*Response)

// CollectorCallback is a type alias for OnEvent callback functions
type CollectorCallback func(*Collector)

// FuncCallback is a type alias for OnFunc callback functions
type FuncCallback func(*Collector)

// EventCallback is a type alias for OnEvent callback functions
type EventCallback func(*Collector)

// ErrorCallback is a type alias for OnError callback functions
type ErrorCallback func(*Response, error)

// ScrapedCallback is a type alias for OnScraped callback functions
type ScrapedCallback func(*Response)
