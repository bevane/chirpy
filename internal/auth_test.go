package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID, _ := uuid.Parse("4aeab7c0-2963-42b3-b420-4007022aabee")
	_, err := MakeJWT(userID, "abcdef", time.Hour)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	userID, _ := uuid.Parse("4aeab7c0-2963-42b3-b420-4007022aabee")
	token, _ := MakeJWT(userID, "abcdef", time.Hour)
	got, err := ValidateJWT(token, "abcdef")
	want := userID
	if err != nil {
		t.Errorf("%v", err)
	}
	if got != want {
		t.Errorf("got %v,\nwant %v", got, want)
	}
}

func TestGetBearerToken(t *testing.T) {
	userID, _ := uuid.Parse("4aeab7c0-2963-42b3-b420-4007022aabee")
	token, _ := MakeJWT(userID, "abcdef", time.Hour)
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)
	got, err := GetBearerToken(headers)
	want := token
	if err != nil {
		t.Errorf("%v", err)
	}
	if got != want {
		t.Errorf("got %v,\nwant %v", got, want)
	}

}
