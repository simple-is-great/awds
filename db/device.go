package db

import (
	"awds/types"

	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListDevices() ([]types.Device, error) {
	devices := []types.Device{}
	result := adapter.db.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func (adapter *DBAdapter) GetDevice(deviceID string) (types.Device, error) {
	var device types.Device
	result := adapter.db.Where("id = ?", deviceID).First(&device)
	if result.Error != nil {
		return device, result.Error
	}

	return device, nil
}

func (adapter *DBAdapter) InsertDevice(device *types.Device) error {
	result := adapter.db.Create(device)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert a device")
	}

	return nil
}

func (adapter *DBAdapter) UpdateDeviceEndpoint(deviceID string, endpoint string) error {
	var record types.Device
	result := adapter.db.Where("id = ?", deviceID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Endpoint = endpoint

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateDeviceDescription(deviceID string, description string) error {
	var record types.Device
	result := adapter.db.Where("id = ?", deviceID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Description = description

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) DeleteDevice(deviceID string) error {
	var device types.Device
	result := adapter.db.Where("id = ?", deviceID).Delete(&device)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to delete a device")
	}

	return nil
}
