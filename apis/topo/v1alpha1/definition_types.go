/*
Copyright 2021 NDD.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DefinitionSpec struct
type DefinitionSpec struct {
	// Properties define the properties of the Definition
	Properties *DefinitionProperties `json:"properties,omitempty"`
}

// A DefinitionStatus represents the observed state of a Definition.
type DefinitionStatus struct {
}

// DefinitionProperties define the properties of the Definition
type DefinitionProperties struct {
	Location       *Location        `json:"location,omitempty"`
	Templates      []*TemplateRule  `json:"templates,omitempty"`
	DiscoveryRules []*DiscoveryRule `json:"discoveryRules,omitempty"`
}

type TemplateRule struct {
	TemplateRef *TemplateReference `json:"templateRef"`
	// +kubebuilder:default=false
	DigitalTwin bool `json:"digitalTwin,omitempty"`
}

// TemplateReference is used to reference a particular template.
type TemplateReference struct {
	// Namespace is the namespace of the template
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name is the name of the template
	Name string `json:"name"`
}

type DiscoveryRule struct {
	DiscoveryRuleRef *DiscoveryRuleReference `json:"discoveryRuleRef"`
	// +kubebuilder:default=false
	DigitalTwin bool `json:"digitalTwin,omitempty"`
}

// DiscoveryRuleReference is used to reference a particular discovery rule.
type DiscoveryRuleReference struct {
	// Namespace is the namespace of the discovery rule
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name is the name of the discovery rule
	Name string `json:"name"`
}

// +kubebuilder:object:root=true

// Definition is the Schema for the Topology API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availability-zone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.topology-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type Definition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DefinitionSpec   `json:"spec,omitempty"`
	Status DefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DefinitionList contains a list of Definitions
type DefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Definition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Definition{}, &DefinitionList{})
}

// Definition type metadata.
var (
	DefinitionKind             = reflect.TypeOf(Definition{}).Name()
	DefinitionGroupKind        = schema.GroupKind{Group: Group, Kind: DefinitionKind}.String()
	DefinitionKindAPIVersion   = DefinitionKind + "." + GroupVersion.String()
	DefinitionGroupVersionKind = GroupVersion.WithKind(DefinitionKind)
)
