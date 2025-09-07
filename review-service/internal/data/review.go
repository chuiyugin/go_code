package data

import (
	"context"
	"errors"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewreviewerRepo .
func NewreviewerRepo(data *Data, logger log.Logger) biz.ReviewerRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

// GetReviewByOrderID 根据订单ID查询评价
func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.OrderID.Eq(orderID)).
		Find()
}

// SaveReply 保存评价回复
func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	// 1、数据校验
	// 1.1 数据合法性校验（已经回复的不允许商家再次回复）
	// 先用评价ID查库，看下是否已评价
	review, err := r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).
		First()
	if err != nil {
		return nil, err
	}
	if review.HasReply == 1 {
		return nil, errors.New("该评价已经回复了")
	}
	// 1.2 水平越权校验（A商家不能够回复B商家的）
	if review.StoreID != reply.StoreID {
		return nil, errors.New("水平越权")
	}
	// 2、更新数据库中的数据（评价表和评价回复表要同时更新）
	// 事务操作
	r.data.query.Transaction(func(tx *query.Query) error {
		// 回复表插入一条数据
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply create reply fail,err:%v\n", err)
			return err
		}
		// 评价表更新 hasReply 字段
		if _, err := tx.ReviewInfo.WithContext(ctx).
			Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).
			Update(tx.ReviewInfo.HasReply, 1); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply update reply fail,err:%v\n", err)
			return err
		}
		return nil
	})
	// 3、返回
	return reply, err
}
