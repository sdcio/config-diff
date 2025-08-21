package types

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type TargetOutput struct {
	TargetName string          `json:"target"`
	TargetPath string          `json:"target-path"`
	Schema     *SchemaOutput   `json:"schema"`
	Intents    []*IntentOutput `json:"intents"`
}

func (t *TargetOutput) ToString() string {
	return fmt.Sprintf("%s [ %s ]\n", t.TargetName, t.Schema.ToString())
}

func (t *TargetOutput) ToStringDetails() string {
	indent := "  "
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "Target: %s ( %s )\n", t.TargetName, t.TargetPath)
	if t.Schema != nil {
		fmt.Fprintf(sb, "%sSchema:\n", indent)
		fmt.Fprintf(sb, "%[1]s%[1]s%[2]s\n", indent, t.Schema.ToStringDetails())
	}
	fmt.Fprintf(sb, "%sIntents:\n", indent)
	for _, i := range t.Intents {
		fmt.Fprintf(sb, "%[1]s%[1]s%[2]s\n", indent, i.ToStringDetails())
	}
	return sb.String()
}

func (t *TargetOutput) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(t)
}

type TargetOutputSlice []*TargetOutput

func (t TargetOutputSlice) ToString() string {
	sb := &strings.Builder{}
	for _, target := range t {
		fmt.Fprint(sb, target.ToString())
	}

	return sb.String()
}

func (t TargetOutputSlice) ToStringDetails() string {
	sb := &strings.Builder{}
	for _, target := range t {
		fmt.Fprint(sb, target.ToStringDetails())
	}

	return sb.String()
}

func (t TargetOutputSlice) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(t)
}
