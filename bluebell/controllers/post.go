package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil { // validator --> binding tag
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("Create post with invaild parm", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c 中取到当前发请求的用户ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64) // 将字符串转换为数字，10进制，64位
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
	}
	// 2. 根据id取出帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表接口的处理函数
func GetPostListHandler(c *gin.Context) {
	// 获取分页和数据
	page, size := getPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// 根据前端传来的参数动态获取帖子参数
// 按照创建时间 或者 分数排列
// 1 获取请求的query string参数
// 2 去redis查询id列表
// 3 根据id去数据库查询帖子详细信息

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数（query string）：/api/v1/post2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数（指定成默认的值）
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取分页和数据
	data, err := logic.GetPostListNew(p) // 更新：合二为一
	if err != nil {
		zap.L().Error("logic.GetPostList2() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// // 根据社区去查询帖子列表
// func GetCommunityPostListHandler(c *gin.Context) {
// 	// 初始化结构体时指定初始参数（指定成默认的值）
// 	p := &models.ParamCommunityPostList{
// 		ParamPostList: &models.ParamPostList{
// 			Page:  1,
// 			Size:  10,
// 			Order: models.OrderTime,
// 		},
// 	}
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
// 		ResponseError(c, CodeInvalidParam)
// 		return
// 	}
// 	// 获取分页和数据
// 	data, err := logic.GetCommunityPostList(p)
// 	if err != nil {
// 		zap.L().Error("logic.GetCommunityPostList() failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	// 返回响应
// 	ResponseSuccess(c, data)
// }
