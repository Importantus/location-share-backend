package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"location-share-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Zeit- und Bufferkonstanten
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Hier kannst du z.B. deine CheckOrigin-Funktion anpassen.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSMessage definiert das standardisierte Nachrichtenformat,
// das über den Websocket-Kanal gesendet wird.
type WSMessage struct {
	Type string      `json:"type"` // z.B. "location_update" oder "new_share"
	Data interface{} `json:"data"` // Payload, z.B. ein Location- oder Share-Objekt
}

// --- Interne Hub-Logik ---

// targetedMessage kapselt eine Nachricht, die an bestimmte Empfänger (UserIDs) gesendet werden soll.
// Ist Recipients leer, so wird die Nachricht an alle Clients gesendet.
type targetedMessage struct {
	Recipients []uuid.UUID
	Payload    []byte
}

// Hub verwaltet alle aktiven Clients und verteilt Nachrichten.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan targetedMessage
}

// newHub erstellt einen neuen Hub.
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan targetedMessage),
	}
}

// run verarbeitet Registrierung, Deregistrierung und Broadcast-Nachrichten.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				// Falls eine Empfängerliste definiert ist, prüfen, ob der Client darin enthalten ist.

				allowed := false
				for _, id := range msg.Recipients {
					if client.userID == id {
						allowed = true
						break
					}
				}
				if !allowed {
					continue
				}

				select {
				case client.send <- msg.Payload:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// --- Client-Implementierung ---

// Client repräsentiert einen Websocket-Client, der an den Hub registriert wird.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
}

// readPump liest Nachrichten von der Websocket-Verbindung.
// In dieser Implementierung ignorieren wir eingehende Nachrichten –
// hier könntest du aber auch z. B. Statusmeldungen oder ACKs verarbeiten.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket read error: %v", err)
			}
			break
		}
	}
}

// writePump schreibt Nachrichten aus dem send-Channel in die Websocket-Verbindung.
// Zudem werden periodisch Ping-Nachrichten gesendet, um die Verbindung aktiv zu halten.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			// Falls weitere Nachrichten im Channel liegen, werden diese angehängt.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte("\n"))
				_, _ = w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// --- Singleton Hub ---

var (
	hubInstance *Hub
	hubOnce     sync.Once
)

// GetHub liefert die Singleton-Instanz des Hubs.
func GetHub() *Hub {
	hubOnce.Do(func() {
		hubInstance = newHub()
		go hubInstance.run()
	})
	return hubInstance
}

// --- Öffentliche Schnittstellen ---

// WSHandler liefert einen Gin-Handler, der die Websocket-Verbindung aufbaut und den Client beim Hub registriert.
// Die Funktion getUserID muss den aktuellen User (als uuid.UUID) aus dem Context extrahieren –
// z. B. mithilfe deiner Auth-Middleware.
func WSHandler(getUserID func(c *gin.Context) (uuid.UUID, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Websocket upgrade error: %v", err)
			return
		}
		client := &Client{
			hub:    GetHub(),
			conn:   conn,
			send:   make(chan []byte, 256),
			userID: userID,
		}
		client.hub.register <- client
		go client.writePump()
		go client.readPump()
	}
}

// BroadcastLocation erstellt eine standardisierte Nachricht mit dem Typ "location_update"
// und sendet sie an die angegebenen Empfänger (userIDs).
func BroadcastLocation(userIDs []uuid.UUID, location models.Location) {
	msg := WSMessage{
		Type: "location_update",
		Data: location,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling location update: %v", err)
		return
	}
	GetHub().broadcast <- targetedMessage{
		Recipients: userIDs,
		Payload:    payload,
	}
}

// BroadcastNewShare erstellt eine Nachricht mit dem Typ "new_share"
// (z. B. wenn ein neuer Standort geteilt wird) und sendet sie an die angegebenen Empfänger.
// Hier wird angenommen, dass models.Share den entsprechenden Datentyp repräsentiert.
func BroadcastNewShare(userIDs []uuid.UUID, share models.Share) {
	msg := WSMessage{
		Type: "share_create",
		Data: share,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling new share: %v", err)
		return
	}
	GetHub().broadcast <- targetedMessage{
		Recipients: userIDs,
		Payload:    payload,
	}
}

// BroadcastShareDeleted erstellt eine Nachricht mit dem Typ "share_deleted"
// (z. B. wenn ein Standort geteilt wird) und sendet sie an die angegebenen Empfänger.
// Hier wird angenommen, dass models.Share den entsprechenden Datentyp repräsentiert.
func BroadcastShareDeleted(userIDs []uuid.UUID, shareID uuid.UUID) {
	msg := WSMessage{
		Type: "share_delete",
		Data: shareID,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling share deleted: %v", err)
		return
	}
	GetHub().broadcast <- targetedMessage{
		Recipients: userIDs,
		Payload:    payload,
	}
}
