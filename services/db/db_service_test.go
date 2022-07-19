package db

import (
	"com.nzair.user/air"
	"fmt"
	"strings"
	"testing"
)

func TestService_CreateAndQueryUser(t *testing.T) {
	dbService, _ := NewDBService("./test.json")
	tester := &air.User{ID: 1, FirstName: "tester", LastName: "New Zealand", Email: "tester@gmail.com", Login: "127.0.0.1"}
	err := dbService.CreateUser(tester)
	if err != nil {
		t.Error("Failed to create user")
	}

	tester1 := &air.User{ID: 2, FirstName: "tester1", LastName: "New Zealand", Email: "tester1@gmail.com", Login: "127.0.0.1"}
	err = dbService.CreateUser(tester1)
	if err != nil {
		t.Error("Failed to create user")
	}

	tester2, err := dbService.UserByEmail("tester@gmail.com")
	if err != nil {
		t.Error("Failed to query user")
	}
	if tester2.ID != tester.ID {
		t.Error(fmt.Scanf("User IDs do not match %d != %d", tester2.ID, tester.ID))
	}
	if strings.Compare(tester2.Email, tester.Email) != 0 {
		t.Error(fmt.Scanf("User Emails do not match %s != %s", tester2.Email, tester.Email))
	}
}
