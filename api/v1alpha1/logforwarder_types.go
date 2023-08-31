/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/samanthajayasinghe/multi-cluster-operator/pkg/clusterx"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LogForwarderSpec defines the desired state of LogForwarder
type LogForwarderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Image       string `json:"image,omitempty"`
	Destination string `json:"destination,omitempty"`
}

// LogForwarderStatus defines the observed state of LogForwarder
type LogForwarderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase string `json:"phase"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LogForwarder is the Schema for the logforwarders API
type LogForwarder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LogForwarderSpec   `json:"spec,omitempty"`
	Status LogForwarderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LogForwarderList contains a list of LogForwarder
type LogForwarderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LogForwarder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LogForwarder{}, &LogForwarderList{})
}

// Job ...
func (t *LogForwarder) Job() batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      t.Name,
			Namespace: t.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{
						{
							Name:  t.Name,
							Image: t.Spec.Image,
						},
					},
				},
			},
		},
	}

	clusterx.SetOwner(t, job)

	return *job
}
