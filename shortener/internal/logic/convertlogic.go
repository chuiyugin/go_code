package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链业务逻辑：输入一个长链接 --> 转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// todo: add your logic here and delete this line
	// 1 校验输入的数据
	// 1.1 数据不能为空（使用validator包进行参数校验）--> 放在converthandler中
	// 1.2 输入的长链接必须是一个能请求通的网址
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("无效的链接")
	}
	// 1.3 判断之前是否已经转链过（数据库中是否已经存在该长链接）
	// 1.3.1 给长链接生成MD5值
	md5Value := md5.Sum([]byte(req.LongUrl)) // 使用项目中封装的md5
	// 1.3.2 根据MD5值去数据库中查是否存在（sql.NullString 可空字符串字段）
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err != nil {
			return nil, fmt.Errorf("该链接已经被转为%s", u.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 1.4 输入的不能是一个短链接（避免循环转链）
	basePath, err := urltool.GetbasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("connect.GetbasePath(req.LongUrl) failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	u2, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err != nil {
			return nil, fmt.Errorf("该链接已经是短链接%s", u2.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	var short string
	for {
		// 2 取号 基于MySQL实现的发号器
		// 每来一个转链请求，就是用 REPLACE INTO 语句往 sequence 表插入一条数据，并且取出主键id作为号码
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence.Next failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		fmt.Println(seq)
		// 3 号码转短链
		// 3.1 安全性(打乱顺序)
		short = base62.IntToBase62(seq)
		// 3.2 避免某些特殊的词
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break // 生成不在黑名单里的短链接就跳出循环
		}
	}
	// 4 存储长短链接映射关系
	if _, err := l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongUrl, Valid: true},
			Md5:  sql.NullString{String: md5Value, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
		},
	); err != nil {
		logx.Errorw("ShortUrlModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 5 返回响应
	// 返回的是 短域名+短链接
	shortUrl := l.svcCtx.Config.ShortDoamin + "/" + short
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
