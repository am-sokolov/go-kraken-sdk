package websocket

import "time"

// MessageEvent represents an inbound WebSocket message and its local receive timestamps.
//
// Raw must be treated as immutable by the caller.
type MessageEvent struct {
	Raw []byte

	// ReceivedAt is the wall-clock time at which the message was read from the socket.
	ReceivedAt time.Time

	// ReceivedMonoNs is a monotonic timestamp in nanoseconds, relative to the client's creation time.
	ReceivedMonoNs int64

	// Parsed indicates whether the message was successfully parsed into a BaseMessage.
	Parsed bool

	// ParseError contains the parsing error message when Parsed is false.
	ParseError string

	// Message is the parsed message when Parsed is true.
	Message BaseMessage
}
