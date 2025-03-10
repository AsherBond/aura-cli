package allowedorigin_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

const (
	instanceId    = "2f49c2b3"
	dataApiId     = "e157301d"
	allowedOrigin = "https://test.com"

	mockPatchResponse = `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
		}
	}`
)

func TestAddAllowedOriginFlagsValidation(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	tests := map[string]struct {
		executedCommand string
		expectedError   string
	}{
		"missing all flags": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s", allowedOrigin),
			expectedError:   "Error: required flag(s) \"data-api-id\", \"instance-id\" not set",
		},
		"missing origin": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin add --data-api-id %s --instance-id %s", dataApiId, instanceId),
			expectedError:   "Error: accepts 1 arg(s), received 0",
		},
		"missing data api id flag": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s --instance-id %s", allowedOrigin, instanceId),
			expectedError:   "Error: required flag(s) \"data-api-id\" not set",
		},
		"missing instance id flag": {
			executedCommand: fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s --data-api-id %s", allowedOrigin, dataApiId),
			expectedError:   "Error: required flag(s) \"instance-id\" not set",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			helper.ExecuteCommand(tt.executedCommand)
			helper.AssertErr(tt.expectedError)
		})
	}
}

func TestAddAllowedOriginWithNoExistingOrigins(t *testing.T) {
	mockGetResponse := `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": []
				}
			}
		}
	}`
	expectedResponse := fmt.Sprintf(`New allowed origins: ["%s"]
{
	"data": {
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "ready",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`, allowedOrigin)

	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusOK, mockGetResponse)
	mockHandler.AddResponse(http.StatusAccepted, mockPatchResponse)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s --instance-id %s --data-api-id %s", allowedOrigin, instanceId, dataApiId))

	mockHandler.AssertCalledTimes(2)
	mockHandler.AssertCalledWithMethod(http.MethodGet)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(fmt.Sprintf("{\"security\":{\"cors_policy\":{\"allowed_origins\":[\"%s\"]}}}", allowedOrigin))

	helper.AssertOut(expectedResponse)
}

func TestAddAllowedOriginWithExistingOrigins(t *testing.T) {
	mockGetResponse := `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": ["https://test1.com", "https://test2.com"]
				}
			}
		}
	}`

	expectedResponse := fmt.Sprintf(`New allowed origins: ["https://test1.com", "https://test2.com", "%s"]
{
	"data": {
		"id": "2f49c2b3",
		"name": "my-data-api-1",
		"status": "ready",
		"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql"
	}
}`, allowedOrigin)

	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusOK, mockGetResponse)
	mockHandler.AddResponse(http.StatusAccepted, mockPatchResponse)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s --instance-id %s --data-api-id %s", allowedOrigin, instanceId, dataApiId))

	mockHandler.AssertCalledTimes(2)
	mockHandler.AssertCalledWithMethod(http.MethodGet)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(fmt.Sprintf("{\"security\":{\"cors_policy\":{\"allowed_origins\":[\"https://test1.com\",\"https://test2.com\",\"%s\"]}}}", allowedOrigin))

	helper.AssertOut(expectedResponse)
}

func TestAddAllowedOriginWithDuplicateOrigin(t *testing.T) {
	mockGetResponse := fmt.Sprintf(`{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": ["https://test1.com", "%s", "https://test2.com"]
				}
			}
		}
	}`, allowedOrigin)

	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusOK, mockGetResponse)
	mockHandler.AddResponse(http.StatusAccepted, mockPatchResponse)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s --instance-id %s --data-api-id %s", allowedOrigin, instanceId, dataApiId))

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertErr(fmt.Sprintf("Error: Origin \"%s\" already exists in allowed origins\n", allowedOrigin))
}

func TestAddAllowedOriginWithOutputTable(t *testing.T) {
	mockGetResponse := `{
		"data": {
			"id": "2f49c2b3",
			"name": "my-data-api-1",
			"status": "ready",
			"url": "https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql",
			"security": {
				"cors_policy": {
					"allowed_origins": ["https://test1.com", "https://test2.com"]
				}
			}
		}
	}`
	expectedResponse := fmt.Sprintf(`New allowed origins: ["https://test1.com", "https://test2.com", "%s"]
┌──────────┬───────────────┬────────┬────────────────────────────────────────────────────────────────────────────────┐
│ ID       │ NAME          │ STATUS │ URL                                                                            │
├──────────┼───────────────┼────────┼────────────────────────────────────────────────────────────────────────────────┤
│ 2f49c2b3 │ my-data-api-1 │ ready  │ https://2f49c2b3.28be6e4d8d3e8360197cb6c1fa1d25d1.graphql.neo4j-dev.io/graphql │
└──────────┴───────────────┴────────┴────────────────────────────────────────────────────────────────────────────────┘
`, allowedOrigin)

	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)

	mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v1/instances/%s/data-apis/graphql/%s", instanceId, dataApiId), http.StatusOK, mockGetResponse)
	mockHandler.AddResponse(http.StatusAccepted, mockPatchResponse)

	helper.ExecuteCommand(fmt.Sprintf("data-api graphql cors-policy allowed-origin add %s --instance-id %s --data-api-id %s --output table", allowedOrigin, instanceId, dataApiId))

	mockHandler.AssertCalledTimes(2)
	mockHandler.AssertCalledWithMethod(http.MethodGet)
	mockHandler.AssertCalledWithMethod(http.MethodPatch)
	mockHandler.AssertCalledWithBody(fmt.Sprintf("{\"security\":{\"cors_policy\":{\"allowed_origins\":[\"https://test1.com\",\"https://test2.com\",\"%s\"]}}}", allowedOrigin))

	helper.AssertOut(expectedResponse)
}
