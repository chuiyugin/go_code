package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

// ReviewerRepo is a Reviewer repo.
type ReviewerRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
}

// ReviewerUsecase is a Reviewer usecase.
type ReviewerUsecase struct {
	repo ReviewerRepo
	log  *log.Helper
}

// NewReviewerUsecase new a Reviewer usecase.
func NewReviewerUsecase(repo ReviewerRepo, logger log.Logger) *ReviewerUsecase {
	return &ReviewerUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ReviewerUsecase) CreateReview(ctx context.Context, r *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview, req:%v\n", r)
	// 1、数据校验
	// 1.1 参数基础校验（validate 参数校验）
	// 1.2 参数业务校验
	reviews, err := uc.repo.GetReviewByOrderID(ctx, r.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("查询数据库失败")
	}
	if len(reviews) > 0 {
		// 已经评价过了
		// return nil, fmt.Errorf("订单%d已经评价过", r.OrderID)
		return nil, v1.ErrorOrderReviewed("订单%d已经评价过", r.OrderID)
	}
	// 2、生成review ID（雪花算法）
	r.ReviewID = snowflake.GenID()
	// 3、查询订单和商品快照信息
	// 实际业务场景需要查询订单服务和商家服务（比如说通过RPC调用订单服务和商家服务）
	// 4、拼装数据入库
	return uc.repo.SaveReview(ctx, r)
}

// CreateReply 创建评价回复
func (uc *ReviewerUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply, param:%v\n", param)
	// 调用data层创建一个评价回复
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}
