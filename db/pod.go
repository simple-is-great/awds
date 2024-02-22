package db

import (
	"awds/types"

	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListPods() ([]types.Pod, error) {
	pods := []types.Pod{}
	result := adapter.db.Find(&pods)
	if result.Error != nil {
		return nil, result.Error
	}

	return pods, nil
}

func (adapter *DBAdapter) GetPod(podID string) (types.Pod, error) {
	var pod types.Pod
	result := adapter.db.Where("id = ?", podID).First(&pod)
	if result.Error != nil {
		return pod, result.Error
	}

	return pod, nil
}

func (adapter *DBAdapter) InsertPod(pod *types.Pod) error {
	result := adapter.db.Create(pod)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert a pod")
	}

	return nil
}

func (adapter *DBAdapter) UpdatePodEndpoint(podID string, endpoint string) error {
	var record types.Pod
	result := adapter.db.Where("id = ?", podID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Endpoint = endpoint

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdatePodDescription(podID string, description string) error {
	var record types.Pod
	result := adapter.db.Where("id = ?", podID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Description = description

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) DeletePod(podID string) error {
	var pod types.Pod
	result := adapter.db.Where("id = ?", podID).Delete(&pod)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to delete a pod")
	}

	return nil
}
