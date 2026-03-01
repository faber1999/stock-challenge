package httpapi

import (
	"fmt"
	"net/http"
	"strings"
)

const openAPITemplate = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Stock Challenge Backend API",
    "description": "API del challenge para sincronizar, consultar y recomendar stocks.",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "%s"
    }
  ],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "operationId": "health",
        "responses": {
          "200": {
            "description": "Servicio saludable",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {
                      "type": "string",
                      "example": "ok"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/db-time": {
      "get": {
        "summary": "Hora actual de la base de datos",
        "operationId": "dbTime",
        "responses": {
          "200": {
            "description": "Hora actual en UTC",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "now": {
                      "type": "string",
                      "format": "date-time",
                      "example": "2026-01-01T12:00:00Z"
                    }
                  }
                }
              }
            }
          },
          "500": {
            "description": "Error consultando la base de datos",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string",
                      "example": "failed to query database"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/stocks/sync": {
      "post": {
        "summary": "Sincroniza stocks desde la API externa",
        "operationId": "syncStocks",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "Maximo de paginas a consultar",
            "schema": {
              "type": "integer",
              "default": 10
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Sincronizacion completa",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SyncResult"
                }
              }
            }
          },
          "502": {
            "description": "Error de API externa"
          }
        }
      }
    },
    "/stocks": {
      "get": {
        "summary": "Lista stocks con filtros y orden",
        "operationId": "listStocks",
        "parameters": [
          {
            "name": "q",
            "in": "query",
            "description": "Busqueda por ticker/company/brokerage",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "action",
            "in": "query",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "sort_by",
            "in": "query",
            "schema": {
              "type": "string",
              "enum": [
                "ticker",
                "company",
                "brokerage",
                "action",
                "rating_from",
                "rating_to",
                "target_from",
                "target_to",
                "synced_at"
              ]
            }
          },
          {
            "name": "order",
            "in": "query",
            "schema": {
              "type": "string",
              "enum": [
                "asc",
                "desc"
              ]
            }
          },
          {
            "name": "limit",
            "in": "query",
            "schema": {
              "type": "integer",
              "default": 20
            }
          },
          {
            "name": "offset",
            "in": "query",
            "schema": {
              "type": "integer",
              "default": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Listado de stocks",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ListStocksResponse"
                }
              }
            }
          }
        }
      }
    },
    "/stocks/{ticker}": {
      "get": {
        "summary": "Detalle por ticker",
        "operationId": "getStockByTicker",
        "parameters": [
          {
            "name": "ticker",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Stock encontrado",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Stock"
                }
              }
            }
          },
          "404": {
            "description": "No encontrado"
          }
        }
      }
    },
    "/stocks/recommendations": {
      "get": {
        "summary": "Top recomendaciones de inversion",
        "operationId": "recommendStocks",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "schema": {
              "type": "integer",
              "default": 5
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Recomendaciones",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/RecommendationsResponse"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Stock": {
        "type": "object",
        "properties": {
          "ticker": {
            "type": "string",
            "example": "AAPL"
          },
          "id": {
            "type": "integer",
            "format": "int64"
          },
          "company": {
            "type": "string",
            "example": "Apple Inc."
          },
          "brokerage": {
            "type": "string",
            "example": "Goldman Sachs"
          },
          "action": {
            "type": "string",
            "example": "Buy"
          },
          "rating_from": {
            "type": "string",
            "example": "Hold"
          },
          "rating_to": {
            "type": "string",
            "example": "Buy"
          },
          "target_from": {
            "type": "number",
            "format": "double"
          },
          "target_to": {
            "type": "number",
            "format": "double"
          },
          "currency": {
            "type": "string",
            "example": "USD"
          },
          "recommend_score": {
            "type": "number",
            "format": "double"
          },
          "synced_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "ticker",
          "company",
          "brokerage",
          "action",
          "rating_from",
          "rating_to",
          "target_from",
          "target_to",
          "currency",
          "recommend_score",
          "synced_at"
        ]
      },
      "SyncResult": {
        "type": "object",
        "properties": {
          "pages_processed": {
            "type": "integer"
          },
          "stocks_saved": {
            "type": "integer"
          }
        }
      },
      "ListStocksResponse": {
        "type": "object",
        "properties": {
          "items": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/Stock"
            }
          },
          "total": {
            "type": "integer"
          },
          "limit": {
            "type": "integer"
          },
          "offset": {
            "type": "integer"
          }
        }
      },
      "Recommendation": {
        "allOf": [
          {
            "$ref": "#/components/schemas/Stock"
          },
          {
            "type": "object",
            "properties": {
              "score": {
                "type": "number",
                "format": "double"
              },
              "upside_pct": {
                "type": "number",
                "format": "double"
              }
            }
          }
        ]
      },
      "RecommendationsResponse": {
        "type": "object",
        "properties": {
          "items": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/Recommendation"
            }
          },
          "total": {
            "type": "integer"
          }
        }
      }
    }
  }
}`

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/swagger/doc.json',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`

func (h *Handler) swaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/swagger" && r.URL.Path != "/swagger/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerUIHTML))
}

func (h *Handler) swaggerDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf(openAPITemplate, serverURLFromRequest(r))))
}

func serverURLFromRequest(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = strings.TrimSpace(strings.Split(forwardedProto, ",")[0])
	}

	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = strings.TrimSpace(strings.Split(forwardedHost, ",")[0])
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
