package chat_completions

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestConvertOpenAIRequestToGeminiCLI_MapsResponseFormatJSONSchema(t *testing.T) {
	input := []byte(`{
		"messages":[{"role":"user","content":"Return JSON"}],
		"response_format":{
			"type":"json_schema",
			"json_schema":{
				"name":"result",
				"strict":true,
				"schema":{
					"type":"object",
					"properties":{"x":{"type":"string"}},
					"required":["x"],
					"additionalProperties":false
				}
			}
		}
	}`)

	output := ConvertOpenAIRequestToGeminiCLI("gemini-2.5-flash", input, false)
	raw := string(output)

	if gjson.Get(raw, "request.generationConfig.responseMimeType").String() != "application/json" {
		t.Fatalf("expected responseMimeType=application/json, got %s", gjson.Get(raw, "request.generationConfig.responseMimeType").Raw)
	}
	if !gjson.Get(raw, "request.generationConfig.responseJsonSchema").Exists() {
		t.Fatalf("expected responseJsonSchema to exist, got %s", gjson.Get(raw, "request.generationConfig").Raw)
	}
	if gjson.Get(raw, "request.generationConfig.responseJsonSchema.additionalProperties").Exists() {
		t.Fatalf("expected Gemini-cleaned schema without additionalProperties, got %s", gjson.Get(raw, "request.generationConfig.responseJsonSchema").Raw)
	}
}

func TestConvertOpenAIRequestToGeminiCLI_MapsResponseFormatJSONObject(t *testing.T) {
	input := []byte(`{
		"messages":[{"role":"user","content":"Return JSON"}],
		"response_format":{"type":"json_object"}
	}`)

	output := ConvertOpenAIRequestToGeminiCLI("gemini-2.5-flash", input, false)
	raw := string(output)

	if gjson.Get(raw, "request.generationConfig.responseMimeType").String() != "application/json" {
		t.Fatalf("expected responseMimeType=application/json, got %s", gjson.Get(raw, "request.generationConfig.responseMimeType").Raw)
	}
}
