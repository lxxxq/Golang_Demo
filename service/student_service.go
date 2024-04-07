package service

import (
	"demo/entity"
	"demo/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetByID(c *gin.Context) {
	userID := c.Param("id")
	num, err := strconv.ParseUint(userID, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式无效"})
		return
	}

	result := repository.SelectByID(num)
	c.JSON(http.StatusOK, result)
}

func GetAll(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	name := c.Query("name")

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page无效"})
		return
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pageSize无效"})
		return
	}

	entities, count, err := repository.SelectAll(pageNum, pageSizeNum, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entities, "total": count})
}

func DropById(c *gin.Context) {
	userID := c.Param("id")

	num, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式无效"})
		return
	}

	err = repository.DeleteByID(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "删除失败"})
	}

	c.JSON(http.StatusOK, "success")
}

func AlterByIds(c *gin.Context) {
	var students []entity.Student

	err := c.BindJSON(&students)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json格式错误"})
		return
	}

	err = repository.UpdateByID(students)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "修改失败"})
		return
	}

	c.JSON(http.StatusOK, "success")
}

func BatchInsert(c *gin.Context) {
	var students []entity.Student

	err := c.BindJSON(&students)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json格式错误"})
		return
	}

	err = repository.InsetData(students)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新增失败"})
		return
	}

	c.JSON(http.StatusOK, "success")
}
