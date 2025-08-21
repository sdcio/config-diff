package cmd

import (
	"fmt"
	"os"

	"github.com/sdcio/config-diff/cmd/interfaces"
)

func WriteOutput(o interfaces.Output) (err error) {
	switch {
	case jsonOutput:
		err = o.WriteToJson(os.Stdout)
	case detailed:
		_, err = fmt.Print(o.ToStringDetails())
	default:
		_, err = fmt.Print(o.ToString())
	}
	return err
}
