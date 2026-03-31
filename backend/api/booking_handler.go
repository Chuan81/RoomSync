package api

import (
	"net/http"
	"roomsync/models"
	"roomsync/repository"
	"roomsync/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingReq struct {
	RoomID    uint      `json:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Title     string    `json:"title" binding:"required"`
}

func CreateBooking(c *gin.Context) {
	var req BookingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		utils.Error(c, http.StatusBadRequest, "Start time must be before end time")
		return
	}

	userID, _ := c.Get("userID")

	// Get room details
	var room models.Room
	if err := repository.DB.First(&room, req.RoomID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Room not found")
		return
	}

	// 提前预订时间校验
	minAdvanceDuration := time.Duration(room.MinAdvanceTime) * time.Minute
	if time.Until(req.StartTime) < minAdvanceDuration {
		utils.Error(c, http.StatusBadRequest, "该会议室要求至少提前 "+strconv.Itoa(room.MinAdvanceTime)+" 分钟预约")
		return
	}

	// 限制活跃预约数量 (pending 或 approved)
	var activeCount int64
	repository.DB.Model(&models.Booking{}).Where(
		"room_id = ? AND user_id = ? AND status IN ?",
		req.RoomID, userID.(uint), []string{"pending", "approved"},
	).Count(&activeCount)
	if activeCount >= int64(room.MaxActiveBookings) {
		utils.Error(c, http.StatusForbidden, "在该房间您已有 "+strconv.Itoa(int(activeCount))+" 条未处理/未开始的预约，达到上限")
		return
	}

	// Conflict detection
	var conflictCount int64
	repository.DB.Model(&models.Booking{}).Where(
		"room_id = ? AND status IN ? AND start_time < ? AND end_time > ?",
		req.RoomID, []string{"pending", "approved"}, req.EndTime, req.StartTime,
	).Count(&conflictCount)

	if conflictCount > 0 {
		utils.Error(c, http.StatusConflict, "Room is already booked for this time slot")
		return
	}

	status := "approved"
	if room.NeedApproval {
		status = "pending"
	}

	booking := models.Booking{
		RoomID:    req.RoomID,
		UserID:    userID.(uint),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title:     req.Title,
		Status:    status,
	}

	if err := repository.DB.Create(&booking).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create booking")
		return
	}

	// TODO: Send notification asynchronously here
	utils.Success(c, booking)
}

func GetMyBookings(c *gin.Context) {
	userID, _ := c.Get("userID")
	var bookings []models.Booking
	if err := repository.DB.Preload("Room").Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	utils.Success(c, bookings)
}

func GetAllBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := repository.DB.Preload("Room").Preload("User").Find(&bookings).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	utils.Success(c, bookings)
}

type ApproveReq struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

func ApproveBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	var req ApproveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var booking models.Booking
	if err := repository.DB.First(&booking, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Booking not found")
		return
	}

	booking.Status = req.Status
	if err := repository.DB.Save(&booking).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update booking status")
		return
	}

	// TODO: Send notification to user about approval/rejection
	utils.Success(c, booking)
}

func CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	userID, _ := c.Get("userID")

	var booking models.Booking
	if err := repository.DB.First(&booking, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Booking not found")
		return
	}

	// Allow users to cancel their own, and admins to cancel any
	role, _ := c.Get("role")
	if role != "admin" && booking.UserID != userID.(uint) {
		utils.Error(c, http.StatusForbidden, "You can only cancel your own bookings")
		return
	}

	booking.Status = "cancelled"
	if err := repository.DB.Save(&booking).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to cancel booking")
		return
	}
	utils.Success(c, booking)
}
