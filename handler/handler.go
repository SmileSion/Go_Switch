package handler

import (
	"net/http"
	"edulimitrate/model"
	"github.com/gin-gonic/gin"
)

type RegionRequest struct {
	RegionCode string `json:"regioncode"`
}

func OpenRegion(c *gin.Context) {
	var req RegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}
	if err := model.InsertRegionCode(req.RegionCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "insert failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "region opened"})
}

func CloseRegion(c *gin.Context) {
	var req RegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}
	if err := model.DeleteRegionCode(req.RegionCode); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "region closed"})
}

func CheckRegion(c *gin.Context) {
	var req RegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}
	exist, err := model.ExistsRegionCode(req.RegionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "check failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "allowed": exist})
}
