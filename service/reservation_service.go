package service

import (
    "spotsync-api/dto"
    "spotsync-api/models"
    "spotsync-api/repository"
    "spotsync-api/utils"
    "gorm.io/gorm"
)

type ReservationService struct {
    db               *gorm.DB
    reservationRepo  *repository.ReservationRepository
    zoneRepo         *repository.ZoneRepository
}

func NewReservationService(db *gorm.DB, reservationRepo *repository.ReservationRepository, zoneRepo *repository.ZoneRepository) *ReservationService {
    return &ReservationService{
        db:              db,
        reservationRepo: reservationRepo,
        zoneRepo:        zoneRepo,
    }
}

// CRITICAL: This method handles the race condition with row locking
func (s *ReservationService) CreateReservation(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
    var reservation *models.Reservation
    var zoneResponse *dto.ZoneResponse

    // Use transaction with row locking to prevent race conditions
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // 1. Lock the parking zone row for update
        zone, err := s.zoneRepo.FindByIDWithLock(tx, req.ZoneID)
        if err != nil {
            return err
        }

        // 2. Count active reservations for this zone within the transaction
        var activeCount int64
        if err := tx.Model(&models.Reservation{}).
            Where("zone_id = ? AND status = ?", req.ZoneID, "active").
            Count(&activeCount).Error; err != nil {
            return err
        }

        // 3. Check if zone is full
        if int(activeCount) >= zone.TotalCapacity {
            return utils.ErrZoneFull
        }

        // 4. Create the reservation
        reservation = &models.Reservation{
            UserID:       userID,
            ZoneID:       req.ZoneID,
            LicensePlate: req.LicensePlate,
            Status:       "active",
        }

        if err := tx.Create(reservation).Error; err != nil {
            return err
        }

        // 5. Prepare response data
        zoneResponse = &dto.ZoneResponse{
            ID:             zone.ID,
            Name:           zone.Name,
            Type:           zone.Type,
            TotalCapacity:  zone.TotalCapacity,
            AvailableSpots: zone.TotalCapacity - int(activeCount) - 1,
            PricePerHour:   zone.PricePerHour,
            CreatedAt:      zone.CreatedAt,
            UpdatedAt:      zone.UpdatedAt,
        }

        return nil // Commit transaction
    })

    if err != nil {
        return nil, err
    }

    return &dto.ReservationResponse{
        ID:           reservation.ID,
        UserID:       reservation.UserID,
        ZoneID:       reservation.ZoneID,
        LicensePlate: reservation.LicensePlate,
        Status:       reservation.Status,
        CreatedAt:    reservation.CreatedAt,
        UpdatedAt:    reservation.UpdatedAt,
        Zone: &dto.ZoneInfo{
            ID:   zoneResponse.ID,
            Name: zoneResponse.Name,
            Type: zoneResponse.Type,
        },
    }, nil
}

func (s *ReservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
    reservations, err := s.reservationRepo.FindByUserID(userID)
    if err != nil {
        return nil, err
    }

    var responses []dto.ReservationResponse
    for _, res := range reservations {
        responses = append(responses, dto.ReservationResponse{
            ID:           res.ID,
            UserID:       res.UserID,
            ZoneID:       res.ZoneID,
            LicensePlate: res.LicensePlate,
            Status:       res.Status,
            CreatedAt:    res.CreatedAt,
            UpdatedAt:    res.UpdatedAt,
            Zone: &dto.ZoneInfo{
                ID:   res.Zone.ID,
                Name: res.Zone.Name,
                Type: res.Zone.Type,
            },
        })
    }

    return responses, nil
}

func (s *ReservationService) CancelReservation(userID uint, reservationID uint, isAdmin bool) error {
    reservation, err := s.reservationRepo.FindByID(reservationID)
    if err != nil {
        return err
    }

    // Check permission: only owner or admin can cancel
    if !isAdmin && reservation.UserID != userID {
        return utils.ErrForbidden
    }

    // Can only cancel active reservations
    if reservation.Status != "active" {
        return utils.ErrInvalidReservationStatus
    }

    // Update status to cancelled
    reservation.Status = "cancelled"
    return s.reservationRepo.Update(reservation)
}

func (s *ReservationService) GetAllReservations() ([]dto.ReservationResponse, error) {
    reservations, err := s.reservationRepo.FindAll()
    if err != nil {
        return nil, err
    }

    var responses []dto.ReservationResponse
    for _, res := range reservations {
        responses = append(responses, dto.ReservationResponse{
            ID:           res.ID,
            UserID:       res.UserID,
            ZoneID:       res.ZoneID,
            LicensePlate: res.LicensePlate,
            Status:       res.Status,
            CreatedAt:    res.CreatedAt,
            UpdatedAt:    res.UpdatedAt,
            Zone: &dto.ZoneInfo{
                ID:   res.Zone.ID,
                Name: res.Zone.Name,
                Type: res.Zone.Type,
            },
            User: &dto.UserInfo{
                ID:    res.User.ID,
                Name:  res.User.Name,
                Email: res.User.Email,
            },
        })
    }

    return responses, nil
}