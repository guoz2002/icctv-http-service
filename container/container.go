package container

import (
	"icctv-http-service/controllers"
	"icctv-http-service/databases"
	"icctv-http-service/middlewares"
	"icctv-http-service/routes"
	"icctv-http-service/services"

	"gorm.io/gorm"
)

// ServiceSet 聚合所有业务服务
type ServiceSet struct {
	Admin     *services.AdminService
	Auth      *services.AuthService
	OrangePi  *services.OrangePiService
	Building  *services.BuildingService
	Device    *services.DeviceService
	PublicNet *services.PublicNetService
}

// Container 提供项目运行所需的依赖
type Container struct {
	DB          *gorm.DB
	Services    ServiceSet
	Controllers routes.ControllerSet
	Middlewares routes.MiddlewareSet
}

// Build 初始化所有依赖
func Build() (*Container, error) {
	db, err := databases.Init()
	if err != nil {
		return nil, err
	}

	serviceSet := ServiceSet{
		Admin:     services.NewAdminService(db),
		Building:  services.NewBuildingService(db),
		Device:    services.NewDeviceService(db),
		PublicNet: services.NewPublicNetService(db),
	}
	serviceSet.OrangePi = services.NewOrangePiService(db, serviceSet.PublicNet)
	serviceSet.Auth = services.NewAuthService(db, serviceSet.Admin, serviceSet.OrangePi, serviceSet.Building)

	ctrlSet := routes.ControllerSet{
		Auth:      controllers.NewAuthController(serviceSet.Auth),
		Admin:     controllers.NewAdminController(serviceSet.Admin),
		OrangePi:  controllers.NewOrangePiController(serviceSet.OrangePi),
		Building:  controllers.NewBuildingController(serviceSet.Building),
		Device:    controllers.NewDeviceController(serviceSet.Device, serviceSet.OrangePi),
		PublicNet: controllers.NewPublicNetController(serviceSet.PublicNet),
	}

	middlewareSet := routes.MiddlewareSet{
		Auth: middlewares.NewAuthMiddleware(serviceSet.Auth),
	}

	return &Container{
		DB:          db,
		Services:    serviceSet,
		Controllers: ctrlSet,
		Middlewares: middlewareSet,
	}, nil
}
