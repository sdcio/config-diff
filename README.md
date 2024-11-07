# Overview

# Execution

```
# ./config-diff <schema-definition> <output-format> <concurrent-validation> <onlyNewOrUpdated>


```
- __\<schema-definition>__: Path of the schema-definition *.yaml file
- __\<output-format>__: Format of the output. (json, json_ietf, xml)
- __\<concurrent-validation>__: Internal tree parameter, selects if the validation of the schema should be run in goroutines. Not really relevant to the user, exposed it for performance measurements. Using a debugger, inspecting validation use flase any other case use -> true 
- __\<onlyNewOrUpdated>__: If set to true, only chnaged or new values will included in output. Otherwise the whole config runnning + new and changed.

Example:
```
./config-diff "/home/mava/projects/config-diff/data/schema/arista-4.33.0.F.yaml" xml true true

```