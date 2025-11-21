package services

// AuthService Methods:
//0. NewAuthService(adminService *AdminService) -> 初始化认证服务依赖
//1. PublicToken() -> 返回公开访问 Token
//2. Login(ctx context.Context, username, password string) -> 校验管理员并签发 JWT
//3. ValidateToken(tokenStr string) -> 解析并验证 JWT

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"icctv-http-service/models"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthServiceInterface 定义认证业务能力
type AuthServiceInterface interface {
	GenerateVideoToken(ctx context.Context, buildingID string, channels []string) (string, error) //1.生成视频访问Token
	Login(ctx context.Context, username, password string) (*AuthToken, error)                     //2.校验账号密码并签发JWT
	ValidateToken(tokenStr string) (*AdminClaims, error)                                          //3.验证JWT并返回Claims
}

// AuthService 认证相关逻辑
type AuthService struct {
	db                  *gorm.DB
	adminService        *AdminService
	orangePiService     *OrangePiService
	buildingService     *BuildingService
	videoTokenSecretKey string
	jwtSecret           []byte
	tokenTTL            time.Duration
}

// AuthToken JWT 返回体
type AuthToken struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// AdminClaims JWT Claims
type AdminClaims struct {
	AdminID  int64  `json:"adminId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// VideoTokenPayload 视频Token的Payload结构
type VideoTokenPayload struct {
	Channels   []string `json:"channels"`
	BuildingID string   `json:"building_id"`
	Exp        int64    `json:"exp"`
	Iat        int64    `json:"iat"`
}

// 0. NewAuthService 构造函数，加载基础配置
func NewAuthService(db *gorm.DB, adminService *AdminService, orangePiService *OrangePiService, buildingService *BuildingService) *AuthService {
	// 读取JWT秘钥
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "icctv-secret"
	}

	// 读取视频Token秘钥
	videoSecret := os.Getenv("VIDEO_TOKEN_SECRET_KEY")
	if videoSecret == "" {
		videoSecret = "&gfbrffuk8_$f(ip1&7@_jt#7d#ocn&xdq7^)ctx8k$^a%k(7e"
	}

	// 读取JWT TTL
	ttlMinutes := 120
	if val := os.Getenv("JWT_TTL_MINUTES"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			ttlMinutes = parsed
		}
	}

	return &AuthService{
		db:                  db,
		adminService:        adminService,
		orangePiService:     orangePiService,
		buildingService:     buildingService,
		videoTokenSecretKey: videoSecret,
		jwtSecret:           []byte(secret),
		tokenTTL:            time.Duration(ttlMinutes) * time.Minute,
	}
}

// 1. GenerateVideoToken 生成视频访问 Token
// 参数：buildingID - 建筑ISmartID, channels - 频道列表
// 逻辑：验证 buildingID 是否存在且关联的 OrangePi 设备存在
func (s *AuthService) GenerateVideoToken(ctx context.Context, buildingID string, channels []string) (string, error) {
	// 检查建筑是否存在
	var building models.Building
	if err := s.db.WithContext(ctx).Where("ismart_id = ?", buildingID).First(&building).Error; err != nil {
		return "", fmt.Errorf("building not found: %s", buildingID)
	}

	// 检查该建筑是否关联的 OrangePi 设备存在
	var deviceCount int64
	if err := s.db.WithContext(ctx).Model(&models.OrangePi{}).Where("ismart_id = ?", buildingID).Count(&deviceCount).Error; err != nil {
		return "", fmt.Errorf("failed to check devices: %w", err)
	}

	if deviceCount == 0 {
		return "", fmt.Errorf("no devices associated with building: %s", buildingID)
	}

	// 生成视频 Token (HMAC-SHA256 签名格式)
	return s.generateHMACToken(buildingID, channels)
}

// generateHMACToken 生成 HMAC-SHA256 签名的视频 Token
// 格式：base64(payload).signature
func (s *AuthService) generateHMACToken(buildingID string, channels []string) (string, error) {
	// 创建 Payload
	now := time.Now().Unix()
	expiry := now + 86400 // 24小时有效期

	payload := VideoTokenPayload{
		Channels:   channels,
		BuildingID: buildingID,
		Exp:        expiry,
		Iat:        now,
	}

	// JSON 序列化（键排序保证一致性）
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Base64 编码 payload
	payloadB64 := base64.URLEncoding.EncodeToString(payloadBytes)

	// 生成 HMAC-SHA256 签名
	h := hmac.New(sha256.New, []byte(s.videoTokenSecretKey))
	h.Write(payloadBytes)
	signature := fmt.Sprintf("%x", h.Sum(nil))

	// 组合 token: payload.signature
	token := fmt.Sprintf("%s.%s", payloadB64, signature)
	return token, nil
}

// 2. Login 验证管理员并生成 JWT
func (s *AuthService) Login(ctx context.Context, username, password string) (*AuthToken, error) {
	admin, err := s.adminService.VerifyCredential(ctx, username, password)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.tokenTTL)
	claims := AdminClaims{
		AdminID:  admin.ID,
		Username: admin.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(admin.ID, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken: tokenStr,
		ExpiresAt:   expiresAt,
	}, nil
}

// 3. ValidateToken 校验 JWT
func (s *AuthService) ValidateToken(tokenStr string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AdminClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return s.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
