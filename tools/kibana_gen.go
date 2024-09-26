package tools

//go:generate go run ../generated/kibana/getschema.go -v main -i ../generated/kibana/schemas -o ../generated/kibana/kibana-filtered.json
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package=kibana -generate=types,client -o ../generated/kibana/kibana.gen.go ../generated/kibana/kibana-filtered.json
