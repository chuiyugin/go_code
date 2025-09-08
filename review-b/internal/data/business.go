package data

import (
	"context"
	v1 "review-b/api/review/v1"
	"review-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type bussinesserRepo struct {
	data *Data
	log  *log.Helper
}

// NewBusinesserRepo .
func NewBusinesserRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &bussinesserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *bussinesserRepo) Reply(ctx context.Context, param *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data] Reply, param:%v\n", param)
	// 之前是查询数据库，而此时是需要调用RPC服务来实现
	reply, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("ReplyReview return , reply:%v err:%v", reply, err)
	if err != nil {
		return 0, err
	}
	return reply.ReplyID, err
}
