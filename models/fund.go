package models

import (
	"fmt"
	"strings"
	"time"
)

type Fund struct {
    ID           int       `json:"id"`
    Name         string    `json:"name"`
    Units        *int      `json:"units"`
    CreatedOn    time.Time `json:"created_on"`
    LastModified time.Time `json:"last_modified"`
}

func (f *Fund) String() string {
	units := "<nil>"
	if f.Units != nil {
		units = fmt.Sprintf("%d", *f.Units)
	}

	return fmt.Sprintf("Fund[ID=%d, Name=%s, Units=%s, CreatedOn=%s, LastModified=%s]",
		f.ID, f.Name, units, f.CreatedOn.Format(time.RFC3339), f.LastModified.Format(time.RFC3339))
}

func (f *Fund) Validate() error {
	if strings.TrimSpace(f.Name) == "" {
		return fmt.Errorf("Name is required.")
	}
	if f.Units == nil || *f.Units <= 0 {
		return fmt.Errorf("Units must be present and greater than zero.")
	}

	return nil
}