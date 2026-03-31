package api

import (
	"net/http"
	"roomsync/models"
	"roomsync/repository"
	"roomsync/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomReq struct {
	Name           string `json:"name" binding:"required"`
	Capacity       int    `json:"capacity" binding:"required"`
	Location       string `json:"location"`
	Equipment      string `json:"equipment"`
	NeedApproval   bool   `json:"need_approval"`
	MinAdvanceTime int    `json:"min_advance_time"`
}

func CreateRoom(c *gin.Context) {
	var req RoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 最小值为 5 分钟
	if req.MinAdvanceTime < 5 {
		req.MinAdvanceTime = 5
	}

	room := models.Room{
		Name:           req.Name,
		Capacity:       req.Capacity,
		Location:       req.Location,
		Equipment:      req.Equipment,
		NeedApproval:   req.NeedApproval,
		MinAdvanceTime: req.MinAdvanceTime,
	}

	if err := repository.DB.Create(&room).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create room")
		return
	}
	utils.Success(c, room)
}

func GetRooms(c *gin.Context) {
	var rooms []models.Room
	// Simple query, could be enhanced with pagination
	if err := repository.DB.Find(&rooms).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch rooms")
		return
	}
	utils.Success(c, rooms)
}

func UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var req RoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var room models.Room
	if err := repository.DB.First(&room, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Room not found")
		return
	}

	room.Name = req.Name
	room.Capacity = req.Capacity
	room.Location = req.Location
	room.Equipment = req.Equipment
	room.NeedApproval = req.NeedApproval
	room.MinAdvanceTime = req.MinAdvanceTime

	if err := repository.DB.Save(&room).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update room")
		return
	}
	utils.Success(c, room)
}

func DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	if err := repository.DB.Delete(&models.Room{}, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete room")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

func GetRoomBookings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var bookings []models.Booking
	// 仅返回已批准或待审批的预约，这些是会占用时间的
	if err := repository.DB.Where("room_id = ? AND status IN ?", id, []string{"approved", "pending"}).Find(&bookings).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch room bookings")
		return
	}
	utils.Success(c, bookings)
}
