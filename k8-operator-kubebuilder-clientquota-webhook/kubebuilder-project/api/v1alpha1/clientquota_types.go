package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientSpec struct {
	Name         string `json:"name"`
	APIKey       string `json:"apiKey"`
	QuotaMinutes int    `json:"quotaMinutes"`
}

type ClientStatus struct {
	Name             string `json:"name"`
	RemainingMinutes int    `json:"remainingMinutes"`
}

// ClientQuotaSpec defines the desired state of ClientQuota
type ClientQuotaSpec struct {
	Namespace string       `json:"namespace"` // namespace to monitor
	Clients   []ClientSpec `json:"clients"`
}

// ClientQuotaStatus defines the observed state of ClientQuota
type ClientQuotaStatus struct {
	Clients []ClientStatus `json:"clients,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClientQuota is the Schema for the clientquotas API
type ClientQuota struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClientQuotaSpec   `json:"spec,omitempty"`
	Status ClientQuotaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClientQuotaList contains a list of ClientQuota
type ClientQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClientQuota `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClientQuota{}, &ClientQuotaList{})
}
