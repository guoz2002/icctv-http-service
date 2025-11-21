package services

// PublicNetService Methods:
//0. NewPublicNetService(db *gorm.DB) -> 初始化公网配置服务
//1. Get(ctx context.Context) -> 获取首条公网配置记录
//2. Update(ctx context.Context, ip string) -> 写入或更新 ExternalIP

import (
	"context"
	"errors"
	"strings"

	"icctv-http-service/models"

	"gorm.io/gorm"
)

// PublicNetServiceInterface 定义公网配置服务能力
type PublicNetServiceInterface interface {
	Get(ctx context.Context) (*models.PublicNetConfig, error)               //1.获取第一条公网配置
	Update(ctx context.Context, ip string) (*models.PublicNetConfig, error) //2.更新或创建公网配置
}

// PublicNetService 公网配置
type PublicNetService struct {
	db *gorm.DB
}

// 0. NewPublicNetService 构造服务实例并注入数据库句柄
func NewPublicNetService(db *gorm.DB) *PublicNetService {
	return &PublicNetService{db: db}
}

// 1. Get 获取当前公网配置（不存在返回 nil）
func (s *PublicNetService) Get(ctx context.Context) (*models.PublicNetConfig, error) {
	var config models.PublicNetConfig
	err := s.db.WithContext(ctx).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// 2. Update 根据传入 IP 创建或更新公网配置
func (s *PublicNetService) Update(ctx context.Context, ip string) (*models.PublicNetConfig, error) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil, errors.New("external ip is required")
	}

	var config models.PublicNetConfig
	err := s.db.WithContext(ctx).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config.ExternalIP = ip
		if err := s.db.WithContext(ctx).Create(&config).Error; err != nil {
			return nil, err
		}
		return &config, nil
	}
	if err != nil {
		return nil, err
	}

	config.ExternalIP = ip
	if err := s.db.WithContext(ctx).Save(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}
