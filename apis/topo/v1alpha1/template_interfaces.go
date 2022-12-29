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
	"fmt"
	"strconv"
)

func (x *Template) GetNumPods() uint32 {
	if x.Spec.Properties.Fabric.Pod == nil {
		return 0
	}
	numPod := uint32(0)
	for _, p := range x.Spec.Properties.Fabric.Pod {
		numPod += *p.PodNumber
	}
	return numPod
}

func (x *FabricTemplate) CheckTemplate(master bool) error {
	if x.Tier1 != nil && x.BorderLeaf != nil {
		return fmt.Errorf("a template cannot have a mix of tier1 and boarderleafs")
	}
	if !master && (x.Tier1 != nil || x.BorderLeaf != nil) {
		return fmt.Errorf("a child template cannot have tier1 or borderleafs defined")
	}
	if x.Pod == nil {
		return nil
	}
	// for a non master template we expect only a single pod definition
	if !master && len(x.Pod) != 1 {
		return fmt.Errorf("a child template can only have 1 pod defined")
	}
	for _, p := range x.Pod {
		if err := p.CheckPodTemplate(master); err != nil {
			return err
		}
	}
	return nil
}

func (x *FabricTemplate) HasDefinitionRef() bool {
	if x.Pod == nil {
		return false
	}
	for _, p := range x.Pod {
		if p.HasDefinitionRef() {
			return true
		}
	}
	return false
}

func (x *FabricTemplate) HasTemplateRef() bool {
	if x.Pod == nil {
		return false
	}
	for _, p := range x.Pod {
		if p.HasTemplateRef() {
			return true
		}
	}
	return false
}

func (x *FabricTemplate) HasReference() bool {
	return x.HasDefinitionRef() || x.HasTemplateRef()
}

func (x *FabricTemplate) HasPodDefinition() bool {
	if x.Pod == nil {
		return false
	}
	for _, p := range x.Pod {
		if p.Tier2 != nil {
			return true
		}
		if p.Tier3 != nil {
			return true
		}
	}
	return false
}

func (x *FabricTemplate) HasTier1() bool      { return x.Tier1 == nil }
func (x *FabricTemplate) HasBorderLeaf() bool { return x.BorderLeaf == nil }

// getSuperSPines identifies the max number of spines in a pod
func (x *FabricTemplate) GetSuperSpines() uint32 {
	var superspines uint32
	for _, podInfo := range x.Pod {
		if superspines < podInfo.Tier2.NodeNumber {
			superspines = podInfo.Tier2.NodeNumber
		}
	}
	return superspines
}

func (x *PodTemplate) CheckPodTemplate(master bool) error {
	// check mix of native definition
	if x.Tier2 != nil || x.Tier3 != nil {
		if x.TemplateRef != nil {
			// this i not allowed
			return fmt.Errorf("podTemplate error: native pod definition can not be mixed with template references")
		}
		//if x.TemplateReference != nil || x.DefinitionReference != nil {
		//	// this i not allowed
		//	return fmt.Errorf("podTemplate error: native pod definition can not be mixed with template/definition references")
		//}
	}
	if master {
		// master template
		if x.HasReference() && x.PodNumber != nil {
			return fmt.Errorf("a template with a reference cannot define the pod number")
		}
		if !x.HasReference() && x.PodNumber == nil {
			return fmt.Errorf("a pod template w/o references should have a podNumber defined")
		}
	} else {
		// this is a child template
		if x.HasReference() {
			return fmt.Errorf("a child template cannot have another child template")
		}
		if x.PodNumber != nil && *x.PodNumber != 1 {
			return fmt.Errorf("a child reference can only define 1 pod instance")
		}
	}
	return nil
}

func (x *PodTemplate) HasReference() bool {
	return x.HasTemplateRef() || x.HasDefinitionRef()
}

func (x *PodTemplate) HasTemplateRef() bool {
	return x.TemplateRef != nil
}

func (x *PodTemplate) HasDefinitionRef() bool {
	return x.TemplateRef != nil
}

func (x *PodTemplate) GetPodNumber() uint32 {
	if x.PodNumber == nil {
		return 1
	}
	return *x.PodNumber
}

func (x *PodTemplate) IsToBeDeployed() bool {
	if _, ok := x.Tag["toBeDeployed"]; !ok {
		return false
	}
	b, err := strconv.ParseBool(x.Tag["toBeDeployed"])
	if err != nil {
		return false
	}
	return b
}

func (x *PodTemplate) SetToBeDeployed(b bool) {
	if x.Tag == nil {
		x.Tag = map[string]string{}
	}
	x.Tag["toBeDeployed"] = strconv.FormatBool(b)
}
