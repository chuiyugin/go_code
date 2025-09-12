package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var secret = []byte("氹氹转菊花园")

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func passwordMd5(password []byte) string {

	h := md5.New()
	h.Write(password) // 密码计算md5
	h.Write(secret)
	PasswordStr := hex.EncodeToString(h.Sum(nil))
	return PasswordStr
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	// todo: add your logic here and delete this line
	// fmt.Printf("req:%#v\n", req)
	logx.Debugv(req) // json.Narshal(req)
	// 参数校验
	if req.Password != req.RePassword {
		return nil, errors.New("两次输入密码不一致")
	}
	// 查询用户名是否存在
	u, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	// 查询数据库失败
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw(
			"user_signup_UserModel.FindOneByUsername failed",
			logx.Field("err", err),
		)
		return nil, errors.New("内部错误")
	}
	// 查询到数据库记录,表示该用户名已被注册
	if u != nil {
		return nil, errors.New("该用户名已被注册")
	}
	// 没查到记录
	// 加密密码
	h := md5.New()
	h.Write([]byte(req.Password)) // 密码计算md5
	h.Write(secret)
	PasswordStr := hex.EncodeToString(h.Sum(nil))
	// 把用户注册信息保存到数据库
	user := &model.Users{
		UserId:   time.Now().Unix(), // 简化
		Username: req.Username,
		Password: PasswordStr,
		Gender:   int64(req.Gender),
	}
	if _, err := l.svcCtx.UsersModel.Insert(context.Background(), user); err != nil {
		logx.Errorf("user_signup_UserModel.FindOneByUsername failed, err:%v", err)
		return nil, err
	}
	return &types.SignupResponse{Message: "success!"}, nil
}
