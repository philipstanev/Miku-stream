package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"

	authdb "github.com/philipstanev/Miku-stream/internal/auth/db"
	"github.com/philipstanev/Miku-stream/internal/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Q *authdb.Queries
	S sessions.Session
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
	return s.S.PutUser(sessions.SessionInfo{UserID: user.ID.String(), Role: user.Role})

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
