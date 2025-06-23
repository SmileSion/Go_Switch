package handler

import (
	"edulimitrate/middleware"
	"edulimitrate/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegionRequest struct {
	RegionCode string `json:"regioncode"`
}

func OpenRegion(c *gin.Context) {
	var req RegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Logger.Printf("[OpenRegion] invalid json: %v, clientIP=%s", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}

	middleware.Logger.Printf("[OpenRegion] request regioncode=%s, clientIP=%s", req.RegionCode, c.ClientIP())

	if err := model.InsertRegionCode(req.RegionCode); err != nil {
		middleware.Logger.Printf("[OpenRegion] insert failed: %v, regioncode=%s", err, req.RegionCode)

		if err.Error() == "region code already exists" {
			c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "insert failed"})
		}
		return
	}

	middleware.Logger.Printf("[OpenRegion] success, regioncode=%s", req.RegionCode)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "region opened"})
}

func CloseRegion(c *gin.Context) {
	var req RegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Logger.Printf("[CloseRegion] invalid json: %v, clientIP=%s", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}

	middleware.Logger.Printf("[CloseRegion] request regioncode=%s, clientIP=%s", req.RegionCode, c.ClientIP())

	if err := model.DeleteRegionCode(req.RegionCode); err != nil {
		middleware.Logger.Printf("[CloseRegion] delete failed: %v, regioncode=%s", err, req.RegionCode)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	middleware.Logger.Printf("[CloseRegion] success, regioncode=%s", req.RegionCode)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "region closed"})
}

func CheckRegion(c *gin.Context) {
	var req RegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Logger.Printf("[CheckRegion] invalid json: %v, clientIP=%s", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}

	middleware.Logger.Printf("[CheckRegion] request regioncode=%s, clientIP=%s", req.RegionCode, c.ClientIP())

	exist, err := model.ExistsRegionCode(req.RegionCode)
	if err != nil {
		middleware.Logger.Printf("[CheckRegion] check failed: %v, regioncode=%s", err, req.RegionCode)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "check failed"})
		return
	}

	middleware.Logger.Printf("[CheckRegion] result: regioncode=%s, allowed=%v", req.RegionCode, exist)
	c.JSON(http.StatusOK, gin.H{"code": 200, "allowed": exist})
}
