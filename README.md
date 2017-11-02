# wecaht - 一个小型的微信公众号api

你可以快速搭建一个订阅号/服务号应用，支持被动回复，主动回复，设置菜单，素材管理，用户管理等

## Installation



    $ go get github.com/slrem/wechat


## Examples 1

被动回复 可以和web框架配合使用
以echo为例

```Go
package main

import (
  "fmt"
  "log"
  "github.com/labstack/echo"
  "github.com/slrem/wechat"
)

  //错误处理
  func ErrorHandler(err error, c wechat.Context) error {
  		return c.Response().Success()
  }

//中间件
func LogUserActive() wechat.Middleware {
  	return func(next wechat.Handler) wechat.Handler {
  		return func(c wechat.Context) error {
        log.Println(c.Request().FromUserName())
  			return next(c)
  		}
  	}
  }

//文本消息处理
func textHandler(c wechat.Context) (err error) {
  	log.Println(c.Request().Content()) //用户发送来的文本消息

    //do somethings
    ...
    return c.Response().Text(c.Request().Content())
}

//关注事件处理
func subscribeHandler(c wechat.Context) (err error) {

	return c.Response().Text("欢迎关注[微笑]")
}

//取消关注事件处理
func unsubscribeHandler(c wechat.Context) (err error) {
	return c.Response().Success() //不做处理
}

//菜单处理事件
func clickMenuHandler(c wechat.Context) (err error) {
  key:=c.Request().EventKey()
	return  c.Response().Image(mediaid)
}

func main()  {
  w, err := wechat.New(
  		"appID", //公众号appid
  		"appsecret", //公众号appsecret
  		"token", //公众号设置的token
  		"encodingAESKey", //公众号加密钥匙
      nil
      )


  w.WechatErrorHandler = ErrorHandler

  w.Use(LogUserActive()) // 记录活跃时间

  w.Text(textHandler)
  w.SubscribeEvent(subscribeHandler)
  w.UnsubscribeEvent(unsubscribeHandler)
  w.MenuClickEvent(clickMenuHandler)

  e := echo.New()
  e.Any("/wechat/:app", func(c echo.Context) (err error) {
    w.Server(c.Response(), c.Request())
  }
}

e.Start(":8080")

```

## Examples 2


```Go
package main

import (
  "fmt"
  "log"

  "github.com/slrem/wechat"
)

func main() {
	w, err := wechat.New(
		"appID", //公众号appid
		"appsecret", //公众号appsecret
		"token", //公众号设置的token
		"encodingAESKey", //公众号加密钥匙
    nil
    )  

	t := w.Trader() //获取一个操作器
	//创建菜单
  menustr:=`{
     "button":[
     {
          "type":"click",
          "name":"今日歌曲",
          "key":"V1001_TODAY_MUSIC"
      },
      {
           "name":"菜单",
           "sub_button":[
           {
               "type":"view",
               "name":"搜索",
               "url":"http://www.soso.com/"
            },
            {
                 "type":"miniprogram",
                 "name":"wxa",
                 "url":"http://mp.weixin.qq.com",
                 "appid":"wx286b93c14bbf93aa",
                 "pagepath":"pages/lunar/index"
             },
            {
               "type":"click",
               "name":"赞一下我们",
               "key":"V1001_GOOD"
            }]
       }]
 }`
  t.CreateMenu(menustr)

//获取粉丝openid

	f, _ := t.GetFans("")
	for _,v:=range f{
    log.Println(v)
  }

//主动发送消息
  t.SendTextMsg("openid", "你好")

//添加一个图片素材
  b,_:=toutil.ReadFile("image.jpg")
  mediaid,url,err:=t.AddImageMaterial(b)

//主动发送消息
  t.SendImageMsg("openid", mediaid)

//群发消息 tagid为0 表示发给全部，其他的为发给属于标签id的所有用户
t.SendTextAll(tagid, "这是群发消息")

```
