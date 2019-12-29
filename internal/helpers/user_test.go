package helpers_test

import (
	"fmt"
	"tbox_backend/internal/helpers"
	"testing"
)

func TestUserHelper_GenerateToken(t *testing.T) {
	userHelper := helpers.NewUserHelper("")
	token, err := userHelper.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(token)
	if token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.N1c2nJo8nCgnpvmf1dYSNq_0tqzl9gs_SBG_aDhVCkE" {
		t.Fatalf("Wrong token")
	}
}

func TestUserHelper_GenerateToken_MultipleTime(t *testing.T) {
	userHelper := helpers.NewUserHelper("abc")
	for i := 0; i < 10; i++ {
		token, err := userHelper.GenerateToken(1)
		if err != nil {
			t.Fatal(err)
		}

		if token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.6dJuEd4gOyQ3aKNHDABXEwQrTYQgNPVDLd0ZcyIUw-Q" {
			t.Fatalf("Wrong token")
		}
	}
}
