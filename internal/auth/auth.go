package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	//"github.com/jackc/pgx/v5"
	//"github.com/jackc/pgx/v5/pgtype"
	//"github.com/jackc/pgx/v5/pgxpool"
	authdb "github.com/philipstanev/Miku-stream/internal/auth/db"
)

type Service struct {
	Q *authdb.Queries
	S Session
}

type SessionInfo struct {
	userID string
	role   string
}

type Session interface {
	GetInfo(sessionID string) (SessionInfo, error)
	PutUser(info SessionInfo) (string, error)
	RemoveUser(sessionID string) error
}

func (s *Service) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.Q.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	stored := user.PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(stored), []byte(password))
	if err != nil {
		return "", err
	}
	return s.S.PutUser(SessionInfo{userID: user.ID.String(), role: user.Role})

}

// nil means passwords match
func (s *Service) checkPassword(ctx context.Context, username string, password string) error {
	user, err := s.Q.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	stored := user.PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(stored), []byte(password))
	return err

}

func hashPassword(password string) string {
	var byteOriginal []byte = []byte(password)
	strCost := os.Getenv("PASSWORD_HASH_COST")
	fmt.Println(strCost)
	cost, err := strconv.Atoi(strCost)
	if err != nil {
		panic(err)
	}
	byteHash, err := bcrypt.GenerateFromPassword(byteOriginal, cost)
	if err != nil {
		panic(err)
	}
	return string(byteHash)
}

func (s *Service) Register(ctx context.Context, password string, username string) (string, error) {
	hash := hashPassword(password)
	pgID, err := s.Q.InsertUser(ctx, authdb.InsertUserParams{PasswordHash: hash, Username: username})
	if err != nil {
		return "", err
	}

	return pgID.String(), err
}
