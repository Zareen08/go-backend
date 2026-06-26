package handler

import (
    "spotsync-api/dto"
    "spotsync-api/service"
    "spotsync-api/utils"
    "strconv"
    "github.com/labstack/echo/v4"
)

type ReservationHandler struct {
    reservationService *service.ReservationService
}

func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
    return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
    userID := c.Get("user_id").(uint)

    var req dto.CreateReservationRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body", err.Error())
    }

    if err := c.Validate(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed", err.Error())
    }

    reservation, err := h.reservationService.CreateReservation(userID, &req)
    if err != nil {
        if err == utils.ErrZoneFull {
            return utils.ErrorResponse(c, 409, "Zone is at full capacity", nil)
        }
        return utils.ErrorResponse(c, 500, "Failed to create reservation", err.Error())
    }

    return utils.SuccessResponse(c, 201, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
    userID := c.Get("user_id").(uint)

    reservations, err := h.reservationService.GetMyReservations(userID)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to retrieve reservations", err.Error())
    }

    return utils.SuccessResponse(c, 200, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
    userID := c.Get("user_id").(uint)
    isAdmin := c.Get("role").(string) == "admin"

    reservationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid reservation ID", nil)
    }

    if err := h.reservationService.CancelReservation(userID, uint(reservationID), isAdmin); err != nil {
        if err == utils.ErrForbidden {
            return utils.ErrorResponse(c, 403, "You can only cancel your own reservations", nil)
        }
        if err == utils.ErrInvalidReservationStatus {
            return utils.ErrorResponse(c, 400, "Only active reservations can be cancelled", nil)
        }
        return utils.ErrorResponse(c, 500, "Failed to cancel reservation", err.Error())
    }

    return utils.SuccessResponse(c, 200, "Reservation cancelled successfully", nil)
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
    reservations, err := h.reservationService.GetAllReservations()
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to retrieve reservations", err.Error())
    }

    return utils.SuccessResponse(c, 200, "All reservations retrieved successfully", reservations)
}