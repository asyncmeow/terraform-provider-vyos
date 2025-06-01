package vyos_models

type InterfacesEthernet struct {
	Addresses *MultiValuedString `json:"address"`
	HwId      string             `json:"hw-id"`
}
