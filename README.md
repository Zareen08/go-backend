# SpotSync API

Smart Parking & EV Charging Reservation System

## Live URL

https://go-backend-x1a8.onrender.com/

## Features

- User registration and authentication with JWT
- Role-based access control (Driver & Admin)
- Parking zone management
- Smart reservation system with race condition prevention
- Real-time availability calculation
- EV charging spot reservation
- PostgreSQL database with GORM
- Clean Architecture implementation

## Tech Stack

- **Go** 1.22+
- **Echo** - Web framework
- **GORM** - ORM with PostgreSQL
- **JWT** - Authentication
- **bcrypt** - Password hashing
- **PostgreSQL** - Database

## Architecture

Handler │────▶│ Service │────▶│ Repository │────▶│ Database │

DTOs (Data Transfer Objects) Models (GORM)

## Local Development

### Prerequisites

- Go 1.22+
- PostgreSQL

### Setup

1. Clone the repository:
bash
git clone https://github.com/yourusername/spotsync-api.git
cd spotsync-api

2. Install dependencies:
bash
go mod download

3. Create .env file

4. Run the application
bash
go run cmd/main.go

### API Endpoints

**Authentication**

Method	Endpoint	Description	Access
POST	/api/v1/auth/register	Register new user	Public
POST	/api/v1/auth/login	Login user	Public

**Parking Zones**

Method	Endpoint	Description	Access
GET	/api/v1/zones	Get all zones	Public
GET	/api/v1/zones/:id	Get zone by ID	Public
POST	/api/v1/zones	Create zone	Admin

**Reservations**

Method	Endpoint	Description	Access
POST	/api/v1/reservations	Create reservation	Authenticated
GET	/api/v1/reservations/my-reservations	Get my reservations	Authenticated
DELETE	/api/v1/reservations/:id	Cancel reservation	Authenticated
GET	/api/v1/reservations	Get all reservations	Admin