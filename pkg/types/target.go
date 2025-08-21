package types

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdcio/config-diff/pkg/utils"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"google.golang.org/protobuf/encoding/protojson"
)

type Target struct {
	config  TargetConfig
	schema  *sdcpb.Schema
	intents Intents
}

func NewTarget(w TargetConfig) *Target {
	wi := &Target{
		config: w,
	}
	err := wi.loadSchemaInfo()
	if err != nil {
		return wi
	}
	err = wi.loadIntent()
	if err != nil {
		return wi
	}

	return wi
}

func (t *Target) AddIntent(i *Intent) error {
	t.intents.AddIntent(i)
	return nil
}

func (t *Target) DeleteIntent(name string) error {
	_, exists := t.intents[name]
	if !exists {
		return fmt.Errorf("error deleting intent %s - not found", name)
	}
	delete(t.intents, name)
	err := os.Remove(t.config.ConfigFileName(name))
	if err != nil {
		return err
	}
	return nil
}

func (t *Target) GetIntents() Intents {
	return t.intents
}

func (t *Target) schemaPersist() error {
	schemaByte, err := protojson.Marshal(t.schema)
	if err != nil {
		return err
	}

	err = os.WriteFile(t.config.SchemaDefinitionFilePath(), schemaByte, 0644)
	return err
}

func (t *Target) intentsPersist() error {
	for name, content := range t.intents {
		data, err := json.Marshal(content)
		if err != nil {
			return err
		}
		err = os.WriteFile(t.config.ConfigFileName(name), data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Target) SetSchema(s *sdcpb.Schema) {
	t.schema = s
}

func (t *Target) Persist() error {
	err := utils.CreateFolder(t.config.TargetPath())
	if err != nil {
		return err
	}
	err = t.schemaPersist()
	if err != nil {
		return err
	}
	err = t.intentsPersist()
	if err != nil {
		return err
	}
	return nil
}

func (t *Target) loadIntent() error {
	matches, err := filepath.Glob(t.config.ConfigFileGlob())
	if err != nil {
		return err
	}
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		ii := &Intent{}
		err = json.Unmarshal(data, ii)
		if err != nil {
			return err
		}
		t.intents.AddIntent(ii)
	}
	return nil
}

func (t *Target) loadSchemaInfo() (err error) {
	t.schema, err = utils.SchemaLoadSdcpbSchemaFile(t.config.SchemaDefinitionFilePath())
	if err != nil {
		return err
	}
	return nil
}

func (t *Target) GetSchema() *sdcpb.Schema {
	return t.schema
}

func (t *Target) Export() *TargetOutput {
	result := &TargetOutput{
		TargetName: t.config.TargetName(),
		Intents:    t.intents.Export(),
		TargetPath: t.config.TargetPath(),
	}
	if t.schema != nil {
		result.Schema = &SchemaOutput{
			Vendor:  t.schema.Vendor,
			Version: t.schema.Version,
		}
	}

	return result
}

type Targets []*Target

func (t *Targets) Add(w *Target) {
	*t = append(*t, w)
}

func (t Targets) Export() TargetOutputSlice {
	result := make([]*TargetOutput, 0, len(t))
	for _, target := range t {
		result = append(result, target.Export())
	}
	return result
}

type TargetConfig interface {
	TargetName() string
	ConfigFileName(intentName string) string
	ConfigFileGlob() string
	SchemaDefinitionFilePath() string
	TargetPath() string
}
