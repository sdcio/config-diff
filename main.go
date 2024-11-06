package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/henderiw/config-diff/schemaclient"
	"github.com/henderiw/config-diff/schemaloader"
	log "github.com/sirupsen/logrus"

	"github.com/sdcio/data-server/pkg/tree"
	treejson "github.com/sdcio/data-server/pkg/tree/importer/json"
	"github.com/sdcio/schema-server/pkg/store/memstore"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

func main() {
	ctx := context.Background()

	args := os.Args
	if len(args) < 1 {
		panic("cannot execute need config and base dir")
	}

	schemastore := memstore.New()
	schemaLoader, err := schemaloader.New(schemastore)
	if err != nil {
		panic(err)
	}
	rsp, err := schemaLoader.LoadSchema(ctx, args[1])
	if err != nil {
		panic(err)
	}
	scb := schemaclient.NewMemSchemaClientBound(schemastore, &sdcpb.Schema{
		Vendor:  rsp.Schema.Vendor,
		Version: rsp.Schema.Version,
	})

	tc := tree.NewTreeContext(tree.NewTreeSchemaCacheClient("dev1", nil, scb), "test")
	root, err := tree.NewTreeRoot(ctx, tc)
	if err != nil {
		panic(err)
	}
	fmt.Println(root.String())

	jsonBytes, err := os.ReadFile("data/config/running/running_srl_01.json")
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

	jsonBytes, err = os.ReadFile("data/config/additions/srl_01.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonBytes, &j)
	if err != nil {
		panic(err)
	}
	jti = treejson.NewJsonTreeImporter(j)
	err = root.ImportConfig(ctx, jti, "one", 20)
	if err != nil {
		panic(err)
	}

	root.FinishInsertionPhase()

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
		x, err := root.ToXML(onlyNewOrUpdated, true, true, false)
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
			j, err = root.ToJson(onlyNewOrUpdated)
		} else {
			j, err = root.ToJsonIETF(onlyNewOrUpdated)
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

	// perform validation
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
