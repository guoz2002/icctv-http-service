package models

// OrangePi OrangePi设备模型
type OrangePi struct {
	ModelFields

	ISmartID                   string `gorm:"type:varchar(100);not null;column:ismart_id;-:migration" json:"ismartid"` // 关联楼栋 ismartId
	Name                       string `gorm:"type:varchar(255);not null" json:"name"`                                  // Orangepi 名称
	ICCTVAuthServiceRemotePort int    `gorm:"not null" json:"icctv_auth_service_remote_port"`                          // 远程认证服务端口
	SSHRemotePort              int    `gorm:"not null" json:"ssh_remote_port"`                                         // SSH 远程端口
	IsActive                   bool   `gorm:"default:true" json:"is_active"`                                           // 是否在用

	// 关联关系
	Building *Building `gorm:"foreignKey:ISmartID;references:ISmartID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"building,omitempty"` // 关联建筑
}

// TableName 指定表名
func (OrangePi) TableName() string {
	return "orangepis"
}
