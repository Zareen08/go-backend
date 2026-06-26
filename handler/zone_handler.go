package handler

import (
    "spotsync-api/dto"
    "spotsync-api/service"
    "spotsync-api/utils"
    "strconv"
    "github.com/labstack/echo/v4"
)

type ZoneHandler struct {
    zoneService *service.ZoneService
}

func NewZoneHandler(zoneService *service.ZoneService) *ZoneHandler {
    return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
    var req dto.CreateZoneRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body", err.Error())
    }

    if err := c.Validate(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Validation failed", err.Error())
    }

    zone, err := h.zoneService.CreateZone(&req)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create zone", err.Error())
    }

    return utils.SuccessResponse(c, 201, "Parking zone created successfully", zone)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
    zones, err := h.zoneService.GetAllZones()
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to retrieve zones", err.Error())
    }

    return utils.SuccessResponse(c, 200, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid zone ID", nil)
    }

    zone, err := h.zoneService.GetZoneByID(uint(id))
    if err != nil {
        return utils.ErrorResponse(c, 404, "Zone not found", nil)
    }

    return utils.SuccessResponse(c, 200, "Parking zone retrieved successfully", zone)
}