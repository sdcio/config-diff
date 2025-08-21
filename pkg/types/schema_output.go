package types

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sdcio/config-diff/cmd/interfaces"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type SchemaOutput struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
}

func NewSchemaOutput(s *sdcpb.Schema) *SchemaOutput {
	return &SchemaOutput{
		Vendor:  s.GetVendor(),
		Version: s.GetVersion(),
	}

}

func (s *SchemaOutput) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(s)
}

func (s *SchemaOutput) ToString() string {
	return fmt.Sprintf("Vendor: %s, Version: %s", s.Vendor, s.Version)
}

func (s *SchemaOutput) ToStringDetails() string {
	return s.ToString()
}

type SchemaOutputSlice []*SchemaOutput

var _ interfaces.Output = SchemaOutputSlice{}

func NewSchemaOutputList(ss []*sdcpb.Schema) SchemaOutputSlice {
	result := SchemaOutputSlice{}
	for _, s := range ss {
		result = append(result, NewSchemaOutput(s))
	}
	return result
}

func (s SchemaOutputSlice) ToString() string {
	sb := &strings.Builder{}
	for _, schema := range s {
		fmt.Fprint(sb, schema.ToString())
	}
	return sb.String()
}

func (s SchemaOutputSlice) ToStringDetails() string {
	return s.ToString()
}

func (s SchemaOutputSlice) WriteToJson(w io.Writer) error {
	jEnc := json.NewEncoder(w)
	jEnc.SetIndent("", "  ")
	return jEnc.Encode(s)
}
