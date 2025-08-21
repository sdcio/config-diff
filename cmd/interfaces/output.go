package interfaces

import "io"

type Output interface {
	ToString() string
	ToStringDetails() string
	WriteToJson(w io.Writer) error
}
