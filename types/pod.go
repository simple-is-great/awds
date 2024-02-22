package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	podIDPrefix string = "pod"
)

// Pod represents an pod, holding all necessary info. to run awdp
type Pod struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Endpoint 	string 	  `json:"endpoint"`
	Description string 	  `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func ValidatePodID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty pod id")
	}

	prefix := fmt.Sprintf("%s_", podIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid pod id - %s", id)
	}
	return nil
}

// NewPodID creates a new Pod ID
func NewPodID() string {
	return fmt.Sprintf("%s_%s", podIDPrefix, xid.New().String())
}