package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	dir, err := filepath.Abs(filepath.Dir(filepath.Dir(".")))
	if err == nil {
		os.Chdir(dir)
	}
	os.Exit(m.Run())
}

func TestEncryptPassword(t *testing.T) {
	password := "test123456"
	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("EncryptPassword failed: %v", err)
	}

	if len(hash) == 0 {
		t.Fatal("EncryptPassword returned empty hash")
	}

	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword failed for correct password")
	}

	if CheckPassword("wrongpassword", hash) {
		t.Fatal("CheckPassword succeeded for wrong password")
	}
}

func TestBeforeSave_AlreadyHashed(t *testing.T) {
	password := "test123456"
	hash, _ := EncryptPassword(password)

	user := &User{
		Username: "testuser",
		Password: hash,
	}

	err := user.BeforeSave(nil)
	if err != nil {
		t.Fatalf("BeforeSave failed: %v", err)
	}

	if user.Password != hash {
		t.Fatal("BeforeSave re-hashed already hashed password")
	}
}

func TestBeforeSave_NotHashed(t *testing.T) {
	password := "test123456"

	user := &User{
		Username: "testuser",
		Password: password,
	}

	err := user.BeforeSave(nil)
	if err != nil {
		t.Fatalf("BeforeSave failed: %v", err)
	}

	if user.Password == password {
		t.Fatal("BeforeSave did not hash password")
	}

	if !CheckPassword(password, user.Password) {
		t.Fatal("CheckPassword failed after BeforeSave")
	}
}