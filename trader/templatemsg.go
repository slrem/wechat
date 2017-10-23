package trader

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
发送模板消息
注:只能是认证的服务号才有此接口权限
*/

//设置所属行业
func (t *Trader) SetIndustry(industryId1, industryId2 int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TemplateURL + "api_set_industry?access_token=" + t.Accesstoken
	var p struct {
		IndustryId1 string `json:"industry_id1"`
		IndustryId2 string `json:"industry_id2"`
	}
	p.IndustryId1 = fmt.Sprint(industryId1)
	p.IndustryId2 = fmt.Sprint(industryId2)
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

//获取设置的行业信息
func (t *Trader) GetIndustry() (jsonstr string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TemplateURL + "get_industry?access_token=" + t.Accesstoken
	b, err := t.Get(surl)
	return string(b), err
}

//获得模板ID
func (t *Trader) GetTemplateId(templateIdShort string) (templateId string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TemplateURL + "api_add_template?access_token=" + t.Accesstoken
	var p struct {
		TemplateIdShort string `json:"template_id_short"`
	}
	p.TemplateIdShort = templateIdShort
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	var r struct {
		ErrCode    int    `json:errcode`
		ErrMsg     string `json:"errmsg"`
		TemplateId string `json:"template_id"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	} else {
		templateId = r.TemplateId
	}
	return
}

//获取模板列表 return josn字符串
func (t *Trader) GetALLTemplate() (jsonstr string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TemplateURL + "get_all_private_template?access_token=" + t.Accesstoken
	b, err := t.Get(surl)
	return string(b), err
}

//删除模板
func (t *Trader) DelTemplate(templateId string) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TemplateURL + "del_private_template?access_token=" + t.Accesstoken
	var p struct {
		TemplateId string `json:"template_id"`
	}
	p.TemplateId = templateId
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

//发送模板消息
func (t *Trader) SendTemplateMsg(jsonContext string) (msgid int, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + t.Accesstoken
	b, err := t.PostJson(surl, jsonContext)
	if err != nil {
		return
	}
	var r struct {
		ErrCode int    `json:errcode`
		ErrMsg  string `json:"errmsg"`
		MsgId   int    `json:"msgid"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	} else {
		msgid = r.MsgId
	}
	return
}
