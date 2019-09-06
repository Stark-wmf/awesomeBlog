package admin

import (
	common "awesomeblog/controllers"
	"awesomeblog/core"
	"awesomeblog/services"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ArticleGet(c *gin.Context) {
	articles, e := services.GetArticleList(0, services.ArticleAll)
	if e != nil {
		core.Error(e.Error())
	}
	datas, e := services.GetCategoryList(services.CategoryNORMAL)

	if e != nil {
		core.Error(e.Error())
	}

	c.JSON(200,gin.H{"articles": articles,"categorys":datas})
	//c.HTML(http.StatusOK, "admin/article/index.html", gin.H{"title": "文章管理", "highlight": "article", "data": articles})
}

func ArticleAddGet(c *gin.Context) {
	datas, e := services.GetCategoryList(services.CategoryNORMAL)

	if e != nil {
		core.Error(e.Error())
	}
     c.JSON(200,datas)
	//c.HTML(http.StatusOK, "admin/article/add.html", gin.H{"title": "文章管理", "highlight": "article", "categorys": datas})
}

func ArticleEditGet(c *gin.Context) {
	datas, e := services.GetCategoryList(services.CategoryNORMAL)

	if e != nil {
		core.Error(e.Error())
	}

	tid := c.Query("id")
	id, e := strconv.Atoi(tid)
	if e != nil || id <= 0 {
		c.Redirect(http.StatusFound, "/admin/article/add")
		return
	}

	article, e := services.GetArticleById(id)
	if e != nil {
		core.Error(e.Error())
		c.Redirect(http.StatusFound, "/admin/article/add")
		return
	}
	//article.Content = html.UnescapeString(article.Content)
	c.JSON(200,gin.H{"title": "文章管理", "highlight": "article", "categorys": datas, "article": article})
	//c.HTML(http.StatusOK, "admin/article/edit.html", gin.H{"title": "文章管理", "highlight": "article", "categorys": datas, "article": article})
	return

}

func ArticleAddPost(c *gin.Context) {
	categoryId, rs := c.GetPostForm("category")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1005, nil))
		return
	}
	category_id, e := strconv.Atoi(categoryId)
	if e != nil {
		c.JSON(http.StatusOK, common.GetMessage(1005, nil))
		return
	}

	if category_id == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1005, nil))
		return
	}

	title, rs := c.GetPostForm("title")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1009, nil))
		return
	}

	title = strings.Replace(title, " ", "", -1)
	if len(title) == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1009, nil))
		return
	}

	content, rs := c.GetPostForm("content")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1010, nil))
		return
	}

	// content = strings.Replace(content, " ", "", -1)
	if len(content) == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1010, nil))
		return
	}

	depiction, rs := c.GetPostForm("depiction")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1011, nil))
		return
	}

	depiction = strings.Replace(depiction, " ", "", -1)
	if len(depiction) == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1011, nil))
		return
	}

	//image, rs := c.GetPostForm("myimage")
	//if !rs {
	//	c.JSON(http.StatusOK, common.GetMessage(1012, nil))
	//	return
	//}
	myfile, e := c.FormFile("imgFile")
	var localname string
	if myfile!=nil {
		if e != nil {
			core.Error(e.Error())
			c.JSON(http.StatusOK, map[string]interface{}{"error": 1007, "msg": "读取文件失败"})
			return
		}
		h := md5.New()
		h.Write([]byte(myfile.Filename))
		cipherStr := h.Sum(nil)
		localname = fmt.Sprintf("%s/%d%s", "images", time.Now().Unix(), hex.EncodeToString(cipherStr))
		if e = c.SaveUploadedFile(myfile, localname); e != nil {
			core.Error(e.Error())
			c.JSON(http.StatusOK, map[string]interface{}{"error": 1008, "msg": "保存文件失败"})
			return
		}
	}
	//image = strings.Replace(image, " ", "", -1)
	//if len(image) == 0 {
	//	c.JSON(http.StatusOK, common.GetMessage(1012, nil))
	//	return
	//}

	ttop, rs := c.GetPostForm("top")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1013, nil))
		return
	}

	top, e := strconv.ParseInt(ttop, 10, 8)

	if e != nil {
		c.JSON(http.StatusOK, common.GetMessage(1013, nil))
		return
	}

	if top != 1 && top != 2 {
		c.JSON(http.StatusOK, common.GetMessage(1013, nil))
		return
	}
	fmt.Println(category_id,title, html.EscapeString(content), int8(top), depiction, localname)
	e = services.AddArticle(category_id, title, html.EscapeString(content), int8(top), depiction, localname)

	if e != nil {
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}

	c.JSON(http.StatusOK, common.GetMessage(0, nil))
	return
}

func ArticleEditPost(c *gin.Context) {
	id, rs := c.GetPostForm("articleId")
	articleId, e := strconv.Atoi(id)
	if e != nil {
		c.JSON(http.StatusOK, common.GetMessage(1016, nil))
		return
	}

	if articleId <= 0 {
		c.JSON(http.StatusOK, common.GetMessage(1016, nil))
		return
	}

	article, e := services.GetArticleById(articleId)
	if e != nil {
		if e == services.NoData {
			c.JSON(http.StatusOK, common.GetMessage(1017, nil))
			return
		}
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(9999, e.Error()))
		return
	}

	categoryId, rs := c.GetPostForm("category")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1005, nil))
		return
	}
	category_id, e := strconv.Atoi(categoryId)
	if e != nil {
		c.JSON(http.StatusOK, common.GetMessage(1005, nil))
		return
	}

	if category_id == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1005, nil))
		return
	}

	title, rs := c.GetPostForm("title")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1009, nil))
		return
	}

	title = strings.Replace(title, " ", "", -1)
	if len(title) == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1009, nil))
		return
	}

	content, rs := c.GetPostForm("content")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1010, nil))
		return
	}

	// content = strings.Replace(content, " ", "", -1)
	if len(content) == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1010, nil))
		return
	}

	depiction, rs := c.GetPostForm("depiction")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1011, nil))
		return
	}

	depiction = strings.Replace(depiction, " ", "", -1)
	if len(depiction) == 0 {
		c.JSON(http.StatusOK, common.GetMessage(1011, nil))
		return
	}

	//image, rs := c.GetPostForm("myimage")
	//if !rs {
	//	c.JSON(http.StatusOK, common.GetMessage(1012, nil))
	//	return
	//}
	//
	//image = strings.Replace(image, " ", "", -1)
	//if len(image) == 0 {
	//	c.JSON(http.StatusOK, common.GetMessage(1012, nil))
	//	return
	//}
	myfile, e := c.FormFile("imgFile")
	if e != nil {
		core.Error(e.Error())
		c.JSON(http.StatusOK, map[string]interface{}{"error": 1007, "msg": "读取文件失败"})
		return
	}
	h := md5.New()
	h.Write([]byte(myfile.Filename))
	cipherStr := h.Sum(nil)
	localname := fmt.Sprintf("%s/%d%s", "images", time.Now().Unix(), hex.EncodeToString(cipherStr))
	if e = c.SaveUploadedFile(myfile, localname); e != nil {
		core.Error(e.Error())
		c.JSON(http.StatusOK, map[string]interface{}{"error": 1008, "msg": "保存文件失败"})
		return
	}

	ttop, rs := c.GetPostForm("top")
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(1013, nil))
		return
	}

	top, e := strconv.ParseInt(ttop, 10, 8)

	if e != nil {
		c.JSON(http.StatusOK, common.GetMessage(1013, nil))
		return
	}

	if top != 1 && top != 2 {
		c.JSON(http.StatusOK, common.GetMessage(1013, nil))
		return
	}

	article.CategoryId = category_id
	article.Title = title
	article.Content = html.EscapeString(content)
	article.Top = int8(top)
	article.Depiction = depiction
	article.Image = localname
	e = services.EditArticle(article)
	if e != nil {
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(99999, e.Error()))
		return
	}

	c.JSON(http.StatusOK, common.GetMessage(0, nil))
	return
}

func ArticleDeletePost(c *gin.Context) {
	id := c.PostForm("id")

	article_id, e := strconv.Atoi(id)
	if e != nil {
		c.JSON(http.StatusOK, common.GetMessage(1016, nil))
		return
	}

	rs, e := services.DeleteArticle(article_id)
	if e != nil {
		core.Error(e.Error())
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}
	if !rs {
		c.JSON(http.StatusOK, common.GetMessage(999, nil))
		return
	}
	c.JSON(http.StatusOK, common.GetMessage(0, nil))
	return
}
