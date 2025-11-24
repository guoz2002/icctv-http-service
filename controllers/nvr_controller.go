package controllers

// NVRController Methods:
//0. NewNVRController(service *services.NVRService) -> 注入 NVRService
//1. List(w http.ResponseWriter, r *http.Request) -> 列出NVR
//2. Create(w http.ResponseWriter, r *http.Request) -> 创建NVR
//3. Update(w http.ResponseWriter, r *http.Request) -> 更新NVR
//4. Delete(w http.ResponseWriter, r *http.Request) -> 删除NVR

import (
	"net/http"
	"strconv"

	"icctv-http-service/models"
	"icctv-http-service/services"
)

// NVRControllerInterface 定义NVR接口能力
type NVRControllerInterface interface {
	List(w http.ResponseWriter, r *http.Request)   //1.查询NVR列表/详情
	Create(w http.ResponseWriter, r *http.Request) //2.创建NVR
	Update(w http.ResponseWriter, r *http.Request) //3.更新NVR
	Delete(w http.ResponseWriter, r *http.Request) //4.删除NVR
}

// NVRController NVR接口
type NVRController struct {
	service *services.NVRService
}

// 0. NewNVRController 构造函数
func NewNVRController(service *services.NVRService) *NVRController {
	return &NVRController{service: service}
}

// 1. List 查询NVR列表或详情
func (c *NVRController) List(w http.ResponseWriter, r *http.Request) {
	// 检查是否有id参数，如果有则返回详情
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id parameter")
			return
		}
		item, err := c.service.GetByID(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondData(w, http.StatusOK, item)
		return
	}

	// 否则返回列表
	items, err := c.service.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondData(w, http.StatusOK, items)
}

// 2. Create 创建NVR
func (c *NVRController) Create(w http.ResponseWriter, r *http.Request) {
	var req models.NVR
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := c.service.Create(r.Context(), req)
	respondResult(w, err, http.StatusCreated, item)
}

// 3. Update 更新NVR
func (c *NVRController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromQuery(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req models.NVR
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := c.service.Update(r.Context(), id, req)
	respondResult(w, err, http.StatusOK, item)
}

type deleteNVRRequest struct {
	ID int64 `json:"id"`
}

// 4. Delete 删除NVR
func (c *NVRController) Delete(w http.ResponseWriter, r *http.Request) {
	var req deleteNVRRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ID == 0 {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := c.service.Delete(r.Context(), req.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondData(w, http.StatusOK, map[string]bool{"deleted": true})
}
