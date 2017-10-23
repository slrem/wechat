package wechat

import (
	"fmt"
	"sync"
)

// type TextHandler func(c Context, message TextMessage) error
// type ImageHandler func(c Context, message ImageMessage) error
// type VoiceHandler func(c Context, message VoiceMessage) error
// type ShortVideoHandler func(c Context, message ShortVideoMessage) error
// type LocationHandler func(c Context, message LocationMessage) error
// type LinkHandler func(c Context, message LinkMessage) error
// type SubscribeEventHandler func(c Context, message SubscribeEventMessage) error
// type UnsubscribeEventHandler func(c Context, message UnsubscribeEventMessage) error
// type ScanSubscribeEventHandler func(c Context, message ScanSubscribeEventMessage) error
// type ScanEventHandler func(c Context, message ScanEventMessage) error
// type LocationEventHandler func(c Context, message LocationEventMessage) error
// type MenuViewEventHandler func(c Context, message MenuViewEventMessage) error
// type MenuClickEventHandler func(c Context, message MenuClickEventMessage) error

type Route struct {
	MsgType MsgType
	Key     string
	Handler Handler
}

type Router struct {
	routes map[string]Route
	mtx    sync.Mutex
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Route),
	}
}

func routeKey(msgType MsgType, key string) string {
	return fmt.Sprintf("%d%s", msgType, key)
}

func (r *Router) Get(msgType MsgType, key string) Handler {
	k := routeKey(msgType, key)

	r.mtx.Lock()
	defer r.mtx.Unlock()

	route, ok := r.routes[k]
	if !ok {
		return nil
	}

	return route.Handler
}

func (r *Router) Add(msgType MsgType, key string, route Route) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.routes[routeKey(msgType, key)] = route
}

func (r *Router) Find(c Context) {
	msgType := c.Request().MsgType()
	key := ""

	if h := r.Get(msgType, key); h != nil {
		c.SetHandler(h)
	} else {
		c.SetHandler(c.Wechat().DefaultHandler())
	}

}
