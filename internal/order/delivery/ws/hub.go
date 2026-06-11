package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/internal/order"
	"kerjadekat/backend/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	writeWait      = 10 * time.Second
	authWait       = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type authFrame struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type client struct {
	userID uuid.UUID
	role   string
	conn   *websocket.Conn
	send   chan []byte
}

// Hub routes Redis Pub/Sub order events to authenticated worker WebSockets.
type Hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[*client]struct{}

	tokens *token.Issuer
	rdb    *redis.Client
}

func NewHub(tokens *token.Issuer, rdb *redis.Client) *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]map[*client]struct{}),
		tokens:  tokens,
		rdb:     rdb,
	}
}

func (h *Hub) register(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.userID] == nil {
		h.clients[c.userID] = make(map[*client]struct{})
	}
	h.clients[c.userID][c] = struct{}{}
}

func (h *Hub) unregister(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.userID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, c.userID)
		}
	}
	close(c.send)
}

func (h *Hub) sendToUsers(userIDs []uuid.UUID, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, id := range userIDs {
		for cl := range h.clients[id] {
			select {
			case cl.send <- payload:
			default:
			}
		}
	}
}

// RunSubscribe listens on orders:new and orders:status, fanning out to connections.
func (h *Hub) RunSubscribe(ctx context.Context) {
	sub := h.rdb.Subscribe(ctx, order.ChannelOrdersNew, order.ChannelOrdersStatus)
	ch := sub.Channel()
	go func() {
		defer sub.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				switch msg.Channel {
				case order.ChannelOrdersNew:
					var evt order.NewOrderAvailable
					if err := json.Unmarshal([]byte(msg.Payload), &evt); err != nil {
						log.Printf("ws hub: bad new_order payload: %v", err)
						continue
					}
					payload, err := json.Marshal(evt)
					if err != nil {
						continue
					}
					h.sendToUsers(evt.WorkerUserIDs, payload)

				case order.ChannelOrdersStatus:
					var evt order.OrderStatusChanged
					if err := json.Unmarshal([]byte(msg.Payload), &evt); err != nil {
						log.Printf("ws hub: bad status payload: %v", err)
						continue
					}
					payload, err := json.Marshal(evt)
					if err != nil {
						continue
					}
					recipients := []uuid.UUID{evt.ConsumerID}
					if evt.WorkerID != nil {
						recipients = append(recipients, *evt.WorkerID)
					}
					h.sendToUsers(recipients, payload)
				}
			}
		}
	}()
}

// HandleWS upgrades GET /api/v1/ws — JWT must arrive in the first post-handshake frame.
func (h *Hub) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	cl, ok := h.authenticateFirstFrame(conn)
	if !ok {
		_ = conn.Close()
		return
	}

	h.register(cl)
	go h.writePump(cl)
	go h.readPump(cl)
}

func (h *Hub) authenticateFirstFrame(conn *websocket.Conn) (*client, bool) {
	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(authWait))

	_, data, err := conn.ReadMessage()
	if err != nil {
		h.writeAuthFailure(conn, "missing auth frame")
		return nil, false
	}

	var frame authFrame
	if err := json.Unmarshal(data, &frame); err != nil || frame.Type != "auth" || frame.Token == "" {
		h.writeAuthFailure(conn, "invalid auth frame")
		return nil, false
	}

	claims, err := h.tokens.ParseAccess(frame.Token)
	if err != nil {
		h.writeAuthFailure(conn, "invalid token")
		return nil, false
	}
	if claims.Role != domain.RoleWorker {
		h.writeAuthFailure(conn, "workers only")
		return nil, false
	}

	_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteJSON(map[string]string{"type": "auth_ok"}); err != nil {
		return nil, false
	}

	return &client{
		userID: claims.UserID,
		role:   claims.Role,
		conn:   conn,
		send:   make(chan []byte, 32),
	}, true
}

func (h *Hub) writeAuthFailure(conn *websocket.Conn, reason string) {
	_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
	_ = conn.WriteJSON(map[string]string{
		"type":    "auth_error",
		"message": reason,
	})
	_ = conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.ClosePolicyViolation, reason),
	)
}

func (h *Hub) readPump(cl *client) {
	defer func() {
		h.unregister(cl)
		_ = cl.conn.Close()
	}()
	cl.conn.SetReadLimit(maxMessageSize)
	_ = cl.conn.SetReadDeadline(time.Now().Add(pongWait))
	cl.conn.SetPongHandler(func(string) error {
		return cl.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		if _, _, err := cl.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Hub) writePump(cl *client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = cl.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-cl.send:
			_ = cl.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = cl.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := cl.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = cl.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cl.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
