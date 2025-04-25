package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Say struct {
	Uid       string `json:"uid"`
	Content   string `json:"content"`
	From      string `json:"from"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
type User struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// List 处理列表请求
// 该方法从数据库中获取分页数据，并将其返回给客户端
func (ctrl *Controller) List(c *gin.Context) {
	// 当前页
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	// 每页数量
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 释放资源
	defer cancel()
	// 执行查询
	rows, err := ctrl.db.QueryContext(
		ctx,
		"SELECT `uuid` as `uid`, `content`, `source` as `from`, `created_at`, `updated_at` FROM `says` ORDER BY `updated_at` DESC LIMIT ? OFFSET ?",
		limit, (page-1)*limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "数据读取失败，" + err.Error(),
		})
		return
	}
	// 释放资源
	defer rows.Close()
	// 定义结果
	var says []Say
	// 遍历结果
	for rows.Next() {
		var say Say
		if err := rows.Scan(&say.Uid, &say.Content, &say.From, &say.CreatedAt, &say.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "数据读取失败，" + err.Error(),
			})
			return
		}
		says = append(says, say)
	}
	// 检查错误
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "数据处理错误，" + err.Error(),
		})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "ok",
		"data":    says,
	})
}

// 用于添加一条说说的接口
// 适配木木老师写的浏览器插件接口(https://github.com/lmm214/memos-bber)，兼容memos的接口
func (ctrl *Controller) Add(c *gin.Context) {
	var userId int64
	if val, exists := c.Get("userid"); exists {
		userId = val.(int64) // 安全获取当前请求的ID
	}
	// 如果用户id依然为空
	if userId == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户不存在",
		})
		return
	}
	// 获取content参数，如果没有提供则默认为空字符串
	content := c.DefaultPostForm("content", "")
	// 检查content是否为空，如果为空则返回错误响应
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "内容不能为空", // 兼容memos的接口
		})
		return
	}
	// 获取visibility参数，如果没有提供则默认为"public"
	visibility := c.DefaultPostForm("visibility", "public")
	visibilityVal := 0
	// 根据visibility值设置visibilityVal
	if visibility == "public" {
		visibilityVal = 1
	}
	// 获取User-Agent头，如果没有提供则默认为"unknown"
	agent := c.GetHeader("User-Agent")
	if agent == "" {
		agent = "unknown"
	}
	// 获取source参数，如果没有提供则默认为User-Agent的值
	source := c.DefaultPostForm("from", "")
	if source == "" {
		source = agent
	}
	// 开始写入数据库
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 生成uuid
	uuid, err := ShortUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据写入失败，" + err.Error(),
		})
		return
	}
	// 执行SQL插入语句
	_, err = ctrl.db.ExecContext(
		ctx,
		"INSERT INTO `says` (`content`, `user_id`, `source`, `agent`, `uuid`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		content, userId, source, agent, uuid, visibilityVal, time.Now().Unix(), time.Now().Unix(),
	)
	// 检查错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据写入失败，" + err.Error(),
		})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (ctrl *Controller) Auth(c *gin.Context) {
	var userId int64
	if val, exists := c.Get("userid"); exists {
		userId = val.(int64) // 安全获取当前请求的ID
	}
	// 如果用户id依然为空
	if userId == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户不存在",
		})
		return
	}
	// 从库里获取用户详情
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := ctrl.db.QueryRowContext(
		ctx,
		"SELECT `id`, `nickname`, `avatar` FROM `users` WHERE `id` = ?",
		userId,
	)
	var user User
	// 扫描结果
	if err := row.Scan(&user.Id, &user.Nickname, &user.Avatar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据读取失败，" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       user.Id,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
	})
}

// Base64编码方案（无损压缩）
// 碰撞概率比原来低10^14个数量级,生成更紧凑的字符串格式
func ShortUUID() (string, error) {
	// 用于生成串的字符串
	const DefaultIDAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// 生成默认22位ID（使用默认字符集）
	return gonanoid.Generate(DefaultIDAlphabet, 22)
}
