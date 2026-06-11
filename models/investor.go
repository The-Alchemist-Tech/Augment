package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Investor struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedOn    time.Time `json:"created_on"`
	LastModified time.Time `json:"last_modified"`
}

func (i *Investor) String() string {
	return fmt.Sprintf("Investor[ID=%d, Username=%s, Email=%s, FirstName=%s, LastName=%s, CreatedOn=%s, LastModified=%s]",
		i.ID, i.Username, i.Email, i.FirstName, i.LastName, i.CreatedOn.Format(time.RFC3339), i.LastModified.Format(time.RFC3339))
}

func (i *Investor) Validate() error {
	if strings.TrimSpace(i.Username) == "" {
		return fmt.Errorf("Username is required")
	}
	if strings.TrimSpace(i.Email) == "" {
		return fmt.Errorf("Email is required")
	}
	if _, err := mail.ParseAddress(i.Email); err != nil {
		return fmt.Errorf("Email is invalid")
	}
	if strings.TrimSpace(i.FirstName) == "" {
		return fmt.Errorf("FirstName is required")
	}
	if strings.TrimSpace(i.LastName) == "" {
		return fmt.Errorf("LastName is required")
	}
	return nil
}