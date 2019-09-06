package blog

import (
	common "awesomeblog/controllers"
	"awesomeblog/core"
	"awesomeblog/services"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"strconv"
)

var commentPageNum = 5

func CommentGet(c *gin.Context) {
	c.JSON(200,gin.H{"title": core.GetConfig().WebsiteName, "highlight": "article"})
	//c.HTML(http.StatusOK, "blog/comment/index.html", gin.H{"title": core.GetConfig().WebsiteName, "highlight": "article"})
}

func CommentAddPost(c *gin.Context) {
	username,err:=core.Get("username")
	if err!=nil{
		c.JSON(http.StatusOK, common.GetMessage(1020, "redis错误了"))
		return
	}
	//username:=
	if username == "" {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	if !core.CheckEmail(username) {
		c.JSON(http.StatusOK, common.GetMessage(1000, nil))
		return
	}

	tmp_id, rs := c.GetPostForm("article_id")
	if !rs {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "非法请求"})
		return
	}

	article_id, e := strconv.Atoi(tmp_id)
	if e != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "非法请求"})
		return
	}

	if article_id <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "非法请求"})
		return
	}

	content, rs := c.GetPostForm("content")
	if !rs {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "输入正确的评论内容"})
		return
	}

	//email, rs := c.GetPostForm("email")
	//if !rs {
	//	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "输入正确的email"})
	//	return
	//}
	//
	//if !core.CheckEmail(email) {
	//	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "输入正确的email"})
	//	return
	//}
	//
	//name, rs := c.GetPostForm("name")
	//if !rs {
	//	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "输入正确的name"})
	//	return
	//}


	e = services.AddComment(article_id, username, html.EscapeString(content), c.Request.RemoteAddr)
	if e != nil {
		core.Error(e.Error())
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "操作成功"})
	return
}

func CommentGetListForPage(c *gin.Context) {
	tmp := c.Query("pageIndex")
	pageIndex, e := strconv.Atoi(tmp)
	if e != nil {
		pageIndex = 1
	}

	tmp = c.Query("article_id")
	article_id, e := strconv.Atoi(tmp)
	if e != nil {
		article_id = 0
	}

	comments, e := services.GetCommentListForPage(pageIndex, commentPageNum, article_id)
	if e != nil {
		core.Error(e.Error())
	}
	c.JSON(http.StatusOK, common.GetMessage(0, comments))
	return
}
