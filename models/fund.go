package models

import "time"

type Fund struct {
    ID           int       `json:"id"`
    Name         string    `json:"name"`
    Units        int       `json:"units"`
    CreatedOn    time.Time `json:"created_on"`
    LastModified time.Time `json:"last_modified"`
}