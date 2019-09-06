package admin

import (
	common "awesomeblog/controllers"
	"awesomeblog/core"
	"awesomeblog/services"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)
func LoginGet(c *gin.Context) {
//	c.HTML(http.StatusOK, "admin/login/index.html", gin.H{"title": "登录"})
}
func Vertify(c *gin.Context) {
	usercode:=c.PostForm("code")
	code,err:=core.Get("code")
	fmt.Println("code"+code)
	if err!=nil{
		c.JSON(http.StatusOK, common.GetMessage(1020, "redis错误了"))
		return
	}
	username,err:=core.Get("username")
	if err!=nil{
		c.JSON(http.StatusOK, common.GetMessage(1020, "redis错误了"))
		return
	}
	if code!=usercode{
		c.JSON(http.StatusOK, common.GetMessage(1021, "验证码错误了"))
		return
	}
	admin,e:=services.QueryDataByUserName(username)
	if e != nil {
		if e != services.NoData {
			c.JSON(http.StatusOK, common.GetMessage(2, "用户尚未注册"))
			return
		}
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(1003, nil))
		return
	}
	services.EditUser(admin)
	c.JSON(http.StatusOK, common.GetMessage(0, gin.H{"username":username,"groupid":admin.GroupId,"nickname":admin.NickName}))
	//	c.HTML(http.StatusOK, "admin/login/index.html", gin.H{"title": "登录"})
}

func LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	if !core.CheckEmail(username) {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	if password == "" {
		c.JSON(http.StatusOK, common.GetMessage(1001, nil))
		return
	}

	admin, e := services.QueryDataByUserName(username)
	if e != nil {
		if e != services.NoData {
			c.JSON(http.StatusOK, common.GetMessage(9999, nil))
			return
		}
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(1003, nil))
		return
	}

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", username, password)))
	t := h.Sum(nil)

	if admin.PassWord != hex.EncodeToString(t) {
		c.JSON(http.StatusOK, common.GetMessage(1002, nil))
		return
	}

	sessionId := core.GenSessionID()
	cookie := &http.Cookie{Name: "session_id", Value: sessionId, Path: "/", HttpOnly: true, Expires: time.Now().AddDate(0, 0, 7)}
	http.SetCookie(c.Writer, cookie)

	jadmin, _ := json.Marshal(admin)

	if _, e := core.Set(sessionId, string(jadmin), time.Hour*24*7); e != nil {
		core.Errorf("set user info to redis ", e.Error())
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}

	c.JSON(http.StatusOK, common.GetMessage(0, gin.H{"username":username,"groupid":1}))
	return
}

func LogoutGet(c *gin.Context) {
	var cookie *http.Cookie
	var err error
	if cookie, err = c.Request.Cookie("session_id"); err != nil {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	_, e := core.Delete(cookie.Value)
	if e != nil {
		core.Error(e.Error())
	}
	_, ee := core.Delete("username")
	if ee != nil {
		core.Error(ee.Error())
	}
	c.Redirect(http.StatusFound, "/admin/login")

}

func RegistePost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")


	_, e := services.QueryDataByUserName(username)

	if e != services.NoData {
			c.JSON(http.StatusOK, common.GetMessage(1, "用户已经注册过了"))
			return
		}
	code:=services.Sendemail(username)

    sessionId := core.GenSessionID()
	cookie := &http.Cookie{Name: "code", Value: sessionId, Path: "/", HttpOnly: true, Expires: time.Now().AddDate(0, 0, 7)}
	http.SetCookie(c.Writer, cookie)


	if _, e := core.Set("code", code, time.Hour*24*7); e != nil {
		core.Errorf("set code info to redis ", e.Error())
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}

	if _, e := core.Set("username", username, time.Hour*24*7); e != nil {
		core.Errorf("set username info to redis ", e.Error())
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}

	if username == "" {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	if !core.CheckEmail(username) {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	if password == "" {
		c.JSON(http.StatusOK, common.GetMessage(1001, nil))
		return
	}

	if nickname == "" {
		c.JSON(http.StatusOK, common.GetMessage(1004, nil))
		return
	}


	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", username, password)))
	t := h.Sum(nil)

	ppw:= hex.EncodeToString(t)



	er:=services.RegisteUser(username,ppw,nickname)
	if er!=nil{
		c.JSON(http.StatusOK, common.GetMessage(1021, "更新错误"))
		return
	}


	//jadmin, _ := json.Marshal(admin)
	//
	//if _, e := core.Set(sessionId, string(jadmin), time.Hour*24*7); e != nil {
	//	core.Errorf("set user info to redis ", e.Error())
	//	c.JSON(http.StatusOK, common.GetMessage(999, nil))
	//	return
	//}

	c.JSON(http.StatusOK, common.GetMessage(0, nil))
	return
}

func UserLoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
   fmt.Println(username+"--")
	if username == "" {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	if !core.CheckEmail(username) {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}
	fmt.Println(password+"--")
	if password == "" {
		c.JSON(http.StatusOK, common.GetMessage(1001, nil))
		return
	}

	admin, e := services.QueryDataByUserName(username)
	if e != nil {
		if e == services.NoData {
			c.JSON(http.StatusOK, common.GetMessage(2, "用户尚未注册"))
			return
		}
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(1003, nil))
		return
	}

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", username, password)))
	t := h.Sum(nil)

	if admin.PassWord != hex.EncodeToString(t) {
		c.JSON(http.StatusOK, common.GetMessage(3, "密码错误"))
		return
	}

	sessionId := core.GenSessionID()
	cookie := &http.Cookie{Name: "session_id", Value: sessionId, Path: "/", HttpOnly: true, Expires: time.Now().AddDate(0, 0, 7)}
	http.SetCookie(c.Writer, cookie)
	c.Request.Header.Set("username",username)

	jadmin, _ := json.Marshal(admin)
	fmt.Println(admin)

	if _, e := core.Set(sessionId, string(jadmin), time.Hour*24*7); e != nil {
		core.Errorf("set user info to redis ", e.Error())
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}

	c.JSON(http.StatusOK, common.GetMessage(1, gin.H{"username":username,"groupid":admin.GroupId,"nickname":admin.NickName}))
	return
}

