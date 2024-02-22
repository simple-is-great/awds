package logic

import (
	"awds/types"

	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListPods() ([]types.Pod, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListPods",
	})

	logger.Debug("received ListPods()")

	return logic.dbAdapter.ListPods()
}

func (logic *Logic) GetPod(podID string) (types.Pod, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetPod",
	})

	logger.Debug("received GetPod()")

	return logic.dbAdapter.GetPod(podID)
}

func (logic *Logic) RegisterPod(pod *types.Pod) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "RegisterPod",
	})

	logger.Debug("received RegisterPod()")

	return logic.dbAdapter.InsertPod(pod)
}

func (logic *Logic) UpdatePodEndpoint(podID string, endpoint string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePodEndpoint",
	})

	logger.Debug("received UpdatePodEndpoint()")

	return logic.dbAdapter.UpdatePodEndpoint(podID, endpoint)
}

func (logic *Logic) UpdatePodDescription(podID string, description string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePodDescription",
	})

	logger.Debug("received UpdatePodDescription()")

	return logic.dbAdapter.UpdatePodDescription(podID, description)
}

func (logic *Logic) DeletePod(podID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeletePod",
	})

	logger.Debug("received DeletePod()")

	pod, err := logic.GetPod(podID)

	if err != nil {
		log.Error("%s does not exist, cannot delete pod", pod)
		return err
	}

	return logic.dbAdapter.DeletePod(podID)
}
