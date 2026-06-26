package repository

import (
    "spotsync-api/models"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type ZoneRepository struct {
    db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) *ZoneRepository {
    return &ZoneRepository{db: db}
}

func (r *ZoneRepository) Create(zone *models.ParkingZone) error {
    return r.db.Create(zone).Error
}

func (r *ZoneRepository) FindAll() ([]models.ParkingZone, error) {
    var zones []models.ParkingZone
    err := r.db.Find(&zones).Error
    return zones, err
}

func (r *ZoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
    var zone models.ParkingZone
    err := r.db.First(&zone, id).Error
    return &zone, err
}

func (r *ZoneRepository) Update(zone *models.ParkingZone) error {
    return r.db.Save(zone).Error
}

func (r *ZoneRepository) Delete(id uint) error {
    return r.db.Delete(&models.ParkingZone{}, id).Error
}

func (r *ZoneRepository) GetActiveReservationCount(zoneID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Reservation{}).
        Where("zone_id = ? AND status = ?", zoneID, "active").
        Count(&count).Error
    return count, err
}

func (r *ZoneRepository) FindByIDWithLock(tx *gorm.DB, id uint) (*models.ParkingZone, error) {
    var zone models.ParkingZone
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, id).Error
    return &zone, err
}