package blog

import (
	common "awesomeblog/controllers"
	"awesomeblog/core"
	"awesomeblog/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"math"
	"net/http"
	"strconv"
)

var articlesPageNum int = 5

func ArticleGet(c *gin.Context) {
	tmp, e := services.GetArticleTotal()
	if e != nil {
		core.Error(e.Error())
	}

	total := math.Ceil(float64(tmp) / float64(articlesPageNum))
	c.JSON(200,total)
//	c.HTML(http.StatusOK, "blog/article/index.html", gin.H{"title": core.GetConfig().WebsiteName, "highlight": "article", "total": total})
}

func ArticleGetListForPage(c *gin.Context) {
	tmp := c.Query("pageIndex")
	pageIndex, e := strconv.Atoi(tmp)
	if e != nil {
		pageIndex = 1
	}

	tmp = c.Query("category_id")
	category_id, e := strconv.Atoi(tmp)
	if e != nil {
		category_id = 0
	}

	articles, e := services.GetArticleListForPage(pageIndex, articlesPageNum, category_id)
	if e != nil {
		core.Error(e.Error())
	}

	ids := make([]int, 0)
	for i := 0; i < len(articles); i++ {
		ids = append(ids, articles[i].ID)
	}

	nums, e := services.GetNumByArticleIds(ids)
	if e != nil {
		core.Error(e.Error())
	}

	type T struct {
		ReadNum    int
		CommentNum int
	}

	tmp1 := make(map[int]T)
	for i := 0; i < len(nums); i++ {
		tmp1[nums[i].ArticleId] = T{ReadNum: nums[i].ReadNum, CommentNum: nums[i].CommentNum}
	}

	for i := 0; i < len(articles); i++ {
		articles[i].ReadNum = tmp1[articles[i].ID].ReadNum
		articles[i].CommentNum = tmp1[articles[i].ID].CommentNum
	}
	c.JSON(http.StatusOK, common.GetMessage(0, articles))
	return
}

func ArticleDetailGet(c *gin.Context) {
	tmp_id := c.Query("id")
	article_id, e := strconv.Atoi(tmp_id)
	if e != nil {
		core.Error(e.Error())
		//c.Redirect(http.StatusFound, "/home/index")
	}

	comments, e := services.GetCommentListForPage(1, commentPageNum, article_id)
	if e != nil {
		core.Error(e.Error())
	}

	article, e := services.GetArticleById(article_id)
	if e != nil {
		core.Error(e.Error())
		//c.Redirect(http.StatusFound, "/home/index")
	}
	article.Content = html.UnescapeString(article.Content)
	num, e := services.GetNumByArticleId(article.ID)
	if e != nil {
		core.Error(e.Error())
	}
	article.ReadNum = num.ReadNum
	article.CommentNum = num.CommentNum

	//articles, e := services.GetArticleList(10, services.ArticleNORMAL)
	if e != nil {
		core.Error(e.Error())
	}

	tmp, e := services.GetCommentTotalByArticle(article_id)
	if e != nil {
		core.Error(e.Error())
	}

	total := math.Ceil(float64(tmp) / float64(commentPageNum))
//	twos, e := services.GetArticleByTwo(article_id)
	fmt.Println(total)
	if e != nil {
		core.Error(e.Error())
	}

	if e := services.UpdateReadNum(article_id); e != nil {
		core.Error(e.Error())
	}
    c.JSON(200,gin.H{"highlight": "home", "article": article, "twos": "weos", "total": total,"comment":comments})
//	c.HTML(http.StatusOK, "blog/article/detail.html", gin.H{"title": core.GetConfig().WebsiteName, "highlight": "home", "article": article, "articles": articles, "twos": twos, "total": total})
}
