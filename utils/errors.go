package utils

import "errors"

var (
    ErrEmailAlreadyExists       = errors.New("email already exists")
    ErrInvalidCredentials       = errors.New("invalid credentials")
    ErrZoneFull                 = errors.New("zone is at full capacity")
    ErrForbidden                = errors.New("forbidden")
    ErrInvalidReservationStatus = errors.New("invalid reservation status")
    ErrReservationNotFound      = errors.New("reservation not found")
)