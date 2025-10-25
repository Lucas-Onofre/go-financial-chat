package broker

type Consumer interface {
	Subscribe(queue string, handler func(message string) error) error
	Close() error
}

type Producer interface {
	Publish(queue string, message string) error
	Close() error
}
