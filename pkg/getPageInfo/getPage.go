package getPageInfo

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetPage 获取分页参数
func GetPage(c *gin.Context) (page, size int64) {
	var err error
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
