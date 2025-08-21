package types

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/sdcio/data-server/pkg/tree/types"
)

type ValidationStatsOutput struct {
	Target string `json:"target"`
	Passed bool   `json:"passed"`
	*types.ValidationStatOverall
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

func NewValidationStatsOutput(targetName string, vr types.ValidationResults, vs *types.ValidationStatOverall) *ValidationStatsOutput {
	jsonResult := &ValidationStatsOutput{
		ValidationStatOverall: vs,
		Target:                targetName,
		Passed:                !vr.HasErrors(),
		Errors:                vr.ErrorsStr(),
		Warnings:              vr.WarningsStr(),
	}
	return jsonResult
}

func (v *ValidationStatsOutput) ToString() string {
	sb := &strings.Builder{}

	fmt.Fprintln(sb, "Validations performed:")
	// sort the map, by getting the keys first
	keys := make([]string, 0, len(v.GetCounter()))
	for typ := range v.GetCounter() {
		keys = append(keys, typ)
	}

	// sorting the keys
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	indent := "  "
	// printing the stats in the sorted order
	for _, typ := range keys {
		fmt.Fprintf(sb, "%s%s: %d\n", indent, typ, v.GetCounter()[typ])
	}

	if len(v.Errors) == 0 {
		fmt.Fprintln(sb, "Successfully Validated!")
	}

	if len(v.Errors) > 0 {
		fmt.Fprintln(sb, "Errors:")
		for _, errStr := range v.Errors {
			fmt.Fprintln(sb, errStr)
		}
	}

	if len(v.Warnings) > 0 {
		fmt.Fprintln(sb, "Warnings:")
		for _, warnStr := range v.Warnings {
			fmt.Fprintln(sb, warnStr)
		}
	}

	return sb.String()
}

func (v *ValidationStatsOutput) ToStringDetails() string {
	return v.ToString()
}

func (v *ValidationStatsOutput) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(v)
}
