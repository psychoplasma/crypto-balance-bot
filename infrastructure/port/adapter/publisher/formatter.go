package publisher

// Formatter is a message formatter to translate the input
// to the corresponding publishing service according to its requirements
type Formatter func(msg interface{}) string
