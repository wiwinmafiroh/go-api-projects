package entity

import (
	"05-go-api-with-middleware/pkg/errs"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type accessRole string

const (
	AdminRole accessRole = "admin"
	UserRole  accessRole = "user"
)

var secretKey = "RAHASIA"

var invalidTokenErr = errs.NewUnauthenticatedError("Invalid token")

type User struct {
	Id        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      accessRole `json:"role"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
}

func (u *User) HashPassword() errs.ErrorMessage {
	salt := 8
	password := []byte(u.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, salt)
	if err != nil {
		return errs.NewInternalServerError("Failed to hash password")
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func (u *User) GenerateToken() string {
	claims := u.claimsToken()

	return u.signToken(claims)
}

func (u *User) claimsToken() jwt.MapClaims {
	return jwt.MapClaims{
		"id":    u.Id,
		"email": u.Email,
		"role":  u.Role,
		"exp":   time.Now().Add(time.Hour * 10).Unix(),
	}
}

func (u *User) signToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	stringToken, _ := token.SignedString([]byte(secretKey))

	return stringToken
}

func (u *User) ValidateToken(bearerToken string) errs.ErrorMessage {
	isBearer := strings.HasPrefix(bearerToken, "Bearer")
	if !isBearer {
		return invalidTokenErr
	}

	splitToken := strings.Split(bearerToken, "Bearer ")
	if len(splitToken) != 2 {
		return invalidTokenErr
	}

	stringToken := splitToken[1]

	token, err := u.parseToken(stringToken)
	if err != nil {
		return err
	}

	var mapClaims jwt.MapClaims

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return invalidTokenErr
	} else {
		mapClaims = claims
	}

	err = u.bindTokenToUserEntity(mapClaims)

	return err
}

func (u *User) parseToken(stringToken string) (*jwt.Token, errs.ErrorMessage) {
	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenErr
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, invalidTokenErr
	}

	return token, nil
}

func (u *User) bindTokenToUserEntity(claims jwt.MapClaims) errs.ErrorMessage {
	if id, ok := claims["id"].(float64); !ok {
		return invalidTokenErr
	} else {
		u.Id = uint(id)
	}

	if email, ok := claims["email"].(string); !ok {
		return invalidTokenErr
	} else {
		u.Email = email
	}

	if role, ok := claims["role"].(string); !ok {
		return invalidTokenErr
	} else {
		u.Role = accessRole(role)
	}

	return nil
}
