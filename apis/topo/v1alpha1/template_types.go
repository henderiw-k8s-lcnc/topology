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

	targetv1 "github.com/henderiw-k8s-lcnc/target/apis/target/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	//targetv1alpha1pb "github.com/yndd/topology/gen/go/apis/topo/v1alpha1"
)

// TBD
// access parameters
/*
  	-> Multi-homing (in tier template for leafs and borderleafs)
	-> Access-speed
	-> Breakout
*/
// uplink parametsr and other
/*
	-> uplink-speed per layer
	-> Racks -> should be a user defined label
*/
// infra-underlay
/*
	-> DSCP for various protocols
	-> Protocols: unicast, multicast
	-> NTP/DNS global
	-> Mgmt IP per device
	-> serial number/mac address (configmap rendered by node reconciler)
*/

type SupportServers struct {
	DnsServers []*string `json:"dnsServers,omitempty"`
	NtPServers []*string `json:"ntpServers,omitempty"`
}

type TemplateSubnet struct {
	IPSubnet       string `json:"ipSubnet,omitempty"`
	SupportServers `json:"inline,omitempty"`
}

type FabricTemplate struct {
	// superspine
	Tier1      *TierTemplate           `json:"tier1,omitempty"`
	BorderLeaf *TierTemplate           `json:"borderLeaf,omitempty"`
	Pod        []*PodTemplate          `json:"pod,omitempty"`
	Settings   *FabricTemplateSettings `json:"settings,omitempty"`
	Tag        map[string]string       `json:"tag,omitempty"`
}

type FabricTemplateSettings struct {
	// max number of uplink per node to the next tier
	// default should be 1 and max is 4
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=4
	// +kubebuilder:default=1
	MaxUplinksTier2ToTier1 uint32 `json:"maxUplinksTier2ToTier1,omitempty"`
	// max number of uplink per node to the next tier
	// default should be 1 and max is 4
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=4
	// +kubebuilder:default=1
	MaxUplinksTier3ToTier2 uint32 `json:"maxUplinksTier3ToTier2,omitempty"`
	// max spines per pod is relevant in a border leaf configuration
	// it ensures the indexes are not impacted as long as the spines per pod
	// are below this limit and if this limit is not changed once the template
	// is instantiated
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=8
	// +kubebuilder:default=2
	MaxSpinesPerPod uint32 `json:"maxSpinesPerPod,omitempty"`
}

type PodTemplate struct {
	// number of pods defined based on this template
	// no default since templates should not define the pod number
	// default should be 1 and max is 16
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=16
	PodNumber *uint32 `json:"num,omitempty"`
	// Tier2 template, that defines the spine parameters in the pod definition
	Tier2 *TierTemplate `json:"tier2,omitempty"`
	// Tier3 template, that defines the leaf parameters in the pod definition
	Tier3 *TierTemplate `json:"tier3,omitempty"`
	// template reference to a template that defines the pod definition
	TemplateRef *TemplateReference `json:"templateRef,omitempty"`
	// definition reference to a template that defines the pod definition
	//DefinitionReference *string           `json:"definitionRef,omitempty"`
	Tag map[string]string `json:"tag,omitempty"`
}

type TierTemplate struct {
	// list to support multiple vendors in a tier - typically criss-cross
	VendorInfo []*FabricTierVendorInfo `json:"vendorInfo,omitempty"`
	// number of nodes in the tier
	// for superspine it is the number of spines in a spine plane
	NodeNumber uint32 `json:"num,omitempty"`
	// number of uplink per node to the next tier
	// default should be 1 and max is 4
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=4
	UplinksPerNode uint32 `json:"uplinkPerNode,omitempty"`
	// +kubebuilder:default=0
	MultiHoming uint32            `json:"multiHoming,omitempty"`
	Tag         map[string]string `json:"tag,omitempty"`
}

type FabricTierVendorInfo struct {
	Platform         string              `json:"platform,omitempty"`
	VendorType       targetv1.VendorType `json:"vendorType,omitempty"`
	InitialSwVersion string              `json:"initialSwVersion,omitempty"`
	Breakout         string              `json:"breakout,omitempty"`
	UpperTierlinks   Range               `json:"upperTierLinks,omitempty"`
	LowerTierlinks   Range               `json:"lowerTierLinks,omitempty"`
	Tag              map[string]string   `json:"tag,omitempty"`
}

type Range struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// TemplateProperties define the properties of the Template
type TemplateProperties struct {
	SupportServers `json:"inline,omitempty"`
	Subnet         *TemplateSubnet `json:"subnet,omitempty"`
	Fabric         *FabricTemplate `json:"fabric,omitempty"`
}

// TemplateSpec struct
type TemplateSpec struct {
	// Properties define the properties of the Template
	Properties TemplateProperties `json:"properties,omitempty"`
}

// A TemplateStatus represents the observed state of a Template.
type TemplateStatus struct {
	TopologyName         string `json:"topology-name,omitempty"`
	//Topology                *NddrTopologyTopology `json:"topology,omitempty"`
}

// +kubebuilder:object:root=true

// Template is the Schema for the Template API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availabilityZone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.oda[?(@.key=='resourceName')].value"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateSpec   `json:"spec,omitempty"`
	Status TemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TemplateList contains a list of Templates
type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}

// Template type metadata.
var (
	TemplateKind             = reflect.TypeOf(Template{}).Name()
	TemplateGroupKind        = schema.GroupKind{Group: Group, Kind: TemplateKind}.String()
	TemplateKindAPIVersion   = TemplateKind + "." + GroupVersion.String()
	TemplateGroupVersionKind = GroupVersion.WithKind(TemplateKind)
)
