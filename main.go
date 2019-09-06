package main

import (
	"awesomeblog/controllers/admin"
	"awesomeblog/controllers/blog"
	"awesomeblog/core"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"time"
)

func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

func HtmlHandler(content string) template.HTML {
	return template.HTML(content)
}

func main() {
//	pwd:=services.Sendemail("2621233431@qq.com")
//	fmt.Println(pwd)
	logpath := flag.String("L", "conf/log.xml", "log config file path")
	err := core.InitLog(*logpath)
	if err != nil {
		fmt.Println("init log error ", err.Error())
		os.Exit(-1)
	}
	defer core.FlushLog()

	var fileName string
	flag.StringVar(&fileName, "S", "conf/conf.yaml", "system config file path")
	err = core.InitConfig(fileName)
	if err != nil {
		fmt.Println("init config ", err.Error())
		os.Exit(-1)
	}

	err = core.InitRedis()
	if err != nil {
		fmt.Println("init redis ", err.Error())
		os.Exit(-1)
	}

	defer core.CloseRedis()

	router := gin.Default()
	router.Use(cors.Default())
	router.SetFuncMap(template.FuncMap{
		"dateFormat":  DateFormat,
		"htmlHandler": HtmlHandler,
	})

	router.Static("/static", "./static")
	router.Static("/images", "./images")
	//router.GET("comment/index", blog.CommentGet)
	router.POST("/register", admin.RegistePost)
	router.POST("/vertify",admin.Vertify)
	router.POST("/login", admin.UserLoginPost)
	router.GET("/logout", admin.LogoutGet)
	router.GET("article/index", admin.ArticleGet)
	router.GET("article/detail", blog.ArticleDetailGet)
	router.GET("article/getList", blog.ArticleGetListForPage)
	router.POST("comment/add", blog.CommentAddPost)
	router.GET("comment/getList", blog.CommentGetListForPage)

	admins := router.Group("/admin")
	//admins.GET("/login", admin.LoginGet)
	//admins.POST("/login", admin.LoginPost)
	//admins.GET("/logout", admin.LogoutGet)

	admins.Use(AuthMiddleWare())
	{
		// admins.GET("/index", admin.AdminGet)
		admins.GET("/article/index", admin.ArticleGet)
		admins.GET("/article/add", admin.ArticleAddGet)
		admins.GET("/article/edit", admin.ArticleEditGet)
	//	admins.POST("/fileupload", admin.Imageupload)
		admins.POST("/article/add", admin.ArticleAddPost)
		admins.POST("/article/delete", admin.ArticleDeletePost)
		admins.POST("/article/edit", admin.ArticleEditPost)


	}

	err = core.InitDB()
	if err != nil {
		fmt.Println("init db error ", err.Error())
		os.Exit(-1)
	}
	defer core.Close()

	addr := core.GetConfig().Addr
	if addr == "" {
		addr = "8080"
	}

	fmt.Println("myblog run addr :", addr)
	if e := router.Run(addr); e != nil {
		panic(e)
	}

}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("session_id"); err == nil {
			value := cookie.Value
			_, err := core.Get(value)

			if err != nil {
				c.Redirect(http.StatusFound, "/admin/login")
				return
			}
			var rs bool
			if rs, err = core.Expire(value, time.Hour*24*7); err != nil {
				c.Redirect(http.StatusFound, "/admin/login")
				return
			}

			if !rs {
				c.Redirect(http.StatusFound, "/admin/login")
				return
			}

			c.Next()
			return

		}
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
}
