package trader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

//上传图文消息内的图片获取URL【订阅号与服务号认证后均可用】
func (t *Trader) UpLoadImg(data []byte) (url string, err error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("media", "filename")
	if err != nil {
		return
	}
	_, err = io.Copy(fw, bytes.NewReader(data))
	if err != nil {
		return
	}
	w.Close()
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := MediaURL + "uploadimg?access_token=" + t.Accesstoken
	req, err := http.NewRequest("POST", surl, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	aaa, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	m := make(map[string]string)
	err = json.Unmarshal(aaa, &m)
	if err != nil {
		return
	}
	if _, ok := m["url"]; ok {
		url = m["url"]
	} else {
		err = errors.New(string(aaa))
	}
	return
}

//上传图文消息素材【订阅号与服务号认证后均可用】
func (t *Trader) UploadNews(articles []NewsArticle) (mediaId string, err error) {
	var a struct {
		Articles []NewsArticle `json:"articles"`
	}
	a.Articles = articles
	str, err := json.Marshal(a)
	if err != nil {
		return
	}
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := MediaURL + "uploadnews?access_token=" + t.Accesstoken
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		return
	}
	var r struct {
		Type    string `json:"type"`
		MediaId string `json:"media_id"`
		CreatAt int    `json:"created_at"`
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	mediaId = r.MediaId
	return
}

//根据标签进行群发【订阅号与服务号认证后均可用】
func (t *Trader) createMsgJson(msgtype string, tagId int, context string) (jsonstr string, err error) {
	type Filter struct {
		IsToAll bool `json:"is_to_all"`
		TagID   int  `json:"tag_id"`
	}
	switch msgtype {
	case mpnewsType:
		var c struct {
			Filter            Filter `json:"filter"`
			MpNews            MpNews `json:"mpnews"`
			MsgType           string `json:"msgtype"`
			SendIgnoreReprint int    `json:"send_ignore_reprint"`
		}
		c.MsgType, c.MpNews.MediaId = mpnewsType, context
		if tagId == 0 {
			c.Filter.IsToAll = true
		} else {
			c.Filter.IsToAll = false
			c.Filter.TagID = tagId
		}
		str, err := json.Marshal(c)
		return string(str), err
	case textType:
		var c struct {
			Filter  Filter `json:"filter"`
			Text    Text   `json:"text"`
			MsgType string `json:"msgtype"`
		}
		c.MsgType, c.Text.Content = textType, context
		if tagId == 0 {
			c.Filter.IsToAll = true
		} else {
			c.Filter.IsToAll = false
			c.Filter.TagID = tagId
		}
		str, err := json.Marshal(c)
		return string(str), err
	case voiceType:
		var c struct {
			Filter  Filter `json:"filter"`
			Voice   Voice  `json:"voice"`
			MsgType string `json:"msgtype"`
		}
		c.Voice.MediaId, c.MsgType = context, voiceType
		if tagId == 0 {
			c.Filter.IsToAll = true
		} else {
			c.Filter.IsToAll = false
			c.Filter.TagID = tagId
		}
		str, err := json.Marshal(c)
		return string(str), err
	case imageType:
		var c struct {
			Filter  Filter `json:"filter"`
			Image   Image  `json:"image"`
			MsgType string `json:"msgtype"`
		}
		c.Image.MediaId, c.MsgType = context, imageType
		if tagId == 0 {
			c.Filter.IsToAll = true
		} else {
			c.Filter.IsToAll = false
			c.Filter.TagID = tagId
		}
		str, err := json.Marshal(c)
		return string(str), err
	case MpVideoType:
		var c struct {
			Filter  Filter `json:"filter"`
			Mpvideo struct {
				MediaId string `json:"media_id"`
			} `json:"mpvideo"`
			MsgType string `json:"msgtype"`
		}
		c.Mpvideo.MediaId, c.MsgType = context, MpVideoType
		if tagId == 0 {
			c.Filter.IsToAll = true
		} else {
			c.Filter.IsToAll = false
			c.Filter.TagID = tagId
		}
		str, err := json.Marshal(c)
		return string(str), err
	case wxcardType:
		var c struct {
			Filter  Filter `json:"filter"`
			WxCard  WxCard `json:"wxcard"`
			MsgType string `json:"msgtype"`
		}
		c.WxCard.CardId, c.MsgType = context, wxcardType
		if tagId == 0 {
			c.Filter.IsToAll = true
		} else {
			c.Filter.IsToAll = false
			c.Filter.TagID = tagId
		}
		str, err := json.Marshal(c)
		return string(str), err
	}
	return
}
func (t *Trader) sendAll(msgtype string, tagId int, mediaId string) (s SendAllResp, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := SendAllURL + t.Accesstoken
	str, err := t.createMsgJson(msgtype, tagId, mediaId)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, str)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &s)
	if s.ErrCode != 0 {
		err = errors.New(string(b))
	}
	return
}

//群发图文
func (t *Trader) SendMpNewsAll(tagId int, mediaId string) (msgid, msgdataid int, err error) {
	m, err := t.sendAll(mpnewsType, tagId, mediaId)
	if err != nil {
		return
	}
	msgid, msgdataid = m.MsgId, m.MsgDataId
	return
}

//群发文本消息
func (t *Trader) SendTextAll(tagId int, text string) (msgid int, err error) {
	m, err := t.sendAll(textType, tagId, text)
	if err != nil {
		return
	}
	msgid = m.MsgId
	return
}

//群发图片
func (t *Trader) SendImageAll(tagId int, mediaId string) (msgid int, err error) {
	m, err := t.sendAll(imageType, tagId, mediaId)
	if err != nil {
		return
	}
	msgid = m.MsgId
	return
}

//群发语音
func (t *Trader) SendVoiceAll(tagId int, mediaId string) (msgid int, err error) {
	m, err := t.sendAll(voiceType, tagId, mediaId)
	if err != nil {
		return
	}
	msgid = m.MsgId
	return
}

//上传视频
func (t *Trader) uploadVideo(mediaId, title, description string) (newMediaId string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/media/uploadvideo?access_token=" + t.Accesstoken
	var c struct {
		MediaId     string `json:"media_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	c.MediaId, c.Title, c.Description = mediaId, title, description
	str, err := json.Marshal(c)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}
	if _, ok := m["media_id"]; ok {
		newMediaId = m["media_id"]
	} else {
		err = errors.New(string(b))
	}
	return
}

//群发视频
func (t *Trader) SendVideoAll(tagId int, mediaId, title, description string) (msgid int, err error) {
	newMediaId, err := t.uploadVideo(mediaId, title, description)
	if err != nil {
		return
	}
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	m, err := t.sendAll(MpVideoType, tagId, newMediaId)
	if err != nil {
		return
	}
	msgid = m.MsgId
	return
}

//群发卡券消息
func (t *Trader) SendWxCardAll(tagId int, cardId string) (msgid int, err error) {
	m, err := t.sendAll(wxcardType, tagId, cardId)
	if err != nil {
		return
	}
	msgid = m.MsgId
	return
}

//删除群发【订阅号与服务号认证后均可用】
// index 第一篇编号为1，该字段填0会删除全部文章
func (t *Trader) DeleteMass(msgid, index int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := DeleteSendALLURL + t.Accesstoken
	var c struct {
		MsgId      int `json:"msg_id"`
		ArticleIdx int `json:"article_idx"`
	}
	c.MsgId, c.ArticleIdx = msgid, index
	str, err := json.Marshal(c)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	var r Res
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	}
	return
}

//预览接口【订阅号与服务号认证后均可用】
func (t *Trader) Preview(msgtype string, useropenid, mediaId string) (err error) {
	var jsonstr []byte
	var p interface{}
	switch msgtype {
	case mpnewsType:
		p = PreviewMpNews{
			ToUser:  useropenid,
			MpNews:  MpNews{MediaId: mediaId},
			MsgType: mpnewsType,
		}

	case textType:
		p = PreviewText{
			ToUser:  useropenid,
			Text:    Text{Content: mediaId},
			MsgType: mpnewsType,
		}

	case videoType:
		p = PreviewVoice{
			ToUser:  useropenid,
			Voice:   Voice{MediaId: mediaId},
			MsgType: mpnewsType,
		}
	case imageType:
		p = PreviewImage{
			ToUser:  useropenid,
			Image:   Image{MediaId: mediaId},
			MsgType: mpnewsType,
		}
	case MpVideoType:
		p = PreviewMpVideo{
			ToUser:  useropenid,
			MpVideo: MpVideo{MediaId: mediaId},
			MsgType: mpnewsType,
		}
	default:
		err = errors.New("消息类型错误")
		return
	}
	jsonstr, err = json.Marshal(p)
	if err != nil {
		return
	}
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := PreviewURL + t.Accesstoken
	b, err := t.PostJson(surl, string(jsonstr))
	if err != nil {
		return
	}
	var r struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		MsgId   int    `json:"msg_id"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	}
	return
}

//查询群发消息发送状态【订阅号与服务号认证后均可用】
//status 消息发送后的状态，SEND_SUCCESS表示发送成功，SENDING表示发送中，SEND_FAIL表示发送失败，DELETE表示已删除
func (t *Trader) GetSendAllStatus(msgid int) (status string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := SendAllStatusURL + t.Accesstoken
	var m struct {
		MsgId string `json:"msg_id"`
	}
	m.MsgId = fmt.Sprint(msgid)
	str, err := json.Marshal(m)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	var r struct {
		MsgId  int    `json:"msg_id"`
		Status string `json:"msg_status"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.MsgId == msgid {
		status = r.Status
	} else {
		err = errors.New(string(b))
	}
	return
}

//获取群发速度
func (t *Trader) GetMassSpeed() (speedgrade, realspeed int, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := MassSpeedURL + "get?access_token=" + t.Accesstoken
	b, err := t.PostJson(surl, "")
	if err != nil {
		return
	}
	v := make(map[string]int)
	err = json.Unmarshal(b, &v)
	if err != nil {
		return
	}
	if _, ok := v["speed"]; ok {
		speedgrade = v["speed"]
		realspeed = v["realspeed"]
	} else {
		err = errors.New(string(b))
	}
	return
}

//设置群发速度 speed 只能是0到4的整数
func (t *Trader) SetMassSpeed(speed int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := MassSpeedURL + "set?access_token=" + t.Accesstoken
	b, err := t.PostJson(surl, `{"speed":`+fmt.Sprint(speed)+`}`)
	if err != nil {
		return
	}
	var r Res
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	}
	return
}
