package service

import (
    "spotsync-api/dto"
    "spotsync-api/models"
    "spotsync-api/repository"
    "spotsync-api/utils"
)

type ZoneService struct {
    zoneRepo *repository.ZoneRepository
}

func NewZoneService(zoneRepo *repository.ZoneRepository) *ZoneService {
    return &ZoneService{zoneRepo: zoneRepo}
}

func (s *ZoneService) CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
    zone := &models.ParkingZone{
        Name:          req.Name,
        Type:          req.Type,
        TotalCapacity: req.TotalCapacity,
        PricePerHour:  req.PricePerHour,
    }

    if err := s.zoneRepo.Create(zone); err != nil {
        return nil, err
    }

    return &dto.ZoneResponse{
        ID:             zone.ID,
        Name:           zone.Name,
        Type:           zone.Type,
        TotalCapacity:  zone.TotalCapacity,
        AvailableSpots: zone.TotalCapacity, // Initially all spots available
        PricePerHour:   zone.PricePerHour,
        CreatedAt:      zone.CreatedAt,
        UpdatedAt:      zone.UpdatedAt,
    }, nil
}

func (s *ZoneService) GetAllZones() ([]dto.ZoneResponse, error) {
    zones, err := s.zoneRepo.FindAll()
    if err != nil {
        return nil, err
    }

    var responses []dto.ZoneResponse
    for _, zone := range zones {
        activeCount, err := s.zoneRepo.GetActiveReservationCount(zone.ID)
        if err != nil {
            return nil, err
        }

        responses = append(responses, dto.ZoneResponse{
            ID:             zone.ID,
            Name:           zone.Name,
            Type:           zone.Type,
            TotalCapacity:  zone.TotalCapacity,
            AvailableSpots: zone.TotalCapacity - int(activeCount),
            PricePerHour:   zone.PricePerHour,
            CreatedAt:      zone.CreatedAt,
            UpdatedAt:      zone.UpdatedAt,
        })
    }

    return responses, nil
}

func (s *ZoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
    zone, err := s.zoneRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    activeCount, err := s.zoneRepo.GetActiveReservationCount(zone.ID)
    if err != nil {
        return nil, err
    }

    return &dto.ZoneResponse{
        ID:             zone.ID,
        Name:           zone.Name,
        Type:           zone.Type,
        TotalCapacity:  zone.TotalCapacity,
        AvailableSpots: zone.TotalCapacity - int(activeCount),
        PricePerHour:   zone.PricePerHour,
        CreatedAt:      zone.CreatedAt,
        UpdatedAt:      zone.UpdatedAt,
    }, nil
}

func (s *ZoneService) UpdateZone(id uint, req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
    zone, err := s.zoneRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    zone.Name = req.Name
    zone.Type = req.Type
    zone.TotalCapacity = req.TotalCapacity
    zone.PricePerHour = req.PricePerHour

    if err := s.zoneRepo.Update(zone); err != nil {
        return nil, err
    }

    activeCount, err := s.zoneRepo.GetActiveReservationCount(zone.ID)
    if err != nil {
        return nil, err
    }

    return &dto.ZoneResponse{
        ID:             zone.ID,
        Name:           zone.Name,
        Type:           zone.Type,
        TotalCapacity:  zone.TotalCapacity,
        AvailableSpots: zone.TotalCapacity - int(activeCount),
        PricePerHour:   zone.PricePerHour,
        CreatedAt:      zone.CreatedAt,
        UpdatedAt:      zone.UpdatedAt,
    }, nil
}

func (s *ZoneService) DeleteZone(id uint) error {
    return s.zoneRepo.Delete(id)
}