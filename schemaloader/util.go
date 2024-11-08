package schemaloader

import (
	"os"

	invv1alpha1 "github.com/sdcio/config-server/apis/inv/v1alpha1"
	"sigs.k8s.io/yaml"
)

func GetConfig(schemaConfigPath string) (*invv1alpha1.Schema, error) {
	b, err := os.ReadFile(schemaConfigPath)
	if err != nil {
		panic(err)
	}

	schema := &invv1alpha1.Schema{}
	if err := yaml.Unmarshal(b, schema); err != nil {
		return nil, err
	}
	return schema, nil
}
