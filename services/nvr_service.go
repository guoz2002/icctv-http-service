package services

// NVRService Methods:
//0. NewNVRService(db *gorm.DB) -> 初始化NVR服务
//1. List(ctx context.Context) -> 查询所有NVR
//2. GetByID(ctx context.Context, id int64) -> 根据ID查询NVR
//3. Create(ctx context.Context, payload models.NVR) -> 创建NVR
//4. Update(ctx context.Context, id int64, payload models.NVR) -> 更新NVR信息
//5. Delete(ctx context.Context, id int64) -> 删除NVR

import (
	"context"

	"icctv-http-service/models"

	"gorm.io/gorm"
)

// NVRServiceInterface 定义NVR业务能力
type NVRServiceInterface interface {
	List(ctx context.Context) ([]models.NVR, error)                                //1.查询NVR列表
	GetByID(ctx context.Context, id int64) (*models.NVR, error)                    //2.根据ID查询NVR
	Create(ctx context.Context, payload models.NVR) (*models.NVR, error)           //3.创建NVR
	Update(ctx context.Context, id int64, payload models.NVR) (*models.NVR, error) //4.更新NVR
	Delete(ctx context.Context, id int64) error                                    //5.删除NVR
}

// NVRService NVR业务逻辑
type NVRService struct {
	db *gorm.DB
}

// 0. NewNVRService 构造函数
func NewNVRService(db *gorm.DB) *NVRService {
	return &NVRService{db: db}
}

// 1. List NVR列表
func (s *NVRService) List(ctx context.Context) ([]models.NVR, error) {
	var items []models.NVR
	if err := s.db.WithContext(ctx).Preload("Building").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// 2. GetByID 根据ID查询NVR
func (s *NVRService) GetByID(ctx context.Context, id int64) (*models.NVR, error) {
	var item models.NVR
	if err := s.db.WithContext(ctx).Preload("Building").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// 3. Create 创建NVR
func (s *NVRService) Create(ctx context.Context, payload models.NVR) (*models.NVR, error) {
	if err := s.db.WithContext(ctx).Create(&payload).Error; err != nil {
		return nil, err
	}
	return &payload, nil
}

// 4. Update 更新NVR
func (s *NVRService) Update(ctx context.Context, id int64, payload models.NVR) (*models.NVR, error) {
	var item models.NVR
	if err := s.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}

	if payload.Name != "" {
		item.Name = payload.Name
	}
	if payload.URL != "" {
		item.URL = payload.URL
	}
	if payload.BuildingID != 0 {
		item.BuildingID = payload.BuildingID
	}
	if payload.AdminUser.Name != "" || payload.AdminUser.Password != "" {
		item.AdminUser = payload.AdminUser
	}
	if len(payload.Users) > 0 {
		item.Users = payload.Users
	}
	if len(payload.RTSPUrls) > 0 {
		item.RTSPUrls = payload.RTSPUrls
	}

	if err := s.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// 5. Delete 删除NVR
func (s *NVRService) Delete(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&models.NVR{}, id).Error
}
