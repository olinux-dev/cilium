// Code generated by go-swagger; DO NOT EDIT.

// Copyright Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Cilium",
    "title": "Cilium Operator",
    "version": "v1beta"
  },
  "basePath": "/v1",
  "paths": {
    "/healthz": {
      "get": {
        "description": "This path will return the status of cilium operator instance.",
        "produces": [
          "text/plain"
        ],
        "tags": [
          "operator"
        ],
        "summary": "Get health of Cilium operator",
        "responses": {
          "200": {
            "description": "Cilium operator is healthy",
            "schema": {
              "type": "string"
            }
          },
          "500": {
            "description": "Cilium operator is not healthy",
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/metrics/": {
      "get": {
        "tags": [
          "metrics"
        ],
        "summary": "Retrieve cilium operator metrics",
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "../openapi.yaml#/definitions/Metric"
              }
            }
          },
          "500": {
            "description": "Metrics cannot be retrieved",
            "x-go-name": "Failed"
          }
        }
      }
    }
  },
  "x-schemes": [
    "unix"
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Cilium",
    "title": "Cilium Operator",
    "version": "v1beta"
  },
  "basePath": "/v1",
  "paths": {
    "/healthz": {
      "get": {
        "description": "This path will return the status of cilium operator instance.",
        "produces": [
          "text/plain"
        ],
        "tags": [
          "operator"
        ],
        "summary": "Get health of Cilium operator",
        "responses": {
          "200": {
            "description": "Cilium operator is healthy",
            "schema": {
              "type": "string"
            }
          },
          "500": {
            "description": "Cilium operator is not healthy",
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/metrics/": {
      "get": {
        "tags": [
          "metrics"
        ],
        "summary": "Retrieve cilium operator metrics",
        "responses": {
          "200": {
            "description": "Success",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/metric"
              }
            }
          },
          "500": {
            "description": "Metrics cannot be retrieved",
            "x-go-name": "Failed"
          }
        }
      }
    }
  },
  "definitions": {
    "metric": {
      "description": "Metric information",
      "type": "object",
      "properties": {
        "labels": {
          "description": "Labels of the metric",
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "name": {
          "description": "Name of the metric",
          "type": "string"
        },
        "value": {
          "description": "Value of the metric",
          "type": "number"
        }
      }
    }
  },
  "x-schemes": [
    "unix"
  ]
}`))
}