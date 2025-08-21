package types

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sdcio/config-diff/cmd/interfaces"
	dsTypes "github.com/sdcio/data-server/pkg/tree/types"
)

type FlagsOutput struct {
	New            bool `json:"new,omitempty"`
	Delete         bool `json:"delete,omitempty"`
	OnlyIntended   bool `json:"OnlyIntended,omitempty"`
	ExplicitDelete bool `json:"ExplicitDelete,omitempty"`
}

var _ interfaces.Output = (*FlagsOutput)(nil)

func NewFlagsOutput(f *dsTypes.UpdateInsertFlags) *FlagsOutput {
	return &FlagsOutput{
		New:            f.GetNewFlag(),
		Delete:         f.GetDeleteFlag(),
		OnlyIntended:   f.GetDeleteOnlyIntendedFlag(),
		ExplicitDelete: f.GetExplicitDeleteFlag(),
	}
}

func (f *FlagsOutput) ToString() string {
	sb := &strings.Builder{}
	fmt.Fprint(sb, "[")
	if f.New {
		fmt.Fprint(sb, "New ")
	}
	if f.Delete {
		fmt.Fprint(sb, "Delete ")
	}
	if f.OnlyIntended {
		fmt.Fprint(sb, "OnlyIntended ")
	}
	if f.ExplicitDelete {
		fmt.Fprint(sb, "ExplicitDelete ")
	}
	fmt.Fprint(sb, "]")
	return sb.String()
}

func (f *FlagsOutput) ToStringDetails() string {
	return f.ToString()
}

func (f *FlagsOutput) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(f)
}
