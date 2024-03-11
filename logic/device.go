package logic

import (
	"awds/types"
	"fmt"

	"github.com/go-resty/resty/v2"
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

func (logic *Logic) GetDeviceInfo(device *types.Device) (*types.Device, error) {
	var response map[string]interface{}
	requestAddr := fmt.Sprintf("http://%s:%s/computing_measure", device.IP, device.Port)
	fmt.Println("requestAddr", requestAddr)
	client := resty.New()
	_, err := client.R().SetResult(&response).Get(requestAddr)
	if err != nil {
		return nil, err
	}

	// save info
	if response["network_latency"].(float64) <= 0 {
		return nil, fmt.Errorf("Network unavailable, value must be positive!")
	}
	device.NetworkLatency = response["network_latency"].(float64)

	if response["cpu"].(float64) <= 0 { 
		return nil, fmt.Errorf("CPU unavailable, value must be postive!")
	}
	device.CPU = response["cpu"].(float64)

	if response["memory"].(float64) <= 0 {
		return nil, fmt.Errorf("Memory unavailable, value must be postive!")
	}
	device.Memory = response["memory"].(float64) // TODO: need to fix key, ram -> memory

	fmt.Println("after getting info", device)

	return device, nil
}

func (logic *Logic) CreateDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateDevice",
	})

	logger.Debug("received CreateDevice()")

	// get device info
	device, err := logic.GetDeviceInfo(device)
	if err != nil {
		return err
	}
	
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
