// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vyos_models

import (
	"encoding/json"
	"errors"
)

var (
	// ErrUnsupportedType is returned if the type is not implemented.
	ErrUnsupportedType = errors.New("unsupported type")

	_ json.Unmarshaler = &MultiValuedString{}
)

type MultiValuedString []string

func (mvs *MultiValuedString) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}

	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}
	switch obj := jsonObj.(type) {
	case string:
		*mvs = MultiValuedString([]string{obj})
		return nil
	case []interface{}:
		s := make([]string, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return ErrUnsupportedType
			}
			s = append(s, value)
		}
		*mvs = MultiValuedString(s)
		return nil
	}
	return ErrUnsupportedType
}
