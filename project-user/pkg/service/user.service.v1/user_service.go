package user_service_v1

import (
	"com.levi/project-common/base"
	"com.levi/project-common/utils"
	"com.levi/project-user/pkg/dao"
	context "context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type UserService struct {
	UnimplementedUserServiceServer
}

// 登陆 Login
func (UserService) Login(ctx context.Context, msg *LoginRequest) (*LoginResponse, error) {
	// 获取参数
	username := msg.Username
	password := msg.Password

	if username == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "username or password is empty")
	}
	user := &dao.User{}
	if err := base.Mysql.Where("name = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	if !utils.BcryptMakeCheck([]byte(password), user.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "failed to check password")
	}
	tokenData, _ := base.JwtService.CreateToken("app", utils.UintToString(user.ID.ID))

	return &LoginResponse{
		Token: tokenData.AccessToken,
	}, nil
}

// GetUser 获取用户信息
func (UserService) GetUser(ctx context.Context, msg *UserRequest) (*UserResponse, error) {
	userId := msg.UserId
	if userId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument %v", userId)
	}
	user, err := dao.UserDao.GetUserInfo(int(userId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found %v", userId)
	}
	return &UserResponse{
		UserId:   utils.UintToInt64(user.ID.ID),
		UserName: user.Name,
	}, nil
}
