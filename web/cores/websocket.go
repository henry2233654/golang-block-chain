package cores

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang-block-chain/services"
	"net/http"
)

const (
	Continue    = 1
	Success     = 2
	ServerError = 3
	ClientError = 4
)

var defaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketConnection(c *gin.Context) (*websocket.Conn, error) {
	conn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type WsResponse struct {
	Status int         `json:"status"`
	Title  string      `json:"title"`
	Data   interface{} `json:"data"`
}

func NewWsResponse(status int, title string, data interface{}) *WsResponse {
	return &WsResponse{
		Status: status,
		Title:  title,
		Data:   data,
	}
}

func NewImportedWsResponse(id uint) *WsResponse {
	imported := struct {
		ID uint `json:"id"`
	}{id}
	return NewWsResponse(Success, "imported", imported)
}

func NewErrorWsResponse(err error) *WsResponse {
	if _, ok := err.(services.ServiceError); ok {
		return NewWsResponse(ClientError, "error", err)
	}
	return NewWsResponse(ServerError, "error", err)
}
