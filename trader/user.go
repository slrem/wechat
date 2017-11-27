package trader

import (
	"encoding/json"
	"errors"
	"strings"
)

/*
  用户管理
*/
//创建标签
func (t *Trader) CreateTag(tagname string) (tagid int, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TagsURL + "create?access_token=" + t.Accesstoken
	var p struct {
		Tag struct {
			Name string `json:"name"`
		} `json:"tag"`
	}
	p.Tag.Name = tagname
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	var r struct {
		Tag struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.Tag.Id == 0 {
		err = errors.New(string(b))
	} else {
		tagid = r.Tag.Id
	}
	return
}

//获取公众号已创建的标签
func (t *Trader) GetTag() (tags []Tag, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TagsURL + "get?access_token=" + t.Accesstoken
	b, err := t.Get(surl)
	if err != nil {
		return
	}
	var r struct {
		Tags []Tag `json:"tags"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if len(r.Tags) == 0 {
		err = errors.New(string(b))
	} else {
		tags = r.Tags
	}
	return
}

//编辑标签
func (t *Trader) UpdateTag(tagid int, tagname string) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TagsURL + "update?access_token=" + t.Accesstoken
	var p struct {
		Tag struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
	}
	p.Tag.Id, p.Tag.Name = tagid, tagname
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

//删除标签
func (t *Trader) DelTag(tagid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := TagsURL + "delete?access_token=" + t.Accesstoken
	var p struct {
		Tag struct {
			Id int `json:"id"`
		} `json:"tag"`
	}
	p.Tag.Id = tagid
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

//获取标签下粉丝列表
/*
tagid 标签id, nextopenid 第一个拉取的OPENID，不填默认从头开始拉取
*/
func (t *Trader) GetUserByTag(tagid int, nextopenid string) (useropenid []string, lastopenid string, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token=" + t.Accesstoken
	var p struct {
		TagId      int    `json:"tagid"`
		NextOpenId string `json:"next_openid"`
	}
	p.TagId, p.NextOpenId = tagid, nextopenid
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
	var r struct {
		Count int `json:"count"`
		Data  struct {
			OpenId []string `json:"openid"`
		} `json:"data"`
		NextOpenId string `json:"next_openid"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	useropenid, lastopenid = r.Data.OpenId, r.NextOpenId
	return
}

//批量为用户打标签
func (t *Trader) BatchTagToUsers(useropenids []string, tagid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=" + t.Accesstoken
	var p struct {
		OpenIds []string `json:"openid_list"`
		TagId   int      `json:"tagid"`
	}
	p.OpenIds, p.TagId = useropenids, tagid
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

//批量为用户取消标签
func (t *Trader) BatchCancelTag(useropenid []string, tagid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=" + t.Accesstoken
	var p struct {
		OpenIds []string `json:"openid_list"`
		TagId   int      `json:"tagid"`
	}
	p.OpenIds, p.TagId = useropenid, tagid
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

//获取用户身上的标签列表
func (t *Trader) GetTagsByUser(useropenid string) (tagids []int, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=" + t.Accesstoken
	var p struct {
		OpenId string `json:"openid"`
	}
	p.OpenId = useropenid
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
	var r struct {
		TagIds []int `json:"tagid_list"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	tagids = r.TagIds
	return
}

//设置用户备注名
func (t *Trader) SetRemark(useropenid string, remark string) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=" + t.Accesstoken
	var p struct {
		OpenId string `json:"openid"`
		Remark string `json:"remark"`
	}
	p.OpenId, p.Remark = useropenid, remark
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

//获取用户基本信息（包括UnionID机制）
func (t *Trader) GetUserInfo(openid string) (user UserInfo, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + t.Accesstoken + "&openid=" + openid + "&lang=zh_CN "
	b, err := t.Get(surl)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return
	}
	if user.Openid == "" {
		err = errors.New(string(b))
	}
	return
}

//获取用户列表
/*
nextopenid 第一个拉取的OPENID，填空默认从头开始拉取
附：关注者数量超过10000时
当公众号关注者数量超过10000时，可通过填写next_openid的值，从而多次拉取列表的方式来满足需求。
*/
func (t *Trader) GetFans(nextopenid string) (fans Fans, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + t.Accesstoken + "&next_openid=" + nextopenid
	b, err := t.Get(surl)
	if err != nil {
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		return
	}

	err = json.Unmarshal(b, &fans)
	return

}
