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
	"unicode"
)

var (
	// Eventually, this should be used instead of merging the others together.
	// Currently it is missing several different schemas and is not complete enough.
	// "https://raw.githubusercontent.com/elastic/kibana/%s/oas_docs/output/kibana.yaml"
	schemaURLs = map[string]string{
		"oas":           "https://raw.githubusercontent.com/elastic/kibana/%s/oas_docs/bundle.json", // alerting and connectors
		"saved_objects": "https://raw.githubusercontent.com/elastic/kibana/%s/packages/core/saved-objects/docs/openapi/bundled.json",
		"data_views":    "https://raw.githubusercontent.com/elastic/kibana/%s/src/plugins/data_views/docs/openapi/bundled.json",
		"fleet":         "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/fleet/common/openapi/bundled.json",
		"slos":          "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/observability_solution/slo/docs/openapi/slo/bundled.json",

		// Does not include connectors
		//"alerting":      "https://raw.githubusercontent.com/elastic/kibana/%s/x-pack/plugins/alerting/docs/openapi/bundled.json",
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
			log.Fatalf("failed to create directory %q: %v", inDir, err)
		}
	}

	outDir, _ := path.Split(outFile)
	if !pathExists(outDir) {
		if err = os.MkdirAll(outDir, 0755); err != nil {
			log.Fatalf("failed to create directory %q: %v", outDir, err)
		}
	}

	baseSchema := Map{
		"components": Map{
			"schemas": Map{},
			"securitySchemes": Map{
				"apiKeyAuth": Map{"in": "header", "name": "Authorization", "type": "apiKey"},
				"basicAuth":  Map{"scheme": "basic", "type": "http"},
			},
		},
		"info": Map{
			"title":   "Kibana HTTP APIs",
			"version": "0.0.0",
		},
		"openapi":  "3.0.0",
		"paths":    Map{},
		"security": Slice{},
		"servers":  Slice{},
		"tags":     Slice{},
	}

	schemas := make(map[string]*Schema, len(schemaURLs))

	for name, urlTemplate := range schemaURLs {
		var bytes []byte
		filename := fmt.Sprintf("%s.json", name)
		filepath := path.Join(inDir, filename)
		url := fmt.Sprintf(urlTemplate, apiVersion)

		// Download the file if it does not exist
		if pathExists(filepath) {
			if bytes, err = os.ReadFile(filepath); err != nil {
				log.Fatalf("failed to read file %q: %v", filepath, err)
			}
		} else {
			log.Printf("Downloading %q to %q", url, filepath)
			if bytes, err = downloadFile(url); err != nil {
				log.Fatalf("failed to download url %q: %v", url, err)
			}
			if err = os.WriteFile(filepath, bytes, 0664); err != nil {
				log.Fatalf("failed to write file %q: %v", filepath, err)
			}
		}

		var schema Schema
		err = json.Unmarshal(bytes, &schema)
		if err != nil {
			log.Fatalf("failed to unmarshal schema from %q: %v", filepath, err)
		}
		schemas[name] = &schema
	}

	// Order the schemas, or the gen is slightly different each time
	names := make([]string, 0, len(schemas))
	for name := range schemas {
		names = append(names, name)
	}
	slices.Sort(names)

	// Merge each child schema
	for _, name := range names {
		schema := schemas[name]

		// Merge each path
		for pathURI, pathInfo := range schema.Paths {
			baseSchema.MustGetMap("paths").Set(pathURI, pathInfo)
			delete(schema.Paths, pathURI)
		}

		// Merge everything in components
		for compKey := range schema.Components {
			compVals := schema.Components.MustGetMap(compKey)
			for key, val := range compVals {
				baseKey := fmt.Sprintf("components.%s.%s", compKey, key)
				baseSchema.Set(baseKey, val)
			}
		}
	}

	// Convert base back to a Schema type
	b, err := json.Marshal(baseSchema)
	if err != nil {
		log.Panic("failed to marshal base schema")
	}
	var schema Schema
	if err = json.Unmarshal(b, &schema); err != nil {
		log.Panic("failed to unmarshal base schema")
	}

	for _, fn := range transformers {
		fn(&schema)
	}

	saveFile(schema, outFile)
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
		log.Fatalf("failed to marshal to file %q: %v", path, err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0664); err != nil {
		log.Fatalf("failed to write file %q: %v", path, err)
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

func toSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	sep := '_'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		if unicode.IsLetter(v) {
			// v is capital letter here
			// irregard first letter
			// add underscore if last letter is capital letter
			// add underscore when previous letter is lowercase
			// add underscore when next letter is lowercase
			if (i != 0 || i == l-1) && (          // head and tail
			(i > 0 && rune(camel[i-1]) >= 'a') || // pre
				(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
				b.WriteRune(sep)
			}
			b.WriteRune(v + diff)
		} else {
			b.WriteRune(v)
		}
	}
	return b.String()
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

func (s Schema) MustGetPath(path string) *Path {
	p, ok := s.Paths[path]
	if !ok {
		log.Panicf("%q not found", path)
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

func (p *Path) GetMethod(method string) Map {
	switch method {
	case "get":
		return p.Get
	case "post":
		return p.Post
	case "put":
		return p.Put
	case "delete":
		return p.Delete
	default:
		log.Panicf("Unhandled method %q", method)
	}
	return nil
}

func (p *Path) SetMethod(method string, value Map) {
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
		log.Panicf("Invalid method %q", method)
	}
}

// ============================================================================

type Map map[string]any

func (m Map) Keys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func (m Map) Has(key string) bool {
	_, ok := m.Get(key)
	return ok
}

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
		log.Panicf("%q not found", key)
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

	log.Panicf("%q is not a slice", key)
	return nil, false
}

func (m Map) MustGetSlice(key string) Slice {
	v, ok := m.GetSlice(key)
	if !ok {
		log.Panicf("%q not found", key)
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

	log.Panicf("%q is not a map", key)
	return nil, false
}

func (m Map) MustGetMap(key string) Map {
	v, ok := m.GetMap(key)
	if !ok {
		log.Panicf("%q not found", key)
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

// Move will move the value from "src" to "dst".
func (m Map) Move(src string, dst string) {
	value := m.MustGet(src)
	m.Set(dst, value)
	m.Delete(src)
}

// Delete will remove the key from the map.
// If key is nested, empty sub-keys will be removed as well.
func (m Map) Delete(key string) bool {
	rootKey, subKeys, found := strings.Cut(key, ".")
	if found {
		if v, ok := m[rootKey]; ok {
			switch t := v.(type) {
			case Slice:
				return t.Delete(subKeys)
			case []any:
				return Slice(t).Delete(subKeys)
			case Map:
				return t.Delete(subKeys)
			case map[string]any:
				return Map(t).Delete(subKeys)
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
		log.Panicf("%q not found", key)
	}
}

func (m Map) CreateRef(name string, key string, schema *Schema) {
	m.MustGet(key) // Check the full path
	refPath := fmt.Sprintf("schemas.%s", name)
	refValue := fmt.Sprintf("#/components/schemas/%s", name)

	// Don't overwrite existing schemas
	if schema.Components.Has(refPath) {
		log.Panicf("Duplicate component %q", refPath)
		return
	}

	var target any
	var targetKey string
	i := strings.LastIndex(key, ".")
	if i == -1 {
		target = m
		targetKey = key
	} else {
		target = m.MustGet(key[:i])
		targetKey = key[i+1:]
	}

	doMap := func(target Map, key string) {
		schema.Components.Set(refPath, target.MustGet(key))
		target.Set(key, Map{"$ref": refValue})
	}

	doSlice := func(target Slice, key string) {
		index, err := strconv.Atoi(key)
		if err != nil {
			log.Panicf("Failed to parse slice index key %q: %v", key, err)
			return
		}
		if index < 0 || index >= len(target) {
			log.Panicf("Slice index is out of bounds (%d, target slice len: %d)", index, len(target))
			return
		}
		schema.Components.Set(refPath, target[index])
		target[index] = Map{"$ref": refValue}
	}

	switch t := target.(type) {
	case map[string]any:
		doMap(Map(t), targetKey)
	case Map:
		doMap(t, targetKey)
	case []any:
		doSlice(Slice(t), targetKey)
	case Slice:
		doSlice(t, targetKey)
	default:
		log.Panicf("Cannot create a ref of target type %T at %q", target, key)
	}
}

func (m Map) Iterate(iteratee func(key string, node Map)) {
	joinPath := func(existing string, next string) string {
		if existing == "" {
			return next
		} else {
			return fmt.Sprintf("%s.%s", existing, next)
		}
	}
	joinIndex := func(existing string, next int) string {
		if existing == "" {
			return fmt.Sprintf("%d", next)
		} else {
			return fmt.Sprintf("%s.%d", existing, next)
		}
	}

	var iterate func(key string, val any)
	iterate = func(key string, val any) {
		switch tval := val.(type) {
		case []any:
			iterate(key, Slice(tval))
		case Slice:
			for i, v := range tval {
				iterate(joinIndex(key, i), v)
			}
		case map[string]any:
			iterate(key, Map(tval))
		case Map:
			for _, k := range tval.Keys() {
				iterate(joinPath(key, k), tval[k])
			}
			iteratee(key, tval)
		}
	}

	iterate("", m)
}

// ============================================================================

type Slice []any

func (s Slice) Get(key string) (any, bool) {
	rootKey, subKeys, found := strings.Cut(key, ".")
	index, err := strconv.Atoi(rootKey)
	if err != nil {
		log.Panicf("Failed to parse slice index key %q: %v", rootKey, err)
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

func (s Slice) GetMap(key string) (Map, bool) {
	value, ok := s.Get(key)
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case Map:
		return t, true
	case map[string]any:
		return t, true
	}

	log.Panicf("%q is not a map", key)
	return nil, false
}

func (s Slice) MustGetMap(key string) Map {
	v, ok := s.GetMap(key)
	if !ok {
		log.Panicf("%q not found", key)
	}
	return v
}

// Delete will remove the index from the slice.
// If key is nested, empty sub-keys will be removed as well.
func (s Slice) Delete(key string) bool {
	rootKey, subKeys, found := strings.Cut(key, ".")
	index, err := strconv.Atoi(rootKey)
	if err != nil {
		log.Panicf("Failed to parse slice index key %q: %v", rootKey, err)
		return false
	}
	if index < 0 || index >= len(s) {
		log.Panicf("Slice index is out of bounds (%d, target slice len: %d)", index, len(s))
		return false
	}
	if found {
		item := (s)[index]
		switch t := item.(type) {
		case Slice:
			return t.Delete(subKeys)
		case []any:
			return Slice(t).Delete(subKeys)
		case Map:
			return t.Delete(subKeys)
		case map[string]any:
			return Map(t).Delete(subKeys)
		}
	} else {
		log.Panicf("Unable to delete from slice directly")
		return true
	}
	return false
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

type TransformFunc func(schema *Schema)

var transformers = []TransformFunc{
	transformFilterPaths,
	transformRemoveKbnXsrf,
	transformRemoveApiVersionParam,
	transformSimplifyContentType,
	transformAlertingPaths,
	transformConnectorPaths,
	transformDataViewPaths,
	transformFleetPaths,
	transformExtractAndSimplifyNested,
}

// transformFilterPaths filters the paths in a schema down to
// a specified list of endpoints and methods.
func transformFilterPaths(schema *Schema) {
	var includePaths = map[string][]string{
		"/api/fleet/agent_policies":                      {"get", "post"},
		"/api/fleet/agent_policies/delete":               {"post"},
		"/api/fleet/agent_policies/{agentPolicyId}":      {"get", "put"},
		"/api/actions/connectors":                        {"get"},
		"/api/actions/connector/{id}":                    {"get", "post", "put", "delete"},
		"/api/alerting/rules/_find":                      {"get"},
		"/api/alerting/rule/{id}":                        {"get", "post", "put", "delete"},
		"/api/data_views":                                {"get"},
		"/api/data_views/data_view":                      {"post"},
		"/api/data_views/data_view/{viewId}":             {"get", "post", "delete"},
		"/api/saved_objects/_import":                     {"post"},
		"/api/fleet/enrollment_api_keys":                 {"get"},
		"/api/fleet/epm/packages":                        {"get", "post"},
		"/api/fleet/epm/packages/{pkgName}/{pkgVersion}": {"get", "post", "put", "delete"},
		"/api/fleet/fleet_server_hosts":                  {"get", "post"},
		"/api/fleet/fleet_server_hosts/{itemId}":         {"get", "put", "delete"},
		"/api/fleet/outputs":                             {"get", "post"},
		"/api/fleet/outputs/{outputId}":                  {"get", "put", "delete"},
		"/api/fleet/package_policies":                    {"get", "post"},
		"/api/fleet/package_policies/{packagePolicyId}":  {"get", "put", "delete"},
		"/api/fleet/service_tokens":                      {"post"},
		"/s/{spaceId}/api/observability/slos":            {"get", "post"},
		"/s/{spaceId}/api/observability/slos/{sloId}":    {"get", "put", "delete"},
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
				pathInfo.SetMethod(method, nil)
			}
		}
	}

	// Go through again, verify each entry exists
	for path, methods := range includePaths {
		pathInfo := schema.MustGetPath(path)
		for _, method := range methods {
			methodInfo := pathInfo.GetMethod(method)
			if methodInfo == nil {
				log.Panicf("Method %q of %q missing", method, path)
			}
		}
	}
}

// transformRemoveKbnXsrf removes the kbn-xsrf header as it
// is already applied in the client.
func transformRemoveKbnXsrf(schema *Schema) {
	for _, pathInfo := range schema.Paths {
		for _, methodInfo := range pathInfo.Endpoints() {
			if params, ok := methodInfo.GetSlice("parameters"); ok {
				for i, _param := range params {
					param := _param.(map[string]any)
					if _ref, hasRef := param["$ref"]; hasRef {
						ref := _ref.(string)
						// Data_views_kbn_xsrf, Saved_objects_kbn_xsrf, etc
						if strings.HasSuffix(ref, "kbn_xsrf") || strings.HasSuffix(ref, "kbn-xsrf") {
							newParams := slices.Delete(params, i, i+1)
							methodInfo.Set("parameters", newParams)
							break
						}
					}
				}
			}
		}
	}
}

// transformRemoveApiVersionParam removes the Elastic API Version
// query parameter header as it is not used for our purposes.
func transformRemoveApiVersionParam(schema *Schema) {
	removeApiVersion := func(fields Map) {
		if params, ok := fields.GetSlice("parameters"); ok {
			for i := len(params) - 1; i >= 0; i-- {
				param := params.MustGetMap(fmt.Sprintf("%d", i))
				if name, ok := param.Get("name"); ok && name.(string) == "elastic-api-version" {
					params = slices.Delete(params, i, i+1)
					fields.Set("parameters", params)
				}
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methodInfo := range pathInfo.Endpoints() {
			removeApiVersion(methodInfo)
		}
	}
}

// transformSimplifyContentType simplifies Content-Type headers such
// as 'application/json; Elastic-Api-Version=2023-10-31'
func transformSimplifyContentType(schema *Schema) {
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

// transformFleetPaths fixes the fleet paths.
func transformFleetPaths(schema *Schema) {
	// Add new_package_policy.vars
	props := schema.Components.MustGetMap("schemas.new_package_policy.properties")
	props.Set("vars.type", "object")

	// Add new_package_policy.secret_references
	props.Set("secret_references", Map{
		"type": "array",
		"items": Map{
			"type": "object",
			"properties": Map{
				"id": Map{
					"type": "string",
				},
			},
		},
	})

	// Set new_package_policy.inputs to a map
	inputs := schema.Components.MustGetMap("schemas.new_package_policy.properties.inputs")
	inputs.Set("type", "object")
	inputs.Move("items", "additionalProperties")

	// Set new_package_policy.inputs.streams to a map
	streams := schema.Components.MustGetMap("schemas.new_package_policy.properties.inputs.additionalProperties.properties.streams")
	streams.Set("type", "object")
	streams.Move("items", "additionalProperties")

	// Define new_package_policy.inputs.streams.enabled/vars
	streams.Set("additionalProperties", map[string]any{
		"properties": map[string]any{
			"enabled": map[string]any{"type": "boolean"},
			"vars":    map[string]any{"type": "object"},
		},
	})

	// For all output_request variants, make "type" a required field
	for _, v := range []string{
		"schemas.output_create_request_elasticsearch.required",
		"schemas.output_create_request_kafka.required",
		"schemas.output_create_request_logstash.required",
		"schemas.output_create_request_remote_elasticsearch.required",
		"schemas.output_update_request_elasticsearch.required",
		"schemas.output_update_request_kafka.required",
		"schemas.output_update_request_logstash.required",
		// "schemas.output_update_request_remote_elasticsearch.required", // MISSING
	} {
		slice := schema.Components.MustGetSlice(v)
		if slice.Contains("type") {
			continue
		}
		slice = append(slice, "type")
		schema.Components.Set(v, slice)
	}
}

// transformDataViewPaths fixes the data_views paths.
func transformDataViewPaths(schema *Schema) {
	path := schema.MustGetPath("/api/data_views/data_view")
	path.Post.Set("operationId", "createDataViewDefault") // has a trailing "w"
}

// transformAlertingPaths fixes the alerting paths.
func transformAlertingPaths(schema *Schema) {
	ruleIdPath := schema.MustGetPath("/api/alerting/rule/{id}")
	ruleFindPath := schema.MustGetPath("/api/alerting/rules/_find")

	// Fix this: "%2Fapi%2Falerting%2Frule%2F%7Bid%7D#0Default"
	ruleIdPath.Get.Set("operationId", "getRuleDefault")
	ruleIdPath.Post.Set("operationId", "createRuleDefault")
	ruleIdPath.Put.Set("operationId", "updateRuleDefault")
	ruleIdPath.Delete.Set("operationId", "deleteRuleDefault")
	ruleFindPath.Get.Set("operationId", "findRuleDefault")

	// When a []string query param is combined with x-go-type-skip-optional-pointers,
	// the generated function still tries to deref the array, which is no longer possible.
	// Make it required, so the resulting function works with a []string instead of a *[]string.
	for _, _param := range ruleFindPath.Get.MustGetSlice("parameters") {
		param := Map(_param.(map[string]any))
		if _name, ok := param.Get("name"); ok {
			name := _name.(string)
			if name == "fields" || name == "filter_consumers" {
				param.Set("required", true)
			}
		}
	}
}

// transformConnectorPaths fixes the connector paths.
func transformConnectorPaths(schema *Schema) {
	connectorIdPath := schema.MustGetPath("/api/actions/connector/{id}")
	connectorsPath := schema.MustGetPath("/api/actions/connectors")

	// Fix this: "%2Fapi%2Factions%2Fconnectors#0Default"
	connectorIdPath.Get.Set("operationId", "getConnectorDefault")
	connectorIdPath.Post.Set("operationId", "createConnectorDefault")
	connectorIdPath.Put.Set("operationId", "updateConnectorDefault")
	connectorIdPath.Delete.Set("operationId", "deleteConnectorDefault")
	connectorsPath.Get.Set("operationId", "getAllConnectorsDefault")
}

// transformExtractAndSimplifyNested does several things.
// - Remove all non-200 error codes, the raw error response is only ever used.
// - Create a ref for each JSON 200 request body and response.
// - Delete any "example" or "examples" keys to slim the output.
// - Delete any enums, validate in the resource or server-side instead.
// - Add "x-go-type-skip-optional-pointer" to any arrays or objects without properties.
func transformExtractAndSimplifyNested(schema *Schema) {
	deleteExampleFn := func(key string, node Map) {
		if node.Has("example") {
			node.Delete("example")
		}
		if node.Has("examples") {
			node.Delete("examples")
		}
	}
	deleteEnumFn := func(key string, node Map) {
		if enum, ok := node.GetSlice("enum"); ok {
			node.Delete("enum")
			// Some enums were missing a type
			if !node.Has("type") {
				switch enum[0].(type) {
				case string:
					node.Set("type", "string")
				case int:
					node.Set("type", "integer")
				default:
					log.Panicf("Unhandled enum type: %T", enum[0])
				}
			}
		}
	}
	removeGoPointersFn := func(key string, node Map) {
		if vtype, ok := node["type"]; ok {
			switch vtype {
			case "array":
				node["x-go-type-skip-optional-pointer"] = true
			case "object":
				if !node.Has("properties") {
					node["x-go-type-skip-optional-pointer"] = true
				}
			}
		}
	}

	for _, pathInfo := range schema.Paths {
		for _, methInfo := range pathInfo.Endpoints() {
			responses := methInfo.MustGetMap("responses")
			for code := range responses {
				// Delete everything but 200 or 204.
				if code != "200" && code != "204" {
					responses.Delete(code)
					continue
				}
			}

			// Ref each method's JSON 200 body and response
			refName := methInfo.MustGet("operationId").(string)
			refName = strings.TrimSuffix(refName, "Default")
			refName = toSnake(refName)
			if req, ok := methInfo.GetMap("requestBody.content.application/json"); ok {
				if req.Has("schema") && !req.Has("schema.$ref") {
					req.CreateRef(refName+"_request", "schema", schema)
				}
			}
			if resp, ok := methInfo.GetMap("responses.200.content.application/json"); ok {
				if resp.Has("schema") && !resp.Has("schema.$ref") {
					resp.CreateRef(refName+"_response_object", "schema", schema)
				}
			}

			methInfo.Iterate(deleteExampleFn)    // Delete any examples
			methInfo.Iterate(deleteEnumFn)       // Delete any enums
			methInfo.Iterate(removeGoPointersFn) // Convert any *[] or *map objects
		}
	}

	schema.Components.Set("examples", Map{})      // Empty components.examples
	schema.Components.Iterate(deleteExampleFn)    // Delete any examples
	schema.Components.Iterate(deleteEnumFn)       // Delete any enums
	schema.Components.Iterate(removeGoPointersFn) // Convert any *[] or *map objects

	// Clean up additionalProperties = false or {}
	schema.Components.Iterate(func(key string, node Map) {
		if v, ok := node["additionalProperties"]; ok {
			switch vv := v.(type) {
			case bool:
				node.Delete("additionalProperties")
			case map[string]any:
				if len(vv) == 0 {
					node.Delete("additionalProperties")
				}
			}
		}
	})

	// Extract everything in component schemas
	componentSchemas := schema.Components.MustGetMap("schemas")
	componentSchemas.Iterate(func(key string, node Map) {
		// Top level
		if !strings.Contains(key, ".") {
			return
		}

		if node.Has("properties") {
			refName := key
			// Edge case(s)
			if refName == "full_agent_policy.properties.output_permissions.additionalProperties" {
				refName = "full_agent_policy_output_permission_item"
			} else {
				refName = strings.ReplaceAll(refName, ".properties", "")
				refName = strings.ReplaceAll(refName, ".additionalProperties", ".item")
				refName = strings.ReplaceAll(refName, ".items", ".item")
				refName = strings.ReplaceAll(refName, ".", "_")
			}
			componentSchemas.CreateRef(refName, key, schema)
		}
	})

	// Dedup component schemas
	/*
		{
			// Collect all refValues
			migrations := make(map[string]string)

			for _, k1 := range componentSchemas.Keys() {
				v1, ok := componentSchemas.GetMap(k1)
				if !ok {
					continue // Deleted
				}

				for _, k2 := range componentSchemas.Keys() {
					v2, ok := componentSchemas.GetMap(k2)
					if !ok {
						continue // Deleted
					}

					if k1 != k2 && reflect.DeepEqual(v1, v2) {
						newRef := "#/components/schemas/" + k1
						oldRef := "#/components/schemas/" + k2
						migrations[oldRef] = newRef
						componentSchemas.Delete(k2)
					}
				}
			}

			updateRefs := func(key string, node Map) {
				if r, ok := node["$ref"]; ok {
					ref := r.(string)
					if newRef, ok := migrations[ref]; ok {
						node["$ref"] = newRef
					}
				}
			}

			for _, pathInfo := range schema.Paths {
				for _, methodInfo := range pathInfo.Endpoints() {
					methodInfo.Iterate(updateRefs)
				}
			}

			schema.Components.Iterate(updateRefs)
		}
	*/
}
