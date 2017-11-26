package trader

type AccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int64  `json:"expires_in"`
}
type VideoDesc struct {
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
}

type Res struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

//客服ZH账号
type (
	KF struct {
		KfAccount    string `json:"kf_account"`
		KfNick       string `json:"kf_nick"`
		KfId         string `json:"kf_id"`
		KfHeadImgUrl string `json:"kf_headimgurl"`
	}
	KfAccountList struct {
		Kflist []KF `json:"kf_list"`
	}
)

type (
	Text struct {
		Content string `json:"content"`
	}
	Image struct {
		MediaId string `json:"media_id"`
	}
	Voice struct {
		MediaId string `json:"media_id"`
	}
	MpVideo struct {
		MediaId string `json:"media_id"`
	}
	Video struct {
		MediaId      string `json:"media_id"`
		ThumbMediaId string `json:"thumb_media_id"`
		Title        string `json:"title"`
		Description  string `json:"description"`
	}
	Music struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		MusicUrl     string `json:"musicurl"`
		HqmusicUrl   string `json:"hqmusicurl"`
		ThumbMediaId string `json:"thumb_media_id"`
	}
	Article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		PicUrl      string `json:"picurl"`
	}
	News struct {
		Articles []Article `json:"articles"`
	}
	MpNews struct {
		MediaId string `json:"media_id"`
	}
	WxCard struct {
		CardId string `json:"card_id"`
	}
	Miniprogrampage struct {
		Title        string `json:"title"`
		AppId        string `json:"appid"`
		PagePath     string `json:"pagepath"`
		ThumbMediaId string `json:"thumb_media_id"`
	}
)

type (
	Customservice struct {
		KFaccount string `json:"kf_account"`
	}

	BassMessage struct {
		Touser  string `json:"touser"`
		MsgType string `json:"msgtype"`
	}
	TextMessage struct {
		BassMessage
		Text          Text          `json:"text"`
		Customservice Customservice `json:"customservice"`
	}
	ImageMessage struct {
		BassMessage
		Image         Image         `json:"image"`
		Customservice Customservice `json:"customservice"`
	}
	VoiceMessage struct {
		BassMessage
		Voice         Voice         `json:"voice"`
		Customservice Customservice `json:"customservice"`
	}
	VideoMessage struct {
		BassMessage
		Video         Video         `json:"video"`
		Customservice Customservice `json:"customservice"`
	}
	MusicMessage struct {
		BassMessage
		Music         Music         `json:"music"`
		Customservice Customservice `json:"customservice"`
	}
	NewsMessage struct {
		BassMessage
		News          News          `json:"news"`
		Customservice Customservice `json:"customservice"`
	}
	MpNewsMessage struct {
		BassMessage
		MpNews        MpNews        `json:"mpnews"`
		Customservice Customservice `json:"customservice"`
	}
	WxCardMessage struct {
		BassMessage
		WxCard        WxCard        `json:"wxcard"`
		Customservice Customservice `json:"customservice"`
	}
	MiniprogrampageMessage struct {
		BassMessage
		Miniprogrampage Miniprogrampage `json:"miniprogrampage"`
		Customservice   Customservice   `json:"customservice"`
	}
)

type (
	NewsArticle struct {
		Title              string `json:"title"`
		ThumbMediaId       string `json:"thumb_media_id"`
		Author             string `json:"author"`
		Digest             string `json:"digest"`
		ShowCoverPic       int    `json:"show_cover_pic"`
		Content            string `json:"content"`
		ContentSourceUrl   string `json:content_source_url`
		NeedOpenComment    int    `json:"need_open_comment"`
		OnlyFansCanComment int    `json:"only_fans_can_comment"`
	}
	NewsList struct {
		Articles []NewsArticle `json:"articles"`
	}
)

type SendAllResp struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	MsgId     int    `json:"msg_id"`
	MsgDataId int    `json:"206227730"`
}

//群发预览结构体
type (
	PreviewMpNews struct {
		ToUser  string `json:"touser"`
		MpNews  MpNews `json:"mpnews"`
		MsgType string `json:"msgtype"`
	}
	PreviewText struct {
		ToUser  string `json:"touser"`
		Text    Text   `json:"text"`
		MsgType string `json:"msgtype"`
	}
	PreviewVoice struct {
		ToUser  string `json:"touser"`
		Voice   Voice  `json:"voice"`
		MsgType string `json:"msgtype"`
	}
	PreviewImage struct {
		ToUser  string `json:"touser"`
		Image   Image  `json:"image"`
		MsgType string `json:"msgtype"`
	}
	PreviewMpVideo struct {
		ToUser  string  `json:"touser"`
		MpVideo MpVideo `json:"mpvideo"`
		MsgType string  `json:"msgtype"`
	}
)

//标签
type (
	Tag struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
)

//用户信息
type UserInfo struct {
	Subscribe     int    `json:"subscribe"`
	Openid        string `json:"openid"`
	NickName      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	HeadimgUrl    string `json:"headimgurl"`
	SubscribeTime int64  `json:"subscribe_time"`
	Unionid       string `json:"unionid"`
	Remark        string `json:"remark"`
	Groupid       int    `json:"groupid"`
	TagidList     []int  `json:"tagid_list"`
}

type Fans struct {
	Total      int    `json:"total"`
	Count      int    `json:"count"`
	NextOpenId string `json:"next_openid"`
	Data       struct {
		OpenId []string `json:"openid"`
	} `json:"data"`
}

//评论
type Comment struct {
	UserCommentID int    `json:"user_comment_id "`
	OpenID        string `json:"openid "`
	CreateTime    int    `json:"create_time "`
	Content       string `json:"content "`
	CommentType   int    `json:"comment_type "`
	Reply         struct {
		Content    string `json:"content "`
		CreateTime int    `json:"create_time "`
	} `json:"reply "`
}
