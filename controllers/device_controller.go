package controllers

// DeviceController Methods:
//0. NewDeviceController(deviceService *services.DeviceService, orangePiService *services.OrangePiService) -> 注入仪表板依赖
//1. Summary(w http.ResponseWriter, r *http.Request) -> 获取设备概览

import (
	"net/http"

	"icctv-http-service/services"
)

// DeviceControllerInterface 定义仪表板接口能力
type DeviceControllerInterface interface {
	Summary(w http.ResponseWriter, r *http.Request) //1.获取设备概览
}

// DeviceController 设备仪表板
type DeviceController struct {
	deviceService   *services.DeviceService
	orangePiService *services.OrangePiService
}

// 0. NewDeviceController 构造函数
func NewDeviceController(deviceService *services.DeviceService, orangePiService *services.OrangePiService) *DeviceController {
	return &DeviceController{
		deviceService:   deviceService,
		orangePiService: orangePiService,
	}
}

// 1. Summary 获取设备信息
func (c *DeviceController) Summary(w http.ResponseWriter, r *http.Request) {
	result, err := c.deviceService.Summary(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondData(w, http.StatusOK, result)
}
