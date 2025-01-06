package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sdcio/config-diff/schemaclient"
	"github.com/sdcio/config-diff/schemaloader"
	log "github.com/sirupsen/logrus"

	"github.com/sdcio/data-server/pkg/tree"
	treejson "github.com/sdcio/data-server/pkg/tree/importer/json"
	"github.com/sdcio/schema-server/pkg/config"
	"github.com/sdcio/schema-server/pkg/store/persiststore"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

func main() {
	ctx := context.Background()

	args := os.Args
	if len(args) < 1 {
		panic("cannot execute need config and base dir")
	}

	schemastore, err := persiststore.New(ctx, "foo", &config.SchemaPersistStoreCacheConfig{})
	if err != nil {
		panic(err)
	}

	schemacr, err := schemaloader.GetConfig(args[1])
	if err != nil {
		panic(err)
	}

	_, err = schemastore.GetSchemaDetails(ctx, &sdcpb.GetSchemaDetailsRequest{
		Schema: &sdcpb.Schema{
			Vendor:  schemacr.Spec.Provider,
			Version: schemacr.Spec.Version,
		},
	})
	if err != nil {
		schemaLoader, err := schemaloader.New(schemastore, &schemaloader.Config{TmpPath: "tmp/tmp", SchemasPath: "tmp/schemas"})
		if err != nil {
			panic(err)
		}
		_, err = schemaLoader.LoadSchema(ctx, schemacr)
		if err != nil {
			panic(err)
		}
	}
	scb := schemaclient.NewMemSchemaClientBound(schemastore, &sdcpb.Schema{
		Vendor:  schemacr.Spec.Provider,
		Version: schemacr.Spec.Version,
	})

	tc := tree.NewTreeContext(tree.NewTreeSchemaCacheClient("dev1", nil, scb), "test")
	root, err := tree.NewTreeRoot(ctx, tc)
	if err != nil {
		panic(err)
	}
	// fmt.Println(root.String())

	//
	// Load running
	//
	jsonBytes, err := os.ReadFile("/home/mava/projects/config-diff/data/config/running/running_eos_02.json")
	if err != nil {
		panic(err)
	}

	var j any
	err = json.Unmarshal(jsonBytes, &j)
	if err != nil {
		panic(err)
	}
	jti := treejson.NewJsonTreeImporter(j)
	err = root.ImportConfig(ctx, jti, tree.RunningIntentName, tree.RunningValuesPrio)
	if err != nil {
		panic(err)
	}

	root.FinishInsertionPhase()
	fmt.Println(root.String())

	// //
	// // Here we load the indent that is supposed to go on top of the running
	// //
	// jsonBytes, err = os.ReadFile("data/config/additions/srl_01.json")
	// if err != nil {
	// 	panic(err)
	// }

	// err = json.Unmarshal(jsonBytes, &j)
	// if err != nil {
	// 	panic(err)
	// }
	// jti = treejson.NewJsonTreeImporter(j)
	// err = root.ImportConfig(ctx, jti, "one", 20)
	// if err != nil {
	// 	panic(err)
	// }

	// root.FinishInsertionPhase()

	//
	// Prepare the output
	//

	output := args[2]
	onlyNewOrUpdated, err := strconv.ParseBool(args[4])
	if err != nil {
		panic(err)
	}

	title := "CONFIG"
	if onlyNewOrUpdated {
		title = "CHANGES"
	}

	fmt.Printf("\n%s IN %q FORMAT:\n\n", title, strings.ToUpper(output))
	switch output {
	case "xml":
		x, err := root.ToXML(onlyNewOrUpdated, true, true, false, true)
		if err != nil {
			panic(err)
		}
		x.Indent(2)
		s, err := x.WriteToString()
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
	case "json", "json_ietf":
		if output == "json" {
			j, err = root.ToJson(onlyNewOrUpdated, true)
		} else {
			j, err = root.ToJsonIETF(onlyNewOrUpdated, true)
		}
		if err != nil {
			panic(err)
		}

		byteDoc, err := json.MarshalIndent(j, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Print(string(byteDoc), "\n", "\n")
	}

	concurrentValidation, err := strconv.ParseBool(args[3])
	if err != nil {
		panic(err)
	}

	//
	// perform validation
	//
	// we use a channel and cumulate all the errors
	validationErrors := []error{}
	validationErrChan := make(chan error)
	go func() {
		root.Validate(ctx, validationErrChan, concurrentValidation)
		close(validationErrChan)
	}()

	// read from the Error channel
	for e := range validationErrChan {
		validationErrors = append(validationErrors, e)
	}

	// check if errors are received
	// If so, join them and return the cumulated errors
	if len(validationErrors) > 0 {
		log.Errorf("cumulated validation errors:\n%v", errors.Join(validationErrors...))
	} else {
		log.Info("VALIDATION COMPLETED WITHOUT ANY ISSUES")
	}
}
