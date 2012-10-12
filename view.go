package main

import (
    "net/http"
    "code.google.com/p/go.net/websocket"
    "github.com/araddon/httpstream"
    "io"
    "crypto/sha1"
    "encoding/base64"
    "fmt"
)

type View struct {
    Id string
    tweet chan *httpstream.Tweet
    regist chan *websocket.Conn
    done chan bool
    conns []*websocket.Conn
}

func NewView(seed string) *View {

    h := sha1.New()
    io.WriteString(h, seed)
    id := base64.StdEncoding.EncodeToString(h.Sum(nil))

    tweet := make(chan *httpstream.Tweet)
    regist := make(chan *websocket.Conn)
    done := make(chan bool)
    conns := make([]*websocket.Conn, 0, 1000)

    view := &View{
        id,
        tweet,
        regist,
        done,
        conns,
    }

    stop := func(w http.ResponseWriter, r *http.Request) {
        done <- true
    }

    onConnected := func(ws *websocket.Conn) {
        regist <- ws
    }

    show := func(w http.ResponseWriter, r *http.Request) {
        templates.ExecuteTemplate(w, "view", view)
    }

    mux.Handle(fmt.Sprintf("/view/%s/ws", id), websocket.Handler(onConnected))
    mux.HandleFunc(fmt.Sprintf("/view/%s/stop", id), show)
    mux.HandleFunc(fmt.Sprintf("/view/%s/show", id), stop)

   return view
}

func (view *View) Start() {
    go func() {
        for {
            select {
            case ws := <-view.regist:
                view.conns = append(view.conns, ws)
            case t := <-view.tweet:
                for _, ws := range view.conns {
                    if ws != nil {
                        websocket.JSON.Send(ws, t)
                    }
                }
            case <-view.done:
                break
            }
        }
    }()
}
