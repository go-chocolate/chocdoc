package openapi

import (
	"encoding/json"
	"strings"
)

type (
	Type              string
	Format            string
	Method            string
	ParameterLocation string
)

const (
	TypeUnknown Type = "unknown"
	TypeNull    Type = "null"
	TypeBoolean Type = "boolean"
	TypeObject  Type = "object"
	TypeArray   Type = "array"
	TypeNumber  Type = "number"
	TypeString  Type = "string"
	TypeInteger Type = "integer"

	FormatBinary  Format = "binary"
	MethodGet     Method = "get"
	MethodHead    Method = "head"
	MethodPost    Method = "post"
	MethodPut     Method = "put"
	MethodPatch   Method = "patch"
	MethodDelete  Method = "delete"
	MethodConnect Method = "connect"
	MethodOptions Method = "options"
	MethodTrace   Method = "trace"

	InPath   ParameterLocation = "path"
	InQuery  ParameterLocation = "query"
	InHeader ParameterLocation = "header"
	InCookie ParameterLocation = "cookie"
)

func (m Method) String() string {
	return strings.ToLower(string(m))
}

type Information struct {
	//REQUIRED. The title of the API
	Title string `json:"title"`

	//REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `json:"version"`

	//A description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	//License information for the exposed API.
	License *License `json:"license,omitempty"`
}

type License struct {
	//REQUIRED. The license name used for the API.
	Name string `json:"name"`

	//A URL to the license used for the API. MUST be in the format of a URL.
	Url string `json:"url"`
}

type Schema struct {
	//The value of this keyword MUST be either a string or an array. If it is an array, elements of the array MUST be strings and MUST be unique.
	//String values MUST be one of the six primitive types ("null", "boolean", "object", "array", "number", or "string"), or "integer" which matches any number with a zero fractional part.
	//An instance validates if and only if the instance is in any of the sets listed for this keyword.
	Type Type `json:"type,omitempty"`

	//Structural validation alone may be insufficient to validate that an instance meets all the requirements of an application. The "format" keyword is defined to allow interoperable semantic validation for a fixed subset of values which are accurately described by authoritative resources, be they RFCs or other external specifications.
	//The value of this keyword is called a format attribute. It MUST be a string. A format attribute can generally only validate a given set of instance types. If the type of the instance to validate is not in this set, validation for this format attribute and instance SHOULD succeed.
	Format Format `json:"format,omitempty"`

	//The value of "properties" MUST be an object. Each value of this object MUST be a valid JSON Schema.
	//This keyword determines how child instances validate for objects, and does not directly validate the immediate instance itself.
	//Validation succeeds if, for each name that appears in both the instance and as a name within this keyword's value, the child instance for that name successfully validates against the corresponding schema.
	//Omitting this keyword has the same behavior as an empty object.
	Properties map[string]*Schema `json:"properties,omitempty"`

	//The value of "items" MUST be either a valid JSON Schema or an array of valid JSON Schemas.
	//This keyword determines how child instances validate for arrays, and does not directly validate the immediate instance itself.
	//If "items" is a schema, validation succeeds if all elements in the array successfully validate against that schema.
	//If "items" is an array of schemas, validation succeeds if each element of the instance validates against the schema at the same position, if any.
	//Omitting this keyword has the same behavior as an empty schema.
	Items *Schema `json:"items,omitempty"`
}

func (s *Schema) String() string {
	return s.JSON()
}

func (s *Schema) JSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Schema) Copy() *Schema {
	var copied = new(Schema)
	*copied = *s
	if s.Properties != nil {
		copied.Properties = map[string]*Schema{}
	}
	if s.Items != nil {
		copied.Items = s.Items.Copy()
	}
	for k, v := range s.Properties {
		copied.Properties[k] = v.Copy()
	}
	return copied
}

func (s *Schema) Extend(schema *Schema, field string) bool {
	switch s.Type {
	case TypeObject:
		if s.Properties[field] != nil {
			s.Properties[field] = schema
			return true
		}
		for _, prop := range s.Properties {
			if prop.Extend(schema, field) {
				return true
			}
		}
	case TypeArray:
		if s.Items.Extend(schema, field) {
			return true
		}
	}
	return false
}

type Content struct {
	//JSON Schema Validation: A Vocabulary for Structural Validation of JSON
	Schema *Schema `json:"schema,omitempty"`
}

type Body struct {
	//REQUIRED. A description of the response. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description"`
	//A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]*Content `json:"content,omitempty"`
}

type Parameter struct {
	//REQUIRED. The name of the parameter. Parameter names are case sensitive.
	//If in is "path", the name field MUST correspond to a template expression occurring within the path field in the Paths Object. See Path Templating for further information.
	//If in is "header" and the name field is "Accept", "Content-Type" or "Authorization", the parameter definition SHALL be ignored.
	//For all other cases, the name corresponds to the parameter name used by the in property.
	Name string `json:"name"`

	//REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	In ParameterLocation `json:"in"`

	//A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	//Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Required bool `json:"required,omitempty"`

	//JSON Schema Validation: A Vocabulary for Structural Validation of JSON
	Schema *Schema `json:"schema,omitempty"`
}

type API struct {
	//A short summary of what the operation does. For maximum readability in the swagger-ui, this field SHOULD be less than 120 characters.
	Summary string `json:"summary"`

	//A verbose explanation of the operation behavior. GFM syntax can be used for rich text representation.
	Description string `json:"description"`

	//A list of MIME types the operation can consume. This overrides the consumes definition at the Swagger Object. An empty value MAY be used to clear the global definition. Value MUST be as described under Mime Types.
	Consumes []string `json:"consumes,omitempty"`

	//A list of MIME types the operation can produce. This overrides the produces definition at the Swagger Object. An empty value MAY be used to clear the global definition. Value MUST be as described under Mime Types.
	Produces []string `json:"produces,omitempty"`

	//A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Tags []string `json:"tags,omitempty"`

	//A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it, but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the Swagger Object's parameters. There can be one "body" parameter at most.
	Parameters []*Parameter `json:"parameters,omitempty"`

	//The request body applicable for this operation. The requestBody is only supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague, requestBody SHALL be ignored by consumers.
	RequestBody *Body `json:"requestBody,omitempty"`

	//Required. The list of possible responses as they are returned from executing this operation.
	Responses map[string]*Body `json:"responses"`
}

type Path map[Method]*API

type Swagger struct {
	//REQUIRED. This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses. The openapi field SHOULD be used by tooling specifications and clients to interpret the OpenAPI document. This is not related to the API info.version string
	Openapi string `json:"openapi"`

	//The object provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
	Info Information `json:"info"`

	//Holds the relative paths to the individual endpoints and their operations. The path is appended to the URL from the Server Object in order to construct the full URL. The Paths MAY be empty, due to ACL constraints.
	Paths map[string]Path `json:"paths,omitempty"`
}

func (s *Swagger) JSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}
