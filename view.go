package main

import (
    "net/http"
    "code.google.com/p/go.net/websocket"
    "github.com/araddon/httpstream"
    "net/url"
)

type View struct {
    tweet chan *httpstream.Tweet
    regist chan *websocket.Conn
    done chan bool
    stop func(w http.ResponseWriter, r *http.Request)
    onConnected func(ws *websocket.Conn)
    conns []*websocket.Conn
    URL *url.URL
}

func NewView(url *url.URL) *View {

    tweet := make(chan *httpstream.Tweet)
    regist := make(chan *websocket.Conn)
    done := make(chan bool)
    conns := make([]*websocket.Conn, 0, 1000)

    stop := func(w http.ResponseWriter, r *http.Request) {
        done <- true
    }

    onConnected := func(ws *websocket.Conn) {
        regist <- ws
    }

    view := &View{
        tweet,
        regist,
        done,
        stop,
        onConnected,
        conns,
        url,
    }
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
