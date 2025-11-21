package services

// AdminService Methods:
//0. NewAdminService(db *gorm.DB) -> 初始化管理员服务
//1. List(ctx context.Context, query models.PaginationQuery) -> 分页查询管理员
//2. GetByID(ctx context.Context, id int64) -> 根据 ID 获取管理员
//3. Create(ctx context.Context, username, password string) -> 新增管理员
//4. Update(ctx context.Context, id int64, username, password string) -> 更新管理员
//5. Delete(ctx context.Context, id int64) -> 删除管理员
//6. VerifyCredential(ctx context.Context, username, password string) -> 校验登录凭证

import (
	"context"
	"errors"
	"strings"

	"icctv-http-service/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminServiceInterface 定义管理员业务能力
type AdminServiceInterface interface {
	List(ctx context.Context, query models.PaginationQuery) ([]models.Adminer, models.PaginationResult, error) //1.分页查询管理员
	GetByID(ctx context.Context, id int64) (*models.Adminer, error)                                            //2.根据 ID 获取
	Create(ctx context.Context, username, password string) (*models.Adminer, error)                            //3.新增管理员
	Update(ctx context.Context, id int64, username, password string) (*models.Adminer, error)                  //4.更新管理员
	Delete(ctx context.Context, id int64) error                                                                //5.删除管理员
	VerifyCredential(ctx context.Context, username, password string) (*models.Adminer, error)                  //6.校验凭证
}

// AdminService 管理员业务逻辑
type AdminService struct {
	db *gorm.DB
}

// 0. NewAdminService 构造函数
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// 1. List 查询管理员列表（带分页）
func (s *AdminService) List(ctx context.Context, query models.PaginationQuery) ([]models.Adminer, models.PaginationResult, error) {
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}

	var (
		admins []models.Adminer
		total  int64
	)

	tx := s.db.WithContext(ctx).Model(&models.Adminer{})
	if err := tx.Count(&total).Error; err != nil {
		return nil, models.PaginationResult{}, err
	}

	order := "created_at desc"
	if query.Asc {
		order = "created_at asc"
	}

	if err := tx.Order(order).
		Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&admins).Error; err != nil {
		return nil, models.PaginationResult{}, err
	}

	return admins, models.PaginationResult{
		Total:    int(total),
		PageNum:  query.PageNum,
		PageSize: query.PageSize,
	}, nil
}

// 2. GetByID 根据ID获取管理员
func (s *AdminService) GetByID(ctx context.Context, id int64) (*models.Adminer, error) {
	var admin models.Adminer
	if err := s.db.WithContext(ctx).First(&admin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &admin, nil
}

// 3. Create 创建管理员
func (s *AdminService) Create(ctx context.Context, username, password string) (*models.Adminer, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &models.Adminer{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := s.db.WithContext(ctx).Create(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// 4. Update 更新管理员信息
func (s *AdminService) Update(ctx context.Context, id int64, username, password string) (*models.Adminer, error) {
	var admin models.Adminer
	if err := s.db.WithContext(ctx).First(&admin, id).Error; err != nil {
		return nil, err
	}

	if username != "" {
		admin.Username = strings.TrimSpace(username)
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		admin.PasswordHash = string(hash)
	}

	if err := s.db.WithContext(ctx).Save(&admin).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}

// 5. Delete 删除管理员
func (s *AdminService) Delete(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&models.Adminer{}, id).Error
}

// 6. VerifyCredential 校验用户名密码
func (s *AdminService) VerifyCredential(ctx context.Context, username, password string) (*models.Adminer, error) {
	var admin models.Adminer
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &admin, nil
}
