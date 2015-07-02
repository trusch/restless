package websocket

import (
    "golang.org/x/net/websocket"
    "net/http"
    "../events"
)

var eventManager *events.EventManager = nil

func WebsocketHandler(ws *websocket.Conn) {
    id := eventManager.On("*@websocket",func(topic string, payload interface{}){
        ws.Write([]byte(topic))
    })
    data := make([]byte,1)
    ws.Read(data)
    // received bytes or connection died: close it.
    eventManager.Off(id)
    ws.Close()
}

func Setup(mgr *events.EventManager){
    eventManager = mgr
    http.Handle("/ws",websocket.Handler(WebsocketHandler))
}