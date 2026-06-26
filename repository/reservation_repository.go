package repository

import (
    "spotsync-api/models"
    "gorm.io/gorm"
)

type ReservationRepository struct {
    db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
    return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(reservation *models.Reservation) error {
    return r.db.Create(reservation).Error
}

func (r *ReservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
    var reservations []models.Reservation
    err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
    return reservations, err
}

func (r *ReservationRepository) FindByID(id uint) (*models.Reservation, error) {
    var reservation models.Reservation
    err := r.db.Preload("User").Preload("Zone").First(&reservation, id).Error
    return &reservation, err
}

func (r *ReservationRepository) FindAll() ([]models.Reservation, error) {
    var reservations []models.Reservation
    err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
    return reservations, err
}

func (r *ReservationRepository) Update(reservation *models.Reservation) error {
    return r.db.Save(reservation).Error
}

func (r *ReservationRepository) GetActiveCountByZone(zoneID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Reservation{}).
        Where("zone_id = ? AND status = ?", zoneID, "active").
        Count(&count).Error
    return count, err
}

func (r *ReservationRepository) CreateWithTransaction(tx *gorm.DB, reservation *models.Reservation) error {
    return tx.Create(reservation).Error
}