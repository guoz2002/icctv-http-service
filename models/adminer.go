package models

// Adminer 管理员账户模型
type Adminer struct {
	ModelFields
	Username     string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"` // 登录用户名(唯一)
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`                    // Bcrypt 哈希，不返回给前端
}

// TableName 指定表名
func (Adminer) TableName() string {
	return "adminers"
}
