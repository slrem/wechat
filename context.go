package wechat

import "net/http"

type Context interface {
	Wechat() *Wechat
	Request() Request
	Response() Response
	SetHandler(h Handler)
}

type Request interface {
	ToUserName() string
	FromUserName() string
	CreateTime() int
	MsgType() MsgType
	Content() string
	MsgId() int64
	PicUrl() string
	MediaId() string
	Format() string
	Recognition() string
	ThumbMediaId() string
	LocationX() float64
	LocationY() float64
	Scale() int
	Label() string
	Title() string
	Description() string
	Url() string
	Event() string
	EventKey() string
	Ticket() string
	Latitude() float32
	Longitude() float32
	Precision() float32

	MenuId() int64
	ScanCodeInfo() ScanCodeInfo
	SendPicsInfo() SendPicsInfo
	SendLocationInfo() SendLocationInfo

	Status() string
}

type Response interface {
	Success() error
	String(s string) error
	Bytes(b []byte) error
	Response(data interface{}) error
	Text(content string) error
	Image(mediaID string) error
	Voice(mediaID string) error
	Video(video Video) error
	Music(music Music) error
	Article(articles ...ArticleItem) error
}

type context struct {
	dft *defaultRequestMessage
	r   *http.Request
	w   http.ResponseWriter
	wc  *Wechat
	dr  defaultResponse

	handler Handler
}

func newContext(w http.ResponseWriter, r *http.Request, wc *Wechat) (c *context) {
	c = &context{
		r:   r,
		w:   w,
		wc:  wc,
		dft: &defaultRequestMessage{},
	}

	c.dr = defaultResponse{c: c, w: w}

	return
}

func (c *context) parse() (err error) {
	data, err := c.wc.body(c.r)
	if err != nil {
		return
	}

	err = c.dft.Unmarshal(data)
	return
}

func (c *context) Handler() Handler {
	return c.handler
}

func (c *context) SetHandler(h Handler) {
	c.handler = h
}

func (c *context) Request() Request {
	return c.dft
}

func (c *context) Wechat() *Wechat {
	return c.wc
}

func (c *context) Response() Response {
	return c.dr
}
