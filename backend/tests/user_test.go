package main

import (
	"testing"

	"216/internal/orchestrator/Entities"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestHash(t *testing.T) {
	password := "secretpassword"
	hashedPassword, err := Entities.Hash(password)
	assert.Nil(t, err)
	assert.NotNil(t, hashedPassword)
}

func TestVerifyPassword(t *testing.T) {
	password := "secretpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.Nil(t, err)

	err = Entities.VerifyPassword(string(hashedPassword), password)
	assert.Nil(t, err)
}

func TestBeforeSave(t *testing.T) {
	user := Entities.User{Name: "John", Email: "john@example.com", Password: "secretpassword"}
	db := &gorm.DB{}
	err := user.BeforeSave(db)
	assert.Nil(t, err)
	assert.NotEqual(t, user.ID, uuid.Nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Password)
}

func TestPrepare(t *testing.T) {
	user := Entities.User{Name: "   John   ", Email: "   john@example.com   "}
	user.Prepare()
	assert.Equal(t, user.Name, "John")
	assert.Equal(t, user.Email, "john@example.com")
}

func TestValidateUpdate(t *testing.T) {
	user := Entities.User{Name: "John", Email: "john@example.com", Password: "secretpassword"}
	err := user.Validate("update")
	assert.Nil(t, err)

	user.Name = ""
	err = user.Validate("update")
	assert.NotNil(t, err)
}

func TestValidateLogin(t *testing.T) {
	user := Entities.User{Name: "John", Email: "john@example.com", Password: "secretpassword"}
	err := user.Validate("login")
	assert.Nil(t, err)

	user.Email = ""
	err = user.Validate("login")
	assert.NotNil(t, err)
}

func TestValidateDefault(t *testing.T) {
	user := Entities.User{Name: "John", Email: "john@example.com", Password: "secretpassword"}
	err := user.Validate("default")
	assert.Nil(t, err)

	user.Password = ""
	err = user.Validate("default")
	assert.NotNil(t, err)
}
