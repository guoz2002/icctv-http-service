package models

// PublicNetConfig 公网配置模型
type PublicNetConfig struct {
	ModelFields
	ExternalIP string `gorm:"type:varchar(50);not null" json:"external_ip"` // 外部IP地址
}

// TableName 指定表名
func (PublicNetConfig) TableName() string {
	return "public_net_configs"
}
