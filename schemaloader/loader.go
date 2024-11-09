package schemaloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdcio/config-server/apis/inv/v1alpha1"
	loader "github.com/sdcio/config-server/pkg/schema"
	"github.com/sdcio/schema-server/pkg/store"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type Config struct {
	TmpPath     string
	SchemasPath string
}

func New(schemastore store.Store, cfg *Config) (*SchemaLoader, error) {
	if err := os.MkdirAll(cfg.TmpPath, 0755|os.ModeDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cfg.SchemasPath, 0755|os.ModeDir); err != nil {
		return nil, err
	}

	return &SchemaLoader{
		schemastore: schemastore,
		cfg:         cfg,
	}, nil
}

type SchemaLoader struct {
	schemastore store.Store
	cfg         *Config
}

func (r *SchemaLoader) LoadSchema(ctx context.Context, schemacr *v1alpha1.Schema) (*sdcpb.CreateSchemaResponse, error) {
	schemaLoader, err := loader.NewLoader(
		filepath.Join(r.cfg.TmpPath),
		filepath.Join(r.cfg.SchemasPath),
		NewNopResolver(),
	)
	if err != nil {
		return nil, err
	}

	schemaLoader.AddRef(ctx, schemacr)
	_, dirExists, err := schemaLoader.GetRef(ctx, schemacr.Spec.GetKey())
	if err != nil {
		return nil, err
	}
	if !dirExists {
		fmt.Println("loading...")
		if err := schemaLoader.Load(ctx, schemacr.Spec.GetKey()); err != nil {
			return nil, err
		}
	}

	return r.schemastore.CreateSchema(ctx, &sdcpb.CreateSchemaRequest{
		Schema: &sdcpb.Schema{
			Name:    "",
			Vendor:  schemacr.Spec.Provider,
			Version: schemacr.Spec.Version,
		},
		File:      schemacr.Spec.GetNewSchemaBase(r.cfg.SchemasPath).Models,
		Directory: schemacr.Spec.GetNewSchemaBase(r.cfg.SchemasPath).Includes,
		Exclude:   schemacr.Spec.GetNewSchemaBase(r.cfg.SchemasPath).Excludes,
	})
}
