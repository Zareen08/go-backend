package service

import (
    "errors"
    "spotsync-api/dto"
    "spotsync-api/models"
    "spotsync-api/repository"
    "spotsync-api/utils"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
    // Check if user exists
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, utils.ErrEmailAlreadyExists
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Set default role
    role := req.Role
    if role == "" {
        role = "driver"
    }

    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     role,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      user.Role,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, utils.ErrInvalidCredentials
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, utils.ErrInvalidCredentials
    }

    // Generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return nil, err
    }

    return &dto.LoginResponse{
        Token: tokenString,
        User: dto.UserResponse{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        },
    }, nil
}

func (s *AuthService) ValidateToken(tokenString string) (uint, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return 0, "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint(claims["user_id"].(float64))
        role := claims["role"].(string)
        return userID, role, nil
    }

    return 0, "", errors.New("invalid token")
}