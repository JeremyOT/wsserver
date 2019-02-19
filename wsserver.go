package wsserver

import (
	"net/http"
	"strings"

	"github.com/JeremyOT/httpserver"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Server handles incoming requests and calls the appropriate handler for the
// request type (socket or simple request).
type Server struct {
	*httpserver.Server

	// HandleRequest is called to handle each incoming request that is not trying
	// to open WebSocket connection.
	HandleRequest http.HandlerFunc

	// HandleWebSocket is called to handle each new WebSocket connection. The connection
	// is closed when HandleSocket returns.
	HandleWebSocket func(*http.Request, *websocket.Conn)
}

func (s *Server) handleRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Connection") == "Upgrade" && strings.Contains(request.Header.Get("Upgrade"), "websocket") {
		s.handleWebsocket(writer, request)
		return
	}
	if s.HandleRequest != nil {
		s.HandleRequest(writer, request)
	}
}

func (s *Server) handleWebsocket(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	if s.HandleWebSocket != nil {
		s.HandleWebSocket(request, conn)
	}
}

// New creates a new Server for handling http requests and websockets
func New() *Server {
	s := &Server{}
	s.Server = httpserver.New(s.handleRequest)
	return s
}
