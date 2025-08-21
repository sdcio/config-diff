package types

import (
	"fmt"

	dsTypes "github.com/sdcio/data-server/pkg/tree/types"
)

type Intent struct {
	Name     string
	Prio     int32
	BasePath string
	Flag     *dsTypes.UpdateInsertFlags
	Format   ConfigFormat
	Data     []byte
}

func NewIntent(name string, prio int32, flag *dsTypes.UpdateInsertFlags) *Intent {
	return &Intent{
		Name: name,
		Prio: prio,
		Flag: flag,
	}
}

func (i *Intent) GetName() string {
	return i.Name
}

func (i *Intent) GetPrio() int32 {
	return i.Prio
}

func (i *Intent) GetFlag() *dsTypes.UpdateInsertFlags {
	return i.Flag
}

func (i *Intent) GetData() []byte {
	return i.Data
}

func (i *Intent) GetBasePath() string {
	return i.BasePath
}

func (i *Intent) SetBasePath(p string) *Intent {
	i.BasePath = p
	return i
}

func (i *Intent) GetFormat() ConfigFormat {
	return i.Format
}

func (i *Intent) SetData(format ConfigFormat, data []byte) *Intent {
	i.Format = format
	i.Data = data
	return i
}

func (i *Intent) String() string {
	return fmt.Sprintf("Name: %s, Prio: %d, Flag: %s, Format: %s", i.GetName(), i.GetPrio(), i.GetFlag(), i.GetFormat())
}

func (i *Intent) Export() *IntentOutput {

	return &IntentOutput{
		Name:     i.GetName(),
		Priority: i.GetPrio(),
		// Flags:    NewFlagsOutput(i.GetFlag()),
	}
}

type Intents map[string]*Intent

func (i *Intents) AddIntent(ii *Intent) {
	if *i == nil {
		*i = make(Intents)
	}
	(*i)[ii.GetName()] = ii
}

func (i Intents) Export() []*IntentOutput {
	result := make([]*IntentOutput, 0, len(i))
	for _, intent := range i {
		result = append(result, intent.Export())
	}
	return result
}
