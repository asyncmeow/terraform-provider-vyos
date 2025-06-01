// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vyos_models

type InterfacesEthernet struct {
	Addresses *MultiValuedString `json:"address"`
	HwId      string             `json:"hw-id"`
}
