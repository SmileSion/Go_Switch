package handler

import (
	"edulimitrate/config"
	"edulimitrate/middleware"
	"edulimitrate/model"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var DecryptedSecretKey string

type RegionRequest struct {
	RegionCode string `json:"regioncode"`
}
type PromptRequest struct {
	Prompt string `json:"prompt"`
}


func InitDecryptedSecret() {
	key, err := middleware.Decrypt(config.Conf.Key.SecretKey)
	if err != nil {
		middleware.Logger.Fatalf("Failed to decrypt secretKey: %v", err)
	}
	DecryptedSecretKey = key
}

func OpenRegion(c *gin.Context) {
	if secret := c.GetHeader("secretKey"); secret != DecryptedSecretKey {
		middleware.Logger.Printf("[OpenRegion] unauthorized access, clientIP=%s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}
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
	if secret := c.GetHeader("secretKey"); secret != DecryptedSecretKey {
		middleware.Logger.Printf("[OpenRegion] unauthorized access, clientIP=%s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}
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

	prompt, err := model.GetPromptWords()
	if err != nil {
		prompt = "" // 出错就返回空，不影响主流程
	}

	middleware.Logger.Printf("[CheckRegion] result: regioncode=%s, allowed=%v，promptwors=%s;", req.RegionCode, exist,prompt)
	c.JSON(http.StatusOK, gin.H{"code": 200, "allowed": exist, "promptwords": prompt})
}

func UpdatePromptWords(c *gin.Context) {
	if secret := c.GetHeader("secretKey"); secret != DecryptedSecretKey {
		middleware.Logger.Printf("[UpdatePromptWords] unauthorized access, clientIP=%s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "unauthorized"})
		return
	}

	var req PromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Logger.Printf("[UpdatePromptWords] invalid json: %v, clientIP=%s", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid json"})
		return
	}

	// 写入 Redis
	err := model.SetPromptWords(req.Prompt)
	if err != nil {
		middleware.Logger.Printf("[UpdatePromptWords] failed to set redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "redis write failed"})
		return
	}

	// 写入本地文件
	err = os.WriteFile("Promptwords", []byte(req.Prompt), 0644)
	if err != nil {
		middleware.Logger.Printf("[UpdatePromptWords] failed to write file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "file write failed"})
		return
	}

	middleware.Logger.Printf("[UpdatePromptWords] success. prompt updated.")
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "prompt updated"})
}
