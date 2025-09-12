package events

type Publisher interface {
	Publish(eventType string, payload interface{}, metadata map[string]string) error
	Close() error
}
