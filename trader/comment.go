package trader

import (
	"encoding/json"
	"errors"
)

/*
  图文评论接口
*/

//打开/关闭已群发文章评论 msgdataid 由SendMpNewsAll返回的字段
// bl ture为开启 fasle为关闭
func (t *Trader) OpenComment(bl bool, msgdataid int, index int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	var action string
	if bl {
		action = "open"
	} else {
		action = "close"
	}
	surl := CommentURL + action + "?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId int `json:"msg_data_id"`
		Index     int `json:"index"`
	}
	p.MsgDataId, p.Index = msgdataid, index
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

//查看指定文章的评论数据
/*
参数	           说明
msg_data_id	  群发返回的msg_data_id
index		      多图文时，用来指定第几篇图文，从0开始，不带默认返回该msg_data_id的第一篇图文
begin		      起始位置
count		     获取数目（>=50会被拒绝）
type		type=0 普通评论&精选评论type=1 普通评论 type=2 精选评论
*/
func (t *Trader) GetCommentList(msgdataid, index, begin, count, commenttype int) (list []Comment, err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := CommentURL + "list?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId int `json:"msg_data_id"`
		Index     int `json:"index"`
		Begin     int `json:"begin"`
		Count     int `json:"count"`
		Type      int `json:"type"`
	}
	p.MsgDataId, p.Index = msgdataid, index
	p.Begin, p.Count, p.Type = begin, count, commenttype
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := t.PostJson(surl, string(str))
	if err != nil {
		return
	}
	var r struct {
		ErrCode int       `json:"errcode"`
		ErrMsg  string    `json:"errmsg"`
		Total   int       `json:"total"`
		Comment []Comment `json:"comment"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
		return
	}
	list = r.Comment
	return
}

//将评论标记精选
func (t *Trader) MarkelectComment(msgdataid, index, usercommentid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := CommentURL + "markelect?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId     int `json:"msg_data_id"`
		Index         int `json:"index"`
		UserCommentID int `json:"user_comment_id "`
	}
	p.MsgDataId, p.Index, p.UserCommentID = msgdataid, index, usercommentid
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

//将评论取消精选
func (t *Trader) UnMarkelectComment(msgdataid, index, usercommentid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := CommentURL + "unmarkelect?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId     int `json:"msg_data_id"`
		Index         int `json:"index"`
		UserCommentID int `json:"user_comment_id "`
	}
	p.MsgDataId, p.Index, p.UserCommentID = msgdataid, index, usercommentid
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

//删除评论
func (t *Trader) DeleteComment(msgdataid, index, usercommentid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := CommentURL + "delete?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId     int `json:"msg_data_id"`
		Index         int `json:"index"`
		UserCommentID int `json:"user_comment_id "`
	}
	p.MsgDataId, p.Index, p.UserCommentID = msgdataid, index, usercommentid
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

//回复评论
func (t *Trader) ReplyComment(msgdataid, index, usercommentid int, content string) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := CommentURL + "reply/add?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId     int    `json:"msg_data_id"`
		Index         int    `json:"index"`
		UserCommentID int    `json:"user_comment_id "`
		Content       string `json:"content"`
	}
	p.MsgDataId, p.Index, p.UserCommentID = msgdataid, index, usercommentid
	p.Content = content
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

//删除回复
func (t *Trader) DeleteReplyComment(msgdataid, index, usercommentid int) (err error) {
	err = t.CheckAccessTokenLive()
	if err != nil {
		return
	}
	surl := CommentURL + "reply/delete?access_token=" + t.Accesstoken
	var p struct {
		MsgDataId     int `json:"msg_data_id"`
		Index         int `json:"index"`
		UserCommentID int `json:"user_comment_id "`
	}
	p.MsgDataId, p.Index, p.UserCommentID = msgdataid, index, usercommentid
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
