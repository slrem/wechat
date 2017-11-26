package trader

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Trader struct {
	AppId                string
	AppSecret            string
	Accesstoken          string
	ExpiresIn            int64
	JsapiTicket          string
	JsapiTicketExpiresIn int64
	mtx                  sync.Mutex
	AccessTokenHandler   Handler
}
type Handler func() (AccessToken, error)

func NewTrader(appid, appsecret string, h Handler) (t *Trader, err error) {
	t = &Trader{
		AppId:              appid,
		AppSecret:          appsecret,
		AccessTokenHandler: h,
	}
	a, err := t.GetAccessToken()
	t.mtx.Lock()
	t.Accesstoken, t.ExpiresIn = a.Access_token, a.Expires_in+time.Now().Unix()
	t.mtx.Unlock()
	return
}

func (t *Trader) GetAccessToken() (a AccessToken, err error) {
	h := t.AccessTokenHandler
	if h != nil {
		a, err = h()
		return
	}
	surl := AccessTokenURL + t.AppId + "&secret=" + t.AppSecret
	res, err := http.Get(surl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &a)
	if err != nil {
		return
	}
	if a.Access_token == "" {
		err = errors.New(string(b))
		return
	}

	return
}

func (t *Trader) CheckAccessTokenLive() (err error) {
	h := t.AccessTokenHandler
	if h != nil {
		return
	}
	var a AccessToken
	if time.Now().Unix() > t.ExpiresIn-300 {
		a, err = t.GetAccessToken()
		t.mtx.Lock()
		t.Accesstoken, t.ExpiresIn = a.Access_token, a.Expires_in+time.Now().Unix()
		t.mtx.Unlock()
	}
	return
}

func (t *Trader) FlushAccessToken() (err error) {
	a, err := t.GetAccessToken()
	if err != nil {
		return
	}
	t.mtx.Lock()
	t.Accesstoken = a.Access_token
	t.mtx.Unlock()
	return
}

func (t *Trader) upload(materialtype string, data []byte, title, introduction string) (mediaId, url string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	str := "filename.jpg"
	s := "media"

	if materialtype == "video" {
		str = "video.mp4"
		// s = "video"
	}
	fw, err := w.CreateFormFile(s, str)
	if err != nil {
		return
	}

	_, err = io.Copy(fw, bytes.NewReader(data))
	if err != nil {
		return
	}

	if materialtype == "video" {
		vw, _ := w.CreateFormField("description")

		videoDesc := &VideoDesc{
			Title:        title,
			Introduction: introduction,
		}
		str, _ := json.Marshal(videoDesc)
		io.Copy(vw, bytes.NewReader(str))
	}

	w.Close()
	surl := UploadURL + t.Accesstoken + `&type=` + materialtype
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
	if strings.Contains(string(aaa), "errcode") {
		err = errors.New(string(aaa))
		return
	}
	m := make(map[string]string)
	err = json.Unmarshal(aaa, &m)
	if err != nil {
		return
	}
	if _, ok := m["media_id"]; ok {
		mediaId = m["media_id"]
		url = m["url"]
	} else {
		err = errors.New(string(aaa))
	}
	return
}

func (t *Trader) Get(surl string) (b []byte, err error) {
	resp, err := http.Get(surl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	return
}

func (t *Trader) PostJson(surl, jsonstr string) (b []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", surl, bytes.NewReader([]byte(jsonstr)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()
	b, err = ioutil.ReadAll(r.Body)
	return
}

func (t *Trader) AddImageMaterial(data []byte) (mediaId, url string, err error) {
	return t.upload("image", data, "", "")
}

func (t *Trader) AddVoiceMaterial(data []byte) (mediaId string, err error) {
	mediaId, _, err = t.upload("voice", data, "", "")
	return
}

func (t *Trader) AddVideoMaterial(data []byte, title, introduction string) (mediaId string, err error) {
	mediaId, _, err = t.upload("video", data, title, introduction)
	return
}

func (t *Trader) AddThumbMaterial(data []byte) (mediaId, url string, err error) {
	return t.upload("thumb", data, "", "")
}

//新增永久图文素材
func (t *Trader) AddNews(nl NewsList) (mediaId string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := AddNewsURL + t.Accesstoken
	b, err := json.Marshal(nl)
	if err != nil {
		return
	}
	data, err := t.PostJson(surl, string(b))
	m := make(map[string]string)
	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}
	if _, ok := m["media_id"]; ok {
		mediaId = m["media_id"]
	} else {
		err = errors.New(string(data))
	}
	return
}

func (t *Trader) sendMsg(msg interface{}) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	surl := SendMsgURL + t.Accesstoken
	res, err := t.PostJson(surl, string(b))
	if err != nil {
		return
	}
	var v Res
	err = json.Unmarshal(res, &v)
	if err != nil {
		return
	}
	if v.ErrCode != 0 {
		err = errors.New(string(res))
	}
	return
}

func (t *Trader) SendTextMsg(touser, text string) error {
	textMsg := &TextMessage{
		Text: Text{Content: text},
	}
	textMsg.Touser, textMsg.MsgType = touser, textType
	return t.sendMsg(textMsg)
}

func (t *Trader) SendImageMsg(touser, imgMediaId string) error {
	imageMsg := &ImageMessage{
		Image: Image{MediaId: imgMediaId},
	}
	imageMsg.Touser, imageMsg.MsgType = touser, imageType
	return t.sendMsg(imageMsg)
}

func (t *Trader) SendVoiceMsg(touser, voiceMediaId string) error {
	v := &VoiceMessage{
		Voice: Voice{MediaId: voiceMediaId},
	}
	v.Touser, v.MsgType = touser, voiceType
	return t.sendMsg(v)
}

func (t *Trader) SendVideoMsg(touser, videoMediaId, thumbMediaId, title, description string) error {
	video := &VideoMessage{
		Video: Video{MediaId: videoMediaId, ThumbMediaId: thumbMediaId, Title: title, Description: description},
	}
	video.Touser, video.MsgType = touser, videoType
	return t.sendMsg(video)
}

func (t *Trader) SendMusicMsg(touser, title, description, musicurl, hqmusicurl, thumb_media_id string) error {
	m := &MusicMessage{
		Music: Music{
			Title:        title,
			Description:  description,
			MusicUrl:     musicurl,
			HqmusicUrl:   hqmusicurl,
			ThumbMediaId: thumb_media_id,
		},
	}
	m.Touser, m.MsgType = touser, musicType
	return t.sendMsg(m)
}

func (t *Trader) SendNewsMsg(touser string, articles []Article) error {
	a := &NewsMessage{
		News: News{
			Articles: articles,
		},
	}
	a.Touser, a.MsgType = touser, newsType
	return t.sendMsg(a)
}

func (t *Trader) SendMPNews(touser, mediaId string) error {
	n := &MpNewsMessage{
		MpNews: MpNews{MediaId: mediaId},
	}
	n.Touser, n.MsgType = touser, mpnewsType
	return t.sendMsg(n)
}

func (t *Trader) SendWxCardMsg(touser, cardId string) error {
	card := &WxCardMessage{
		WxCard: WxCard{CardId: cardId},
	}
	card.Touser, card.MsgType = touser, wxcardType
	return t.sendMsg(card)
}

func (t *Trader) SendMiniProgrampageMsg(touser, title, appid, pagePath, thumbMediaId string) error {
	mp := &MiniprogrampageMessage{
		Miniprogrampage: Miniprogrampage{
			Title:        title,
			AppId:        appid,
			PagePath:     pagePath,
			ThumbMediaId: thumbMediaId,
		},
	}
	mp.Touser, mp.MsgType = touser, miniprogrampageType
	return t.sendMsg(mp)
}

//开关客服输入状态
func (t *Trader) Typing(touser string, b bool) (err error) {
	var aa struct {
		Touser  string `json:"touser"`
		Command string `json:"command"`
	}
	aa.Touser = touser
	if b {
		aa.Command = "Typing"
	} else {
		aa.Command = "CancelTyping"
	}
	d, err := json.Marshal(aa)
	if err != nil {
		return
	}
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TypingURL + t.Accesstoken
	res, err := t.PostJson(surl, string(d))
	var m Res
	err = json.Unmarshal(res, &m)
	if err != nil {
		return
	}
	if m.ErrCode != 0 {
		err = errors.New(string(res))
	}
	return
}

func (t *Trader) kfAccount(action, kfaccount, nickname, password string) (err error) {
	var data struct {
		KfAccount string `json:"kf_account"`
		NickName  string `json:"nickname"`
		PassWord  string `json:"password"`
	}
	data.KfAccount, data.NickName, data.PassWord = kfaccount, nickname, password
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := KFaccountURL + action + "?access_token=" + t.Accesstoken
	res, err := t.PostJson(surl, string(b))
	if err != nil {
		return
	}
	var r Res
	err = json.Unmarshal(res, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(res))
	}
	return
}

//添加客服账号
func (t *Trader) AddKfAccount(kfaccount, nickname, password string) error {
	return t.kfAccount("add", kfaccount, nickname, password)
}

//修改客服账号
func (t *Trader) UpdateKfAccount(kfaccount, nickname, password string) error {
	return t.kfAccount("update", kfaccount, nickname, password)
}

//删除客服账号
func (t *Trader) DelKfAccount(kfaccount, nickname, password string) error {
	return t.kfAccount("del", kfaccount, nickname, password)
}

//设置客服账号头像
func (t *Trader) SetKfAccountheadImg(kfaccount string, imgdata []byte) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := SetKfAccountheadimgURL + "access_token=" + t.Accesstoken + "&kf_account=" + kfaccount
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("upload", "filename")
	if err != nil {
		return
	}
	_, err = io.Copy(fw, bytes.NewReader(imgdata))
	if err != nil {
		return
	}
	w.Close()
	resp, err := http.Post(surl, w.FormDataContentType(), buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
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

//获取所有客服账号
func (t *Trader) GetKfList() (list KfAccountList, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := GetkfListURL + t.Accesstoken
	resp, err := http.Get(surl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &list)
	if err != nil {
		return
	}
	if len(list.Kflist) == 0 {
		err = errors.New(string(b))
	}
	return
}

//菜单
func (t *Trader) bessMenu(action, menujson string) (b []byte, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := Menu + action + "?access_token=" + t.Accesstoken
	switch action {
	case "create":
		b, err = t.PostJson(surl, menujson)
	case "get":
		b, err = t.Get(surl)
	case "delete":
		b, err = t.Get(surl)
	case "addconditional":
		b, err = t.PostJson(surl, menujson)
	case "delconditional":
		b, err = t.PostJson(surl, menujson)
	case "trymatch":
		b, err = t.PostJson(surl, menujson)
	}

	return
}

//创建菜单 参数为菜单json字符串
func (t *Trader) CreateMenu(contentjson string) (err error) {
	b, err := t.bessMenu("create", contentjson)
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

//获取菜单 得到菜单json字符串
func (t *Trader) GetMenu() (menujson string, err error) {
	b, err := t.bessMenu("get", "")
	return string(b), err
}

//删除菜单
func (t *Trader) DelMenu() (err error) {
	b, err := t.bessMenu("delete", "")
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

//创建个性化菜单
func (t *Trader) AddConditionalMenu(menujson string) (menuid string, err error) {
	b, err := t.bessMenu("addconditional", menujson)
	if err != nil {
		return
	}
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}
	if _, ok := m["menuid"]; ok {
		menuid = m["menuid"]
	} else {
		err = errors.New(string(b))
	}
	return
}

//删除个性化菜单
func (t *Trader) DelconditionalMenu(menuid string) (err error) {
	var a struct {
		Menuid string `json:"menuid"`
	}
	a.Menuid = menuid
	str, err := json.Marshal(a)
	if err != nil {
		return
	}
	b, err := t.bessMenu("delconditional", string(str))
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

//测试个性化菜单匹配结果
func (t *Trader) trymatchMenu(userid string) (menujson string, err error) {
	var user struct {
		UserId string `json:"user_id"`
	}
	user.UserId = userid
	str, err := json.Marshal(user)
	if err != nil {
		return
	}
	b, err := t.bessMenu("trymatch", string(str))
	if err != nil {
		return
	}
	return string(b), err
}

//获取永久素材内容
func (t *Trader) GetMaterialInfo(mediaid string) (data []byte, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=" + t.Accesstoken
	var p struct {
		MediaId string `json:"media_id"`
	}
	p.MediaId = mediaid
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	if strings.Contains(string(str), "errcode") {
		err = errors.New(string(b))
		return
	}
	return b, err
}

//删除永久素材
func (t *Trader) DelMaterial(mediaid string) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/material/del_material?access_token=" + t.Accesstoken
	var p struct {
		MediaId string `json:"media_id"`
	}
	p.MediaId = mediaid
	str, err := json.Marshal(p)
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

//修改永久图文素材 index:要更新的文章在图文消息中的位置（多图文消息时，此字段才有意义），第一篇为0
func (t *Trader) UpdateNews(mediaid string, index int, article NewsArticle) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/material/update_news?access_token=" + t.Accesstoken
	var p struct {
		MediaId  string      `json:"media_id"`
		Index    int         `json:"index"`
		Articles NewsArticle `json:"articles"`
	}
	p.MediaId, p.Index, p.Articles = mediaid, index, article
	str, err := json.Marshal(p)
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

//获取素材总数
func (t *Trader) GetMaterialCount() (newscount, imagecount, videocount, voicecount int, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token=" + t.Accesstoken
	b, err := t.Get(surl)
	if err != nil {
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		return
	}
	var r struct {
		VoiceCount int `json:"voice_count"`
		VideoCount int `json:"video_count"`
		ImageCount int `json:"image_count"`
		NewsCount  int `json:"news_count"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	newscount, imagecount = r.NewsCount, r.ImageCount
	videocount, voicecount = r.VideoCount, r.VoiceCount
	return
}

//获取素材列表详情
/*
	type 素材的类型，图片（image）、视频（video）、语音 （voice）、图文（news）
	offset 从全部素材的该偏移位置开始返回，0表示从第一个素材 返回
	count 返回素材的数量，取值在1到20之间
	data 返回的json字符串 需要自己解析
*/
func (t *Trader) BatchGetMaterial(materialtype string, offset int, count int) (data string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=" + t.Accesstoken
	var p struct {
		Type   string `json:"type"`
		OffSet int    `json:"offset"`
		Count  int    `json:"count"`
	}
	p.Type, p.OffSet, p.Count = materialtype, offset, count
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		return
	}
	return string(b), err
}
