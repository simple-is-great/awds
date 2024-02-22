package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	deviceIDPrefix string = "dev"
)

// Device represents a device, holding all necessary info. about device
// may divide into ip, port later
type Device struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Endpoint 	string	  `json:"endpoint"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func ValidateDeviceID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty device id")
	}

	prefix := fmt.Sprintf("%s_", deviceIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid device id - %s", id)
	}
	return nil
}

// NewDeviceID creates a new Device ID
func NewDeviceID() string {
	return fmt.Sprintf("%s_%s", deviceIDPrefix, xid.New().String())
}