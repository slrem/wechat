package wechat

import (
	"encoding/xml"
	"errors"
	"net/http"
	"time"
)

type defaultResponse struct {
	c Context
	w http.ResponseWriter
}

func (dr defaultResponse) String(s string) (err error) {
	err = dr.Bytes([]byte(s))
	return
}

func (dr defaultResponse) Success() (err error) {
	err = dr.Bytes([]byte("success"))
	return
}

func (dr defaultResponse) Bytes(b []byte) (err error) {
	_, err = dr.w.Write(b)
	return
}

func (dr defaultResponse) Response(data interface{}) (err error) {
	b, err := xml.Marshal(data)
	if err != nil {
		return
	}

	if dr.c.Wechat().securityMode {
		b, err = dr.c.Wechat().Encrypt(b)
		if err != nil {
			return
		}
	}
	return dr.Bytes(b)
}

func (dr defaultResponse) Text(content string) error {
	return dr.Response(NewTextResponseMessage(
		dr.c.Request().FromUserName(),
		dr.c.Request().ToUserName(),
		content,
	))
}

func (dr defaultResponse) Image(mediaId string) error {
	return dr.Response(NewImageResponseMessage(
		dr.c.Request().FromUserName(),
		dr.c.Request().ToUserName(),
		mediaId,
	))
}

func (dr defaultResponse) Voice(mediaId string) error {
	return dr.Response(NewVoiceResponseMessage(
		dr.c.Request().FromUserName(),
		dr.c.Request().ToUserName(),
		mediaId,
	))
}

func (dr defaultResponse) Video(video Video) error {
	return dr.Response(NewVideoResponseMessage(
		dr.c.Request().FromUserName(),
		dr.c.Request().ToUserName(),
		video,
	))
}

func (dr defaultResponse) Music(music Music) error {
	return dr.Response(NewMusicResponseMessage(
		dr.c.Request().FromUserName(),
		dr.c.Request().ToUserName(),
		music,
	))
}

func (dr defaultResponse) Article(articles ...ArticleItem) (err error) {
	article, err := NewArticleResponseMessage(
		dr.c.Request().FromUserName(),
		dr.c.Request().ToUserName(),
		articles...,
	)
	if err != nil {
		return
	}

	return dr.Response(article)
}

type CDATAString struct {
	CDATA string `xml:",cdata"`
}

type MediaId struct {
	MediaId CDATAString
}

type TextResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      CDATAString
}

func NewTextResponseMessage(to, from, content string) TextResponseMessage {
	return TextResponseMessage{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      CDATAString{CDATA: content},
	}
}

type ImageResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Image        MediaId
}

func NewImageResponseMessage(to, from, mediaId string) ImageResponseMessage {
	return ImageResponseMessage{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		Image:        MediaId{MediaId: CDATAString{CDATA: mediaId}},
	}
}

type VoiceResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Voice        MediaId
}

func NewVoiceResponseMessage(to, from, mediaId string) VoiceResponseMessage {
	return VoiceResponseMessage{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "voice",
		Voice:        MediaId{MediaId: CDATAString{CDATA: mediaId}},
	}
}

type Video struct {
	MediaId     CDATAString
	Title       CDATAString
	Description CDATAString
}

func NewVideo(mediaId, title, description string) Video {
	return Video{
		MediaId:     CDATAString{CDATA: mediaId},
		Title:       CDATAString{CDATA: title},
		Description: CDATAString{CDATA: description},
	}
}

type VideoResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Video        Video
}

func NewVideoResponseMessage(to, from string, video Video) VideoResponseMessage {
	return VideoResponseMessage{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "video",
		Video:        video,
	}
}

type Music struct {
	Title        CDATAString
	Description  CDATAString
	MusicUrl     CDATAString
	HQMusicUrl   CDATAString
	ThumbMediaId string
}

func NewMusic(title, description, musicUrl, HQMusicUrl, thumbMediaId string) Music {
	return Music{
		Title:        CDATAString{CDATA: title},
		Description:  CDATAString{CDATA: description},
		MusicUrl:     CDATAString{CDATA: musicUrl},
		HQMusicUrl:   CDATAString{CDATA: HQMusicUrl},
		ThumbMediaId: thumbMediaId,
	}
}

type MusicResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Music        Music
}

func NewMusicResponseMessage(to, from string, music Music) MusicResponseMessage {
	return MusicResponseMessage{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "music",
		Music:        music,
	}
}

type Article struct {
	Title       CDATAString
	Description CDATAString
	PicUrl      CDATAString
	Url         CDATAString
}

type ArticleItem struct {
	Item Article `xml:"item"`
}

func NewArticleItem(title, description, picUrl, url string) ArticleItem {
	return ArticleItem{
		Item: Article{
			Title:       CDATAString{CDATA: title},
			Description: CDATAString{CDATA: description},
			PicUrl:      CDATAString{CDATA: picUrl},
			Url:         CDATAString{CDATA: url},
		},
	}
}

type ArticleResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	ArticleCount int
	Articles     []ArticleItem
}

var ArticleCountOverError = errors.New("article count over 10")

func NewArticleResponseMessage(to, from string, articles ...ArticleItem) (a ArticleResponseMessage, err error) {
	count := len(articles)

	if count > 10 {
		err = ArticleCountOverError
		return
	}

	a = ArticleResponseMessage{
		ToUserName:   to,
		FromUserName: from,
		CreateTime:   time.Now().Unix(),
		MsgType:      "news",
		ArticleCount: count,
		Articles:     articles,
	}
	return
}
