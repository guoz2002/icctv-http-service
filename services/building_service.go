package services

// BuildingService Methods:
//0. NewBuildingService(db *gorm.DB) -> 初始化建筑服务
//1. List(ctx context.Context) -> 查询所有建筑及其设备
//2. Create(ctx context.Context, payload models.Building) -> 创建建筑
//3. Update(ctx context.Context, id int64, payload models.Building) -> 更新建筑信息
//4. Delete(ctx context.Context, id int64) -> 删除建筑

import (
	"context"
	"errors"

	"icctv-http-service/models"

	"gorm.io/gorm"
)

// 错误定义
var (
	ErrBuildingNotFound = errors.New("building not found")
	ErrOrangePiNotFound = errors.New("orangepi not found")
	ErrNVRNotFound      = errors.New("nvr not found")
	ErrAlreadyBound     = errors.New("already bound to another building")
	ErrNotBound         = errors.New("not bound to any building")
)

// BuildingServiceInterface 定义建筑业务能力
type BuildingServiceInterface interface {
	List(ctx context.Context) ([]models.Building, error)                                     //1.查询建筑列表
	Create(ctx context.Context, payload models.Building) (*models.Building, error)           //2.创建建筑
	Update(ctx context.Context, id int64, payload models.Building) (*models.Building, error) //3.更新建筑
	Delete(ctx context.Context, id int64) error                                              //4.删除建筑
	BindOrangePi(ctx context.Context, buildingId int64, orangePiId int64) error              //5.绑定OrangePi到建筑
	UnbindOrangePi(ctx context.Context, orangePiId int64) error                              //6.解绑OrangePi
	UpdateBind(ctx context.Context, orangePiId int64, newBuildingId int64) error             //7.更新绑定关系
	GetOrangePisByBuildingID(ctx context.Context, buildingId int64) ([]models.OrangePi, error) //8.查询Building关联的OrangePi
	BindNVR(ctx context.Context, buildingId int64, nvrId int64) error                        //9.绑定NVR到建筑
	UnbindNVR(ctx context.Context, nvrId int64) error                                        //10.解绑NVR
	GetNVRsByBuildingID(ctx context.Context, buildingId int64) ([]models.NVR, error)         //11.查询Building关联的NVR
}

// BuildingService 建筑业务逻辑
type BuildingService struct {
	db *gorm.DB
}

// 0. NewBuildingService 构造函数
func NewBuildingService(db *gorm.DB) *BuildingService {
	return &BuildingService{db: db}
}

// 1. List 建筑列表
func (s *BuildingService) List(ctx context.Context) ([]models.Building, error) {
	var items []models.Building
	if err := s.db.WithContext(ctx).Preload("OrangePis").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// 2. Create 创建建筑
func (s *BuildingService) Create(ctx context.Context, payload models.Building) (*models.Building, error) {
	if err := s.db.WithContext(ctx).Create(&payload).Error; err != nil {
		return nil, err
	}
	return &payload, nil
}

// 3. Update 更新建筑
func (s *BuildingService) Update(ctx context.Context, id int64, payload models.Building) (*models.Building, error) {
	var item models.Building
	if err := s.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}

	if payload.ISmartID != "" {
		item.ISmartID = payload.ISmartID
	}
	if payload.Name != "" {
		item.Name = payload.Name
	}
	if payload.Remark != "" {
		item.Remark = payload.Remark
	}

	if err := s.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// 4. Delete 删除建筑
func (s *BuildingService) Delete(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&models.Building{}, id).Error
}

// 5. BindOrangePi 绑定OrangePi到建筑
func (s *BuildingService) BindOrangePi(ctx context.Context, buildingId int64, orangePiId int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查建筑是否存在
		var building models.Building
		if err := tx.First(&building, buildingId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBuildingNotFound
			}
			return err
		}

		// 检查建筑是否有有效的 ISmartID
		if building.ISmartID == "" {
			return errors.New("building has no ismart_id")
		}

		// 检查OrangePi是否存在
		var orangePi models.OrangePi
		if err := tx.First(&orangePi, orangePiId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrangePiNotFound
			}
			return err
		}

		// 检查OrangePi是否已绑定到其他建筑
		if orangePi.ISmartID != "" && orangePi.ISmartID != building.ISmartID {
			return ErrAlreadyBound
		}

		// 如果已经绑定到当前建筑，直接返回成功
		if orangePi.ISmartID == building.ISmartID {
			return nil
		}

		// 更新OrangePi的ISmartID
		orangePi.ISmartID = building.ISmartID
		return tx.Save(&orangePi).Error
	})
}

// 6. UnbindOrangePi 解绑OrangePi
func (s *BuildingService) UnbindOrangePi(ctx context.Context, orangePiId int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查OrangePi是否存在
		var orangePi models.OrangePi
		if err := tx.First(&orangePi, orangePiId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrangePiNotFound
			}
			return err
		}

		// 检查OrangePi是否已绑定
		if orangePi.ISmartID == "" {
			return ErrNotBound
		}

		// 清空ISmartID（设为空字符串）
		orangePi.ISmartID = ""
		return tx.Save(&orangePi).Error
	})
}

// 7. UpdateBind 更新绑定关系
func (s *BuildingService) UpdateBind(ctx context.Context, orangePiId int64, newBuildingId int64) error {
	// 检查新建筑是否存在
	var newBuilding models.Building
	if err := s.db.WithContext(ctx).First(&newBuilding, newBuildingId).Error; err != nil {
		return err
	}

	// 检查OrangePi是否存在
	var orangePi models.OrangePi
	if err := s.db.WithContext(ctx).First(&orangePi, orangePiId).Error; err != nil {
		return err
	}

	// 更新OrangePi的ISmartID
	orangePi.ISmartID = newBuilding.ISmartID
	return s.db.WithContext(ctx).Save(&orangePi).Error
}

// 8. GetOrangePisByBuildingID 查询Building关联的所有OrangePi
func (s *BuildingService) GetOrangePisByBuildingID(ctx context.Context, buildingId int64) ([]models.OrangePi, error) {
	var building models.Building
	if err := s.db.WithContext(ctx).Preload("OrangePis").First(&building, buildingId).Error; err != nil {
		return nil, err
	}
	return building.OrangePis, nil
}

// 9. BindNVR 绑定NVR到建筑
func (s *BuildingService) BindNVR(ctx context.Context, buildingId int64, nvrId int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查建筑是否存在
		var building models.Building
		if err := tx.First(&building, buildingId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBuildingNotFound
			}
			return err
		}

		// 检查NVR是否存在
		var nvr models.NVR
		if err := tx.First(&nvr, nvrId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNVRNotFound
			}
			return err
		}

		// 检查NVR是否已绑定到其他建筑
		if nvr.BuildingID != 0 && nvr.BuildingID != buildingId {
			return ErrAlreadyBound
		}

		// 如果已经绑定到当前建筑，直接返回成功
		if nvr.BuildingID == buildingId {
			return nil
		}

		// 更新NVR的BuildingID
		nvr.BuildingID = buildingId
		return tx.Save(&nvr).Error
	})
}

// 10. UnbindNVR 解绑NVR
func (s *BuildingService) UnbindNVR(ctx context.Context, nvrId int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查NVR是否存在
		var nvr models.NVR
		if err := tx.First(&nvr, nvrId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNVRNotFound
			}
			return err
		}

		// 检查NVR是否已绑定
		if nvr.BuildingID == 0 {
			return ErrNotBound
		}

		// 清空BuildingID（设为0）
		nvr.BuildingID = 0
		return tx.Save(&nvr).Error
	})
}

// 11. GetNVRsByBuildingID 查询Building关联的所有NVR
func (s *BuildingService) GetNVRsByBuildingID(ctx context.Context, buildingId int64) ([]models.NVR, error) {
	var building models.Building
	if err := s.db.WithContext(ctx).Preload("NVRs").First(&building, buildingId).Error; err != nil {
		return nil, err
	}
	return building.NVRs, nil
}
