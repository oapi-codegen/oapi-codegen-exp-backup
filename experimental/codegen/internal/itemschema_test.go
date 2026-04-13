package codegen

import (
	"strings"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func schemaPathStrings(schemas []*SchemaDescriptor) []string {
	paths := make([]string, len(schemas))
	for i, s := range schemas {
		paths[i] = s.Path.String()
	}
	return paths
}

func TestIsSequentialMediaType(t *testing.T) {
	tests := []struct {
		contentType string
		want        bool
	}{
		{"text/event-stream", true},
		{"application/jsonl", true},
		{"application/x-ndjson", true},
		{"application/json-seq", true},
		{"multipart/mixed", true},
		{"application/json", false},
		{"text/plain", false},
		{"application/xml", false},
		{"multipart/form-data", false},
		{"application/x-www-form-urlencoded", false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			got := IsSequentialMediaType(tt.contentType)
			assert.Equal(t, tt.want, got)
		})
	}
}

const itemSchemaResponseSpec = `openapi: "3.1.0"
info:
  title: ItemSchema Test API
  version: "1.0"
paths:
  /events:
    get:
      operationId: streamEvents
      responses:
        "200":
          description: OK
          content:
            text/event-stream:
              schema:
                type: string
                description: Raw SSE stream
              itemSchema:
                type: object
                properties:
                  id:
                    type: integer
                  message:
                    type: string
`

const itemSchemaRequestBodySpec = `openapi: "3.1.0"
info:
  title: ItemSchema RequestBody Test API
  version: "1.0"
paths:
  /ingest:
    post:
      operationId: ingestEvents
      requestBody:
        required: true
        content:
          application/jsonl:
            schema:
              type: string
              description: Raw JSONL stream
            itemSchema:
              type: object
              properties:
                event:
                  type: string
                timestamp:
                  type: integer
      responses:
        "202":
          description: Accepted
`

func TestGatherItemSchema_Response(t *testing.T) {
	doc, err := libopenapi.NewDocument([]byte(itemSchemaResponseSpec))
	require.NoError(t, err, "Failed to parse spec")

	contentTypeMatcher := NewContentTypeMatcher(DefaultContentTypes())
	schemas, err := GatherSchemas(doc, contentTypeMatcher, OutputOptions{})
	require.NoError(t, err, "Failed to gather schemas")

	// Verify the itemSchema was gathered
	var foundItemSchema bool
	for _, s := range schemas {
		pathStr := s.Path.String()
		if strings.Contains(pathStr, "itemSchema") {
			foundItemSchema = true
			assert.NotNil(t, s.Schema)
			assert.NotNil(t, s.Schema.Properties)
			break
		}
	}
	assert.True(t, foundItemSchema, "expected itemSchema to be gathered from response; gathered paths: %v", schemaPathStrings(schemas))
}

func TestGatherItemSchema_RequestBody(t *testing.T) {
	doc, err := libopenapi.NewDocument([]byte(itemSchemaRequestBodySpec))
	require.NoError(t, err, "Failed to parse spec")

	contentTypeMatcher := NewContentTypeMatcher(DefaultContentTypes())
	schemas, err := GatherSchemas(doc, contentTypeMatcher, OutputOptions{})
	require.NoError(t, err, "Failed to gather schemas")

	// Verify the itemSchema was gathered
	var foundItemSchema bool
	for _, s := range schemas {
		pathStr := s.Path.String()
		if strings.Contains(pathStr, "itemSchema") {
			foundItemSchema = true
			assert.NotNil(t, s.Schema)
			assert.NotNil(t, s.Schema.Properties)
			break
		}
	}
	assert.True(t, foundItemSchema, "expected itemSchema to be gathered from request body; gathered paths: %v", schemaPathStrings(schemas))
}

func TestGatherOperations_ItemSchema_Response(t *testing.T) {
	doc, err := libopenapi.NewDocument([]byte(itemSchemaResponseSpec))
	require.NoError(t, err, "Failed to parse spec")

	ctx := NewCodegenContext()
	contentTypeMatcher := NewContentTypeMatcher(DefaultContentTypes())
	ops, err := GatherOperations(doc, ctx, contentTypeMatcher)
	require.NoError(t, err, "Failed to gather operations")
	require.Len(t, ops, 1, "Expected 1 operation")

	op := ops[0]
	assert.Equal(t, "streamEvents", op.OperationID)

	// Find the SSE response content
	require.NotEmpty(t, op.Responses, "Operation responses should not be empty")
	resp := op.Responses[0]
	require.NotEmpty(t, resp.Contents, "Response content should not be empty")

	var sseContent *ResponseContentDescriptor
	for _, c := range resp.Contents {
		if c.ContentType == "text/event-stream" {
			sseContent = c
			break
		}
	}
	require.NotNil(t, sseContent, "expected text/event-stream content")

	assert.True(t, sseContent.IsSequential)
	assert.NotNil(t, sseContent.ItemSchema, "expected ItemSchema to be populated")
	assert.NotNil(t, sseContent.ItemSchema.Schema)
	assert.NotNil(t, sseContent.ItemSchema.Schema.Properties)
}

func TestGatherOperations_ItemSchema_RequestBody(t *testing.T) {
	doc, err := libopenapi.NewDocument([]byte(itemSchemaRequestBodySpec))
	require.NoError(t, err, "Failed to parse spec")

	ctx := NewCodegenContext()
	contentTypeMatcher := NewContentTypeMatcher(DefaultContentTypes())
	ops, err := GatherOperations(doc, ctx, contentTypeMatcher)
	require.NoError(t, err, "Failed to gather operations")
	require.Len(t, ops, 1, "Expected 1 operation")

	op := ops[0]
	assert.Equal(t, "ingestEvents", op.OperationID)

	// Find the JSONL request body
	require.NotEmpty(t, op.Bodies, "Operation body should not be empty")

	var jsonlBody *RequestBodyDescriptor
	for _, b := range op.Bodies {
		if b.ContentType == "application/jsonl" {
			jsonlBody = b
			break
		}
	}
	require.NotNil(t, jsonlBody, "expected application/jsonl body")

	assert.True(t, jsonlBody.IsSequential)
	assert.NotNil(t, jsonlBody.ItemSchema, "expected ItemSchema to be populated")
	assert.NotNil(t, jsonlBody.ItemSchema.Schema)
	assert.NotNil(t, jsonlBody.ItemSchema.Schema.Properties)
}

func TestGatherOperations_NonSequential_HasNoItemSchema(t *testing.T) {
	const spec = `openapi: "3.1.0"
info:
  title: Test
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    name:
                      type: string
`
	doc, err := libopenapi.NewDocument([]byte(spec))
	require.NoError(t, err, "Failed to parse spec")

	ctx := NewCodegenContext()
	contentTypeMatcher := NewContentTypeMatcher(DefaultContentTypes())
	ops, err := GatherOperations(doc, ctx, contentTypeMatcher)
	require.NoError(t, err, "Failed to gather operations")
	require.Len(t, ops, 1, "Expected 1 operation")

	resp := ops[0].Responses[0]
	require.NotEmpty(t, resp.Contents)

	jsonContent := resp.Contents[0]
	assert.False(t, jsonContent.IsSequential, "Contents should not be sequential")
	assert.Nil(t, jsonContent.ItemSchema, "Item Schema should be empty")
}

func TestGenerate_ItemSchema(t *testing.T) {
	const spec = `openapi: "3.1.0"
info:
  title: Streaming API
  version: "1.0"
paths:
  /events:
    get:
      operationId: streamEvents
      responses:
        "200":
          description: OK
          content:
            text/event-stream:
              schema:
                type: string
              itemSchema:
                type: object
                properties:
                  id:
                    type: integer
                  message:
                    type: string
  /ingest:
    post:
      operationId: ingestData
      requestBody:
        required: true
        content:
          application/x-ndjson:
            schema:
              type: string
            itemSchema:
              $ref: '#/components/schemas/Event'
      responses:
        "202":
          description: Accepted
components:
  schemas:
    Event:
      type: object
      properties:
        type:
          type: string
        data:
          type: string
`

	doc, err := libopenapi.NewDocument([]byte(spec))
	require.NoError(t, err, "Failed to parse spec")

	cfg := Configuration{
		PackageName: "testpkg",
	}

	code, err := Generate(doc, []byte(spec), cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, code)

	// Verify the itemSchema types are generated
	assert.Contains(t, code, "Event")
}
