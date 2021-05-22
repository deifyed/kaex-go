package api

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDefaultPVC() v1.PersistentVolumeClaim {
	return v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1.PersistentVolumeClaimSpec{
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceRequestsStorage: resource.Quantity{
						Format: "1Gi",
					},
				},
			},
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
		},
	}
}

func CreatePVCName(app Application, path string) string {
	cleanPath := strings.Replace(path, "/", "", -1)

	return fmt.Sprintf("%s-%s", app.Name, cleanPath)
}

func CreatePersistentVolume(app Application, path string, size string) (v1.PersistentVolumeClaim, error) {
	volume := generateDefaultPVC()

	volume.ObjectMeta.Name = CreatePVCName(app, path)
	volume.ObjectMeta.Namespace = app.Namespace

	capacity, err := createStorageRequest(size)
	if err != nil {
		return v1.PersistentVolumeClaim{}, err
	}
	volume.Spec.Resources.Requests = capacity

	return volume, nil
}

func createStorageRequest(requestSize string) (v1.ResourceList, error) {
	quantity, err := resource.ParseQuantity("1Gi")
	if requestSize != "" {
		quantity, err = resource.ParseQuantity(requestSize)

		if err != nil {
			return nil, err
		}
	}

	return v1.ResourceList{
		v1.ResourceStorage: quantity,
	}, nil
}
