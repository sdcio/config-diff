package types

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sdcio/config-diff/cmd/interfaces"
)

type IntentOutput struct {
	Name     string `json:"target"`
	Priority int32  `json:"priority"`
	// Flags    *FlagsOutput `json:"flag"`
}

var _ interfaces.Output = (*IntentOutput)(nil)

func (i *IntentOutput) ToString() string {
	return fmt.Sprintf("Name: %s, Priority: %d", i.Name, i.Priority)
}

func (i *IntentOutput) ToStringDetails() string {
	return i.ToString()
}

func (i *IntentOutput) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(i)
}
