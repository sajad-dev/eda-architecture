package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// HandelFuncType defines the function signature for handling WebSocket requests
type HandelFuncType func(w http.ResponseWriter, r *http.Request, ws *Websocket)

// MiddlewareFuncType defines the function signature for HTTP middleware
type MiddlewareFuncType func(http.Handler) http.Handler

// TriggerBody represents the structure of a trigger event
type TriggerBody struct {
	Name     string      `json:"name"` // Name of the event
	Channel  string      `json:"channel,omitempty"` // Single channel name (if applicable)
	Channels []string    `json:"channels,omitempty"` // List of channels to broadcast to
	Data     interface{} `json:"data"` // Event payload
}

// Message represents a WebSocket message structure
type Message struct {
	Event   string      `json:"event"` // Event type
	Channel string      `json:"channel,omitempty"` // Associated channel (if any)
	Data    interface{} `json:"data"` // Message content
}

// MessageChan represents a message with connection details
type MessageChan struct {
	connection *websocket.Conn // WebSocket connection instance
	public_key string // Public key for authentication or identification
	message    *Message // Pointer to the message payload
}

// Websocket manages WebSocket connections
type Websocket struct {
	ServerMux CustomServeMux // Custom multiplexer for handling WebSocket routes
	Clients   map[string]map[string][]*websocket.Conn // Active WebSocket clients grouped by public key and channel
	CountConn int // Connection counter
}

// Addr represents a WebSocket route and its middlewares
type Addr struct {
	Pattern        string // Route pattern
	Handler        HandelFuncType // Associated handler function
	MiddlewareList []MiddlewareFuncType // List of middleware functions applied to the route
}

// CustomServeMux is a thread-safe HTTP request multiplexer
type CustomServeMux struct {
	Mux *http.ServeMux // HTTP multiplexer for handling requests
	Mu  sync.RWMutex   // Read-write mutex for thread safety
}

// WebSocket upgrader settings
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}
