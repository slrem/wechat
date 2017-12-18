package trader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Jsapi_ticket struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Ticket     string `json:"ticket"`
	Expires_in int64  `json:"Expires_in"`
}

func (t *Trader) isJTAlive() bool {
	if t.JsapiTicket == "" || t.JsapiTicketExpiresIn-time.Now().Unix() < 120 {
		return false
	}
	return true
}

func (t *Trader) GetJsapiTicket() (ticket string, err error) {
	t.mtx.Lock()
	if !t.isJTAlive() {
		err = t.httpGetJsapi_ticket()
	}
	ticket = t.JsapiTicket
	t.mtx.Unlock()
	return
}
func (t *Trader) SetJsapiTicket(ticket string, Expires_in int64) {
	t.JsapiTicket, t.JsapiTicketExpiresIn = ticket, Expires_in
}

func (t *Trader) httpGetJsapi_ticket() (err error) {
	res, err := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + t.Accesstoken + "&type=jsapi")
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var jt Jsapi_ticket
	err = json.Unmarshal(b, &jt)
	if err != nil {
		return
	}
	if jt.ErrCode != 0 {
		err = errors.New(string(b))
		return
	}
	t.SetJsapiTicket(jt.Ticket, jt.Expires_in+time.Now().Unix())
	return
}

type WebConf struct {
	AppId     string
	Timestamp int64
	Noncestr  string
	Signature string
}

func (t *Trader) WebConfig(url string) (wf WebConf, err error) {
	wf.AppId = t.AppId
	ticket, err := t.GetJsapiTicket()
	if err != nil {
		return
	}
	wf.Timestamp = time.Now().Unix()
	wf.Noncestr = GetRandStr(16)
	str := "jsapi_ticket=" + ticket + "&noncestr=" +
		wf.Noncestr + "&timestamp=" + fmt.Sprint(wf.Timestamp) + "&url=" + url
	wf.Signature = sha(str)
	return
}
