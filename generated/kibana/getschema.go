//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
)

var (
	// Eventually, this should be used instead of merging the others together.
	// Currently it is missing several different schemas and is not complete enough.
	// "https://raw.githubusercontent.com/elastic/kibana/%s/oas_docs/output/kibana.yaml"
	schemaURLs = map[string]string{
		"base":          "https://raw.githubusercontent.com/elastic/kibana/%s/oas_docs/bundle.json",
		"saved_objects": "https://raw.githubusercontent.com/elastic/kibana/%s/packages/core/saved-objects/docs/openapi/bundled.json",
		"data_views":    "https://raw.githubusercontent.com/elastic/kibana/%s/src/plugins/data_views/docs/openapi/bundled.json",
		"alerting":      "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/alerting/docs/openapi/bundled.json",
		"fleet":         "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/fleet/common/openapi/bundled.json",
		"slos":          "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/observability_solution/slo/docs/openapi/slo/bundled.json",
	}
)

func main() {
	_inDir := flag.String("i", "", "input dir")
	_outFile := flag.String("o", "", "output file")
	_apiVersion := flag.String("v", "main", "api version")
	flag.Parse()

	inDir := *_inDir
	outFile := *_outFile
	apiVersion := *_apiVersion

	if inDir == "" || outFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error

	if !pathExists(inDir) {
		if err = os.MkdirAll(inDir, 0755); err != nil {
			log.Fatalf("failed to create directory %s: %v", inDir, err)
		}
	}

	outDir, _ := path.Split(outFile)
	if !pathExists(outDir) {
		if err = os.MkdirAll(outDir, 0755); err != nil {
			log.Fatalf("failed to create directory %s: %v", outDir, err)
		}
	}

	var baseSchema Map
	schemas := make(map[string]*Schema, len(schemaURLs))

	for name, urlTemplate := range schemaURLs {
		var bytes []byte
		filename := fmt.Sprintf("%s.json", name)
		filepath := path.Join(inDir, filename)
		url := fmt.Sprintf(urlTemplate, apiVersion)

		// Download the file if it does not exist
		if pathExists(filepath) {
			if bytes, err = os.ReadFile(filepath); err != nil {
				log.Fatalf("failed to read file %s: %v", filepath, err)
			}
		} else {
			log.Printf("Downloading %s to %s", url, filepath)
			if bytes, err = downloadFile(url); err != nil {
				log.Fatalf("failed to download url %s: %v", url, err)
			}
			if err = os.WriteFile(filepath, bytes, 0664); err != nil {
				log.Fatalf("failed to write file %s: %v", filepath, err)
			}
		}

		if name == "base" {
			var schema Map
			err = json.Unmarshal(bytes, &schema)
			if err != nil {
				log.Fatalf("failed to unmarshal schema from %s: %v", filepath, err)
			}
			baseSchema = schema
		} else {
			var schema Schema
			err = json.Unmarshal(bytes, &schema)
			if err != nil {
				log.Fatalf("failed to unmarshal schema from %s: %v", filepath, err)
			}
			schemas[name] = &schema
		}
	}

	baseSchema["paths"] = make(Map)

	// Order the schemas, or the gen is slightly different each time
	schemaOrder := make([]string, 0, len(schemas))
	for key := range schemas {
		schemaOrder = append(schemaOrder, key)
	}
	slices.Sort(schemaOrder)

	for _, name := range schemaOrder {
		schema := schemas[name]
		for _, fn := range transformers {
			fn(name, schema)
		}

		// Merge each path
		for pathURI, pathInfo := range schema.Paths {
			baseSchema.MustGetMap("paths").Set(pathURI, pathInfo)
			delete(schema.Paths, pathURI)
		}

		// Merge everything in components
		for compKey := range schema.Components {
			// schemas, examples, parameters
			compVals := schema.Components.MustGetMap(compKey)
			for key, val := range compVals {
				baseKey := fmt.Sprintf("components.%s.%s", compKey, key)
				baseSchema.Set(baseKey, val)
			}
		}
	}

	saveFile(baseSchema, outFile)
}

// Check if the given path exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Download the file at url.
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func saveFile(obj any, path string) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		log.Fatalf("failed to marshal to file %s: %v", path, err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0664); err != nil {
		log.Fatalf("failed to write file %s: %v", path, err)
	}
}

func deepCopy[T any](src T) T {
	var dest T

	jsonStr, err := json.Marshal(src)
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(jsonStr, &dest)
	if err != nil {
		log.Panic(err)
	}

	return dest
}

// ============================================================================

type Schema struct {
	Paths      map[string]*Path `json:"paths"`
	Version    string           `json:"openapi"`
	Tags       []Map            `json:"tags,omitempty"`
	Servers    []Map            `json:"servers,omitempty"`
	Components Map              `json:"components,omitempty"`
	Security   []Map            `json:"security,omitempty"`
	Info       Map              `json:"info"`
}

func (s Schema) MustGetPath(compName string) *Path {
	p, ok := s.Paths[compName]
	if !ok {
		log.Panicf("%s not found", compName)
	}
	return p
}

// ============================================================================

type Path struct {
	Parameters []Map `json:"parameters,omitempty"`
	Get        Map   `json:"get,omitempty"`
	Post       Map   `json:"post,omitempty"`
	Put        Map   `json:"put,omitempty"`
	Delete     Map   `json:"delete,omitempty"`
}

func (p *Path) Endpoints() map[string]Map {
	endpoints := map[string]Map{}
	if p.Get != nil {
		endpoints["get"] = p.Get
	}
	if p.Post != nil {
		endpoints["post"] = p.Post
	}
	if p.Put != nil {
		endpoints["put"] = p.Put
	}
	if p.Delete != nil {
		endpoints["delete"] = p.Delete
	}
	return endpoints
}

func (p *Path) Set(method string, value Map) {
	switch method {
	case "get":
		p.Get = value
	case "post":
		p.Post = value
	case "put":
		p.Put = value
	case "delete":
		p.Delete = value
	default:
		log.Panicf("Invalid method %s", method)
	}
}

// ============================================================================

type Map map[string]any

func (m Map) Get(key string) (any, bool) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		switch t := m[rootKey].(type) {
		case Map:
			return t.Get(subKeys)
		case map[string]any:
			return Map(t).Get(subKeys)
		case Slice:
			return t.Get(subKeys)
		case []any:
			return Slice(t).Get(subKeys)
		default:
			rootKey = key
		}
	}

	value, ok := m[rootKey]
	return value, ok
}

func (m Map) MustGet(key string) any {
	v, ok := m.Get(key)
	if !ok {
		log.Panicf("%s not found", key)
	}
	return v
}

func (m Map) GetSlice(key string) (Slice, bool) {
	value, ok := m.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Slice:
		return t, true
	case []any:
		return t, true
	}

	log.Panicf("%s is not a slice", key)
	return nil, false
}

func (m Map) MustGetSlice(key string) Slice {
	v, ok := m.GetSlice(key)
	if !ok {
		log.Panicf("%s not found", key)
	}
	return v
}

func (m Map) GetMap(key string) (Map, bool) {
	value, ok := m.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Map:
		return t, true
	case map[string]any:
		return t, true
	}

	log.Panicf("%s is not a map", key)
	return nil, false
}

func (m Map) MustGetMap(key string) Map {
	v, ok := m.GetMap(key)
	if !ok {
		log.Panicf("%s not found", key)
	}
	return v
}

// Set will set key to the value of "value".
func (m Map) Set(key string, value any) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		if v, ok := m[rootKey]; ok {
			switch t := v.(type) {
			case Map:
				t.Set(subKeys, value)
			case map[string]any:
				Map(t).Set(subKeys, value)
			}
		} else {
			subMap := Map{}
			subMap.Set(subKeys, value)
			m[rootKey] = subMap
		}
	} else {
		m[rootKey] = value
	}
}

// Move will move the value from "key" to "target".
// If "key" does not exist, the operation is a no-op.
func (m Map) Move(key, target string) {
	value, ok := m.Get(key)
	if !ok {
		return
	}
	m.Set(target, value)
	m.Delete(key)
}

// Delete will remove the key from the map.
// If key is nested, empty sub-keys will be removed as well.
func (m Map) Delete(key string) bool {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		if v, ok := m[rootKey]; ok {
			switch t := v.(type) {
			case Map:
				t.Delete(subKeys)
			case map[string]any:
				Map(t).Delete(subKeys)
			}
		}
	} else {
		delete(m, rootKey)
		return true
	}
	return false
}

func (m Map) MustDelete(key string) {
	if !m.Delete(key) {
		log.Panicf("%s not found", key)
	}
}

func (m Map) CreateRef(name string, key string, schema *Schema) {
	compPath := fmt.Sprintf("schemas.%s", name)
	compRef := fmt.Sprintf("#/components/schemas/%s", name)
	i := strings.LastIndex(key, ".")
	mapPath := key[:i]
	mapField := key[i+1:]

	m.MustGet(key) // Check the full path
	target := m.MustGetMap(mapPath)

	// Don't overwrite existing schemas
	if _, ok := schema.Components.Get(compPath); !ok {
		value := target.MustGet(mapField)
		schema.Components.Set(compPath, value)
	}

	target.MustDelete(mapField)
	target.Set(mapField+".$ref", compRef)
}

// ============================================================================

type Slice []any

func (s Slice) Get(key string) (any, bool) {
	indexStr, subKeys, found := strings.Cut(key, ".")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		log.Panicf("Failed to parse slice index key %q: %v", indexStr, err)
		return nil, false
	}
	if index < 0 || index >= len(s) {
		log.Panicf("Slice index is out of bounds (%d, target slice len: %d)", index, len(s))
		return nil, false
	}

	if found {
		switch t := s[index].(type) {
		case Map:
			return t.Get(subKeys)
		case map[string]any:
			return Map(t).Get(subKeys)
		case Slice:
			return t.Get(subKeys)
		case []any:
			return Slice(t).Get(subKeys)
		}
	}

	value := s[index]
	return value, true
}

// Contains returns true if s contains value.
func (s Slice) Contains(value string) bool {
	for _, v := range s {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if value == s {
			return true
		}
	}

	return false
}

// ============================================================================

type TransformFunc func(name string, schema *Schema)

var transformers = []TransformFunc{
	transformFilterPaths,
	transformRemoveKbnXsrf,
	transformSimplifyContentType,
	transformFleetEpmInlinePackageDefinitions,
	transformFleetNewPackagePolicy,
	transformFleetOutputTypeRequired,
	//transformDataViewPaths,
	//transformAlertingRulePaths,
	//transformConnectorPaths,
	transformRemoveUnnecessaryGoPointers,
}

// transformFilterPaths filters the paths in a schema down to
// a specified list of endpoints and methods.
func transformFilterPaths(_ string, schema *Schema) {
	var includePaths = map[string][]string{
		"/agent_policies":                             {"get", "post"},
		"/agent_policies/delete":                      {"post"},
		"/agent_policies/{agentPolicyId}":             {"get", "put"},
		"/api/actions/connectors":                     {"get"},
		"/api/actions/connector":                      {"post"},
		"/api/actions/connector/{connectorId}":        {"get", "post", "put", "delete"},
		"/api/alerting/rules/_find":                   {"get"},
		"/api/alerting/rule":                          {"post"},
		"/api/alerting/rule/{ruleId}":                 {"get", "post", "put", "delete"},
		"/api/data_views":                             {"get"},
		"/api/data_views/data_view":                   {"post"},
		"/api/data_views/data_view/{viewId}":          {"get", "post", "delete"},
		"/api/saved_objects/_import":                  {"post"},
		"/enrollment_api_keys":                        {"get"},
		"/epm/packages":                               {"get", "post"},
		"/epm/packages/{pkgName}/{pkgVersion}":        {"get", "post", "put", "delete"},
		"/fleet_server_hosts":                         {"get", "post"},
		"/fleet_server_hosts/{itemId}":                {"get", "put", "delete"},
		"/outputs":                                    {"get", "post"},
		"/outputs/{outputId}":                         {"get", "put", "delete"},
		"/package_policies":                           {"get", "post"},
		"/package_policies/{packagePolicyId}":         {"get", "put", "delete"},
		"/s/{spaceId}/api/observability/slos":         {"get", "post"},
		"/s/{spaceId}/api/observability/slos/{sloId}": {"get", "put", "delete"},
		"/service_tokens":                             {"post"},
	}

	for path, pathInfo := range schema.Paths {
		// Remove paths not in filter list.
		if _, exists := includePaths[path]; !exists {
			delete(schema.Paths, path)
			continue
		}

		// Filter out endpoints not if filter list
		allowedMethods := includePaths[path]
		for method := range pathInfo.Endpoints() {
			if !slices.Contains(allowedMethods, method) {
				switch method {
				case "get":
					pathInfo.Get = nil
				case "post":
					pathInfo.Post = nil
				case "put":
					pathInfo.Put = nil
				case "delete":
					pathInfo.Delete = nil
				default:
					log.Panicf("Unhandled method %s", method)
				}
			}
		}
	}
}

// transformRemoveKbnXsrf removes the kbn-xsrf header as it
// is already applied in the client.
func transformRemoveKbnXsrf(_ string, schema *Schema) {
	removeKbnXsrf := func(fields Map) {
		if params, ok := fields.GetSlice("parameters"); ok {
			for i, _param := range params {
				param := _param.(map[string]any)
				if _ref, hasRef := param["$ref"]; hasRef {
					ref := _ref.(string)
					// Data_views_kbn_xsrf, Saved_objects_kbn_xsrf, etc
					if strings.HasPrefix(ref, "#/components/parameters/") && strings.HasSuffix(ref, "kbn_xsrf") {
						newParams := append(params[:i], params[i+1:]...)
						fields.Set("parameters", newParams)
						break
					}
				}
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methodInfo := range pathInfo.Endpoints() {
			removeKbnXsrf(methodInfo)
		}
	}
}

// transformSimplifyContentType simplifies Content-Type headers such
// as 'application/json; Elastic-Api-Version=2023-10-31'
func transformSimplifyContentType(_ string, schema *Schema) {
	simplifyContentType := func(fields Map) {
		if content, ok := fields.GetMap("content"); ok {
			for key := range content {
				idx := strings.Index(key, ";")
				if idx != -1 {
					newKey := key[:idx]
					content.Move(key, newKey)
				}
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, endpoint := range pathInfo.Endpoints() {
			if req, ok := endpoint.GetMap("requestBody"); ok {
				simplifyContentType(req)
			}
			if resp, ok := endpoint.GetMap("responses"); ok {
				for _, _respInfo := range resp {
					respInfo := _respInfo.(map[string]any)
					simplifyContentType(respInfo)
				}
			}
		}
	}

	if responses, ok := schema.Components.GetMap("responses"); ok {
		for key := range responses {
			respInfo := responses.MustGetMap(key)
			simplifyContentType(respInfo)
		}
	}
}

// transformFleetEpmInlinePackageDefinitions relocates inline type definitions for the
// EPM endpoints to the dedicated schemas section of the OpenAPI schema. This needs
// to be done as there is a bug in the OpenAPI generator which causes types to
// be generated with invalid names:
func transformFleetEpmInlinePackageDefinitions(name string, schema *Schema) {
	if name != "fleet" {
		return
	}

	epmPkgsPath := schema.MustGetPath("/epm/packages")
	epmPkgPath := schema.MustGetPath("/epm/packages/{pkgName}/{pkgVersion}")

	epmPkgsPath.Post.CreateRef("package_install_source", "responses.200.content.application/json.schema.properties._meta.properties.install_source", schema)
	epmPkgsPath.Post.CreateRef("package_item", "responses.200.content.application/json.schema.properties.items", schema)
	epmPkgPath.Get.CreateRef("package_status", "responses.200.content.application/json.schema.allOf.1.properties.status", schema)
	epmPkgPath.Post.CreateRef("package_install_source", "responses.200.content.application/json.schema.properties._meta.properties.install_source", schema)
	epmPkgPath.Post.CreateRef("package_item", "responses.200.content.application/json.schema.properties.items", schema)
	epmPkgPath.Put.CreateRef("package_item", "responses.200.content.application/json.schema.properties.items", schema)
	epmPkgPath.Delete.CreateRef("package_item", "responses.200.content.application/json.schema.properties.items", schema)
	schema.Components.CreateRef("package_policy_request_input_stream", "schemas.package_policy_request.properties.inputs.additionalProperties.properties.streams.additionalProperties", schema)
	schema.Components.CreateRef("package_policy_request_input", "schemas.package_policy_request.properties.inputs.additionalProperties", schema)
	schema.Components.CreateRef("package_policy_package_info", "schemas.new_package_policy.properties.package", schema)
	schema.Components.CreateRef("package_item_type", "schemas.package_item.items.properties.type", schema)
}

// transformFleetNewPackagePolicy fixes several fields on the "new_package_policy"
// component schemas so as to conform with the Fleet API.
// - Add the missing "vars" object
// - Add the missing "secret_references" array and child object
// - Convert "inputs" to an object from an array
// - Convert "inputs.streams" to an object from an array
// - Add the missing "inputs.streams" object
// - Add the missing "secret_references" array
// - Extract the above nested objects
func transformFleetNewPackagePolicy(name string, schema *Schema) {
	if name != "fleet" {
		return
	}

	// Add vars
	props := schema.Components.MustGetMap("schemas.new_package_policy.properties")
	if _, ok := props.Get("vars"); !ok {
		props.Set("vars.type", "object")
	}

	// Add secret_references
	if _, ok := props.Get("secret_references"); !ok {
		props.Set("secret_references", map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type": "string",
					},
				},
			},
		})
	}

	// Convert inputs array to map
	inputs := schema.Components.MustGetMap("schemas.new_package_policy.properties.inputs")
	inputs.Set("type", "object")
	inputs.Move("items", "additionalProperties")

	// Extract nested
	schema.Components.CreateRef("new_package_policy_input", "schemas.new_package_policy.properties.inputs.additionalProperties", schema)

	// Convert input.streams array to map
	streams := schema.Components.MustGetMap("schemas.new_package_policy_input.properties.streams")
	streams.Set("type", "object")
	streams.Move("items", "additionalProperties")

	// Define the missing object
	streams.Set("additionalProperties", map[string]any{
		"properties": map[string]any{
			"enabled": map[string]any{"type": "boolean"},
			"vars":    map[string]any{"type": "object"},
		},
	})

	// Extract nested
	schema.Components.CreateRef("new_package_policy_input_stream", "schemas.new_package_policy_input.properties.streams.additionalProperties", schema)
}

// transformFleetOutputTypeRequired ensures that the "type" key is
// in the list of required keys for the output type.
func transformFleetOutputTypeRequired(name string, schema *Schema) {
	if name != "fleet" {
		return
	}

	path := []string{
		"schemas.output_create_request_elasticsearch.required",
		"schemas.output_create_request_kafka.required",
		"schemas.output_create_request_logstash.required",
		"schemas.output_create_request_remote_elasticsearch.required",
		"schemas.output_update_request_elasticsearch.required",
		"schemas.output_update_request_kafka.required",
		"schemas.output_update_request_logstash.required",
		// "schemas.output_update_request_remote_elasticsearch.required", // MISSING
	}

	for _, v := range path {
		slice := schema.Components.MustGetSlice(v)
		if slice.Contains("type") {
			continue
		}
		slice = append(slice, "type")
		schema.Components.Set(v, slice)
	}
}

// transformDataViewPaths clones the DataView paths to support a spaceID
// parameter. Additionally, cleans up some naming conventions.
func transformDataViewPaths(schema *Schema) {
	{
		if _, ok := schema.Components.GetMap("parameters.Data_views_space_id"); ok {
			log.Panic("parameters.Data_views_space_id already exists")
		}

		// Same as SLOs_space_id, so copy it
		param := schema.Components.MustGetMap("parameters.SLOs_space_id")
		param = deepCopy(param)
		schema.Components.Set("parameters.Data_views_space_id", param)
	}

	fixPath := func(path string) {
		newPath := fmt.Sprintf("/s/{spaceId}%s", path)
		if _, ok := schema.Paths[newPath]; ok {
			log.Panicf(`path "%s" already exists`, newPath)
		}
		pathInfo := schema.MustGetPath(path)
		pathInfo = deepCopy(pathInfo)
		schema.Paths[newPath] = pathInfo

		// Remove "Default" from the operation ID
		for _, endpoint := range pathInfo.Endpoints() {
			opId := endpoint.MustGet("operationId").(string)
			opId = strings.ReplaceAll(opId, "Default", "")
			endpoint.Set("operationId", opId)
		}

		// Add the spaceId path parameter
		for _, endpoint := range pathInfo.Endpoints() {
			if params, ok := endpoint.GetSlice("parameters"); ok {
				params = append(params, map[string]any{"$ref": "#/components/parameters/Data_views_space_id"})
				endpoint.Set("parameters", params)
			} else {
				params := []any{map[string]any{"$ref": "#/components/parameters/Data_views_space_id"}}
				endpoint.Set("parameters", params)
			}
		}
	}

	schema.MustGetPath("/api/data_views/data_view").
		Post.Set("operationId", "createDataViewDefault") // has a trailing "w"

	fixPath("/api/data_views")
	fixPath("/api/data_views/data_view")
	fixPath("/api/data_views/data_view/{viewId}")
}

// transformAlertingRulePaths clones the AlertingRule paths to support a spaceID
// parameter.
func transformAlertingRulePaths(schema *Schema) {
	{
		if _, ok := schema.Components.GetMap("parameters.Alerting_rules_space_id"); ok {
			log.Panic("parameters.Alerting_rules_space_id already exists")
		}

		// Same as SLOs_space_id, so copy it
		param := schema.Components.MustGetMap("parameters.SLOs_space_id")
		param = deepCopy(param)
		schema.Components.Set("parameters.Alerting_rules_space_id", param)
	}

	fixPath := func(path string) {
		newPath := fmt.Sprintf("/s/{spaceId}%s", path)
		if _, ok := schema.Paths[newPath]; ok {
			log.Panicf(`path "%s" already exists`, newPath)
		}
		pathInfo := schema.MustGetPath(path)
		// Add "Default" to the original operation ID
		for _, endpoint := range pathInfo.Endpoints() {
			if endpoint == nil {
				continue
			}
			opId := endpoint.MustGet("operationId").(string)
			opId += "Default"
			endpoint.Set("operationId", opId)
		}

		pathInfo = deepCopy(pathInfo)
		schema.Paths[newPath] = pathInfo

		// Remove "Default" from the operation ID
		for _, endpoint := range pathInfo.Endpoints() {
			if endpoint == nil {
				continue
			}
			opId := endpoint.MustGet("operationId").(string)
			opId = strings.ReplaceAll(opId, "Default", "")
			endpoint.Set("operationId", opId)
		}

		// Add the spaceId path parameter
		for _, endpoint := range pathInfo.Endpoints() {
			if endpoint == nil {
				continue
			}
			if params, ok := endpoint.GetSlice("parameters"); ok {
				params = append(params, map[string]any{"$ref": "#/components/parameters/Alerting_rules_space_id"})
				endpoint.Set("parameters", params)
			} else {
				params := []any{map[string]any{"$ref": "#/components/parameters/Alerting_rules_space_id"}}
				endpoint.Set("parameters", params)
			}
		}
	}

	fixPath("/api/alerting/rule")
	fixPath("/api/alerting/rule/{ruleId}")
}

// transformConnectorPaths clones the Connector paths to support a spaceID
// parameter.
func transformConnectorPaths(schema *Schema) {
	{
		if _, ok := schema.Components.GetMap("parameters.Connectors_space_id"); ok {
			log.Panic("parameters.Connectors_space_id already exists")
		}

		// Same as SLOs_space_id, so copy it
		param := schema.Components.MustGetMap("parameters.SLOs_space_id")
		param = deepCopy(param)
		schema.Components.Set("parameters.Connectors_space_id", param)
	}

	fixPath := func(path string) {
		newPath := fmt.Sprintf("/s/{spaceId}%s", path)
		if _, ok := schema.Paths[newPath]; ok {
			log.Panicf(`path "%s" already exists`, newPath)
		}
		pathInfo := schema.MustGetPath(path)
		// Add "Default" to the original operation ID
		for _, endpoint := range pathInfo.Endpoints() {
			opId := endpoint.MustGet("operationId").(string)
			opId += "Default"
			endpoint.Set("operationId", opId)
		}

		pathInfo = deepCopy(pathInfo)
		schema.Paths[newPath] = pathInfo

		// Remove "Default" from the operation ID
		for _, endpoint := range pathInfo.Endpoints() {
			if endpoint == nil {
				continue
			}
			opId := endpoint.MustGet("operationId").(string)
			opId = strings.ReplaceAll(opId, "Default", "")
			endpoint.Set("operationId", opId)
		}

		// Add the spaceId path parameter
		for _, endpoint := range pathInfo.Endpoints() {
			if endpoint == nil {
				continue
			}
			if params, ok := endpoint.GetSlice("parameters"); ok {
				params = append(params, map[string]any{"$ref": "#/components/parameters/Connectors_space_id"})
				endpoint.Set("parameters", params)
			} else {
				params := []any{map[string]any{"$ref": "#/components/parameters/Connectors_space_id"}}
				endpoint.Set("parameters", params)
			}
		}
	}

	fixPath("/api/actions/connector")
	fixPath("/api/actions/connector/{connectorId}")
	fixPath("/api/actions/connectors")
}

// transformRemoveUnnecessaryGoPointers removes pointers from
// map and slice objects.
func transformRemoveUnnecessaryGoPointers(_ string, schema *Schema) {
	var iterate func(val any)
	iterate = func(val any) {
		switch tval := val.(type) {
		case []any:
			for _, v := range tval {
				iterate(v)
			}
		case map[string]any:
			iterate(Map(tval))
		case Map:
			for _, v := range tval {
				iterate(v)
			}
			if vtype, ok := tval["type"]; ok {
				switch vtype {
				case "array":
					tval["x-go-type-skip-optional-pointer"] = true
				case "object":
					if _, ok := tval["properties"]; !ok {
						tval["x-go-type-skip-optional-pointer"] = true
					}
				}
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methInfo := range pathInfo.Endpoints() {
			iterate(methInfo)
		}
	}
	iterate(schema.Components)
}

/*
// transformFixPackageSearchResult removes unneeded fields from the
// SearchResult struct. These fields are also causing parsing errors.
func transformFixPackageSearchResult(schema *Schema) {
	properties, ok := schema.Components.GetMap("schemas.search_result.properties")
	if !ok {
		panic("properties not found")
	}
	properties.Delete("icons")
	properties.Delete("installationInfo")
}
*/
