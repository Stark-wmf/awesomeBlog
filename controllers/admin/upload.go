package admin

import (
	"awesomeblog/core"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Imageupload(c *gin.Context) {
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

	c.JSON(http.StatusOK, map[string]interface{}{"error": 0, "url": "/" + localname})
	return

}
