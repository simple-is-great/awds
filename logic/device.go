package logic

import (
	"awds/types"

	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListDevices() ([]types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListDevices",
	})

	logger.Debug("received ListDevices()")

	return logic.dbAdapter.ListDevices()
}

func (logic *Logic) GetDevice(deviceID string) (types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetDevice",
	})

	logger.Debug("received GetDevice()")

	return logic.dbAdapter.GetDevice(deviceID)
}

func (logic *Logic) CreateDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateDevice",
	})

	logger.Debug("received CreateDevice()")

	return logic.dbAdapter.InsertDevice(device)
}

func (logic *Logic) UpdateDeviceEndpoint(deviceID string, endpoint string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIP",
	})

	logger.Debug("received UpdateDeviceIP()")

	return logic.dbAdapter.UpdateDeviceEndpoint(deviceID, endpoint)
}

func (logic *Logic) UpdateDeviceDescription(deviceID string, description string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceDescription",
	})

	logger.Debug("received UpdateDeviceDescription()")

	return logic.dbAdapter.UpdateDeviceDescription(deviceID, description)
}

func (logic *Logic) DeleteDevice(deviceID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteDevice",
	})

	logger.Debug("received DeleteDevice()")

	device, err := logic.GetDevice(deviceID)

	if err != nil {
		log.Error("%s does not exist, cannot delete device", device)
		return err
	}
	
	return logic.dbAdapter.DeleteDevice(deviceID)
}
