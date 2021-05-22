package api

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDefaultDeployment() v1.Deployment {
	return v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1.DeploymentSpec{
			Replicas: nil,
			Selector: &metav1.LabelSelector{},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "",
					Annotations: nil,
				},
				Spec: v12.PodSpec{
					Volumes: nil,
				},
			},
		},
	}
}

func CreateContainers(app Application) []v12.Container {
	var envVars []v12.EnvVar
	for key, value := range app.Environment {
		envVars = append(envVars, v12.EnvVar{Name: key, Value: value})
	}

	volumeMounts := make([]v12.VolumeMount, len(app.Volumes))
	for index, volume := range app.Volumes {
		for path := range volume {
			volumeMounts[index] = v12.VolumeMount{
				Name:      CreatePVCName(app, path),
				MountPath: path,
			}
		}
	}

	containers := []v12.Container{{
		Name:         app.Name,
		Image:        fmt.Sprintf("%s:%s", app.Image, app.Version),
		Env:          envVars,
		VolumeMounts: volumeMounts,
	}}

	return containers
}

func CreateVolumes(app Application) []v12.Volume {
	volumes := make([]v12.Volume, len(app.Volumes))

	for index, volume := range app.Volumes {
		for path := range volume {
			volumes[index] = v12.Volume{
				Name: CreatePVCName(app, path),
				VolumeSource: v12.VolumeSource{
					PersistentVolumeClaim: &v12.PersistentVolumeClaimVolumeSource{
						ClaimName: CreatePVCName(app, path),
					},
				},
			}

			break
		}
	}

	return volumes
}

func CreateDeployment(app Application) (v1.Deployment, error) {
	deployment := generateDefaultDeployment()

	deployment.ObjectMeta.Name = app.Name
	deployment.ObjectMeta.Namespace = app.Namespace

	if app.Replicas == 0 {
		app.Replicas = 1
	}
	deployment.Spec.Replicas = &app.Replicas

	deployment.Spec.Selector.MatchLabels = map[string]string{
		"app": app.Name,
	}

	if app.ImagePullSecret != "" {
		deployment.Spec.Template.Spec.ImagePullSecrets = []v12.LocalObjectReference{
			{Name: app.ImagePullSecret},
		}
	}

	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{
		"app": app.Name,
	}

	deployment.Spec.Template.Spec.Volumes = CreateVolumes(app)
	deployment.Spec.Template.Spec.Containers = CreateContainers(app)

	return deployment, nil
}
