package notification

// Formatter is a message formatter
type Formatter interface {
	Format(msg interface{}) string
}
