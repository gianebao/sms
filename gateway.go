package sms

// Gateway represents a generic SMS gateway
type Gateway interface {
	send(to string, message Message, callback string) (interface{}, error)
}

// Send creates a message using the provided `gateway`
func Send(g Gateway, to string, message Message, callback string) (interface{}, error) {
	return g.send(to, message, callback)
}
