package types

import (
	"fmt"
	"io"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"google.golang.org/protobuf/encoding/protojson"
)

type BlameResultOutput struct {
	*sdcpb.BlameTreeElement
}

func NewBlameResultOutput(e *sdcpb.BlameTreeElement) *BlameResultOutput {
	return &BlameResultOutput{
		BlameTreeElement: e,
	}
}

func (b *BlameResultOutput) ToString() string {
	return fmt.Sprintln(b.BlameTreeElement.ToString())
}

func (b *BlameResultOutput) ToStringDetails() string {
	return b.ToString()
}

func (b *BlameResultOutput) WriteToJson(w io.Writer) error {
	marshaler := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	jsonBytes, err := marshaler.Marshal(b.BlameTreeElement)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}
