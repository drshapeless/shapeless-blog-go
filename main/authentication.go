package main

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (app *application) createUser(username, password string) error {
	file, err := os.Create(app.userPath(username))
	if err != nil {
		return err
	}

	defer file.Close()

	pw, err := hashPassword(password)
	if err != nil {
		return err
	}

	file.WriteString(pw)
	return nil
}

func (app *application) authenticateUser(username, password string) error {
	bytes, err := os.ReadFile(app.userPath(username))
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(bytes, []byte(password))
	return err
}

func (app *application) writeTokenHash(token string) error {
	file, err := os.Create(app.tokenPath())
	if err != nil {
		return err
	}

	defer file.Close()

	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	if err != nil {
		return err
	}

	app.cache.TokenHash = bytes

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) getTokenHash() ([]byte, error) {
	bytes, err := ioutil.ReadFile(app.tokenPath())
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (app *application) compareToken(token string) error {
	tokenHash, err := app.getTokenHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(tokenHash, []byte(token))
	if err != nil {
		return err
	}
	return nil
}
