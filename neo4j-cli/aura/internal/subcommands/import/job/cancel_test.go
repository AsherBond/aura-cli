package job_test

import (
	"fmt"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
	"net/http"
	"testing"
)

func TestCancelImportJob(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs/87d485b4-73fc-4a7f-bb03-720f4672947e/cancellation", http.StatusOK, `
		{
			"data": {"id": "87d485b4-73fc-4a7f-bb03-720f4672947e"}
		}
	`)

	helper.SetConfigValue("aura.beta-enabled", true)

	helper.ExecuteCommand("import job cancel --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee 87d485b4-73fc-4a7f-bb03-720f4672947e")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodPost)

	helper.AssertErr("")
	helper.AssertOutJson(`
		{
			"data": {"id": "87d485b4-73fc-4a7f-bb03-720f4672947e"}
		}
	`)
}

func TestCancelImportJobError(t *testing.T) {
	projectId := "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
	jobId := "87d485b4-73fc-4a7f-bb03-720f4672947e"
	testCases := map[string]struct {
		executeCommand      string
		statusCode          int
		expectedCalledTimes int
		expectedError       string
		returnBody          string
	}{
		"correct command with error response 1": {
			executeCommand:      fmt.Sprintf("import job cancel --project-id=%s %s", projectId, jobId),
			statusCode:          http.StatusBadRequest,
			expectedCalledTimes: 1,
			expectedError:       "Error: [The job 87d485b4-73fc-4a7f-bb03-720f4672947e has requested to cancel]",
			returnBody: `{
				"errors": [
					{
					"message": "The job 87d485b4-73fc-4a7f-bb03-720f4672947e has requested to cancel",
					"reason": "Requested"
					}
				]
			}`,
		},
		"correct command with error response 2": {
			executeCommand:      fmt.Sprintf("import job cancel --project-id=%s %s", projectId, jobId),
			statusCode:          http.StatusMethodNotAllowed,
			expectedCalledTimes: 1,
			expectedError:       "Error: [string]",
			returnBody: `{
				"errors": [
					{
					"message": "string",
					"reason": "string",
					"field": "string"
					}
				]
			}`,
		},
		"incorrect command with missing project id": {
			executeCommand:      fmt.Sprintf("import job cancel %s", jobId),
			statusCode:          http.StatusBadRequest,
			expectedCalledTimes: 0,
			expectedError:       "Error: required flag(s) \"project-id\" not set",
			returnBody:          ``,
		},
		"incorrect command with missing job id": {
			executeCommand:      fmt.Sprintf("import job cancel --project-id=%s", projectId),
			statusCode:          http.StatusBadRequest,
			expectedCalledTimes: 0,
			expectedError:       "Error: accepts 1 arg(s), received 0",
			returnBody:          ``,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			helper := testutils.NewAuraTestHelper(t)
			defer helper.Close()

			mockHandler := helper.NewRequestHandlerMock(fmt.Sprintf("/v2beta1/projects/%s/import/jobs/%s/cancellation", projectId, jobId), testCase.statusCode, testCase.returnBody)

			helper.SetConfigValue("aura.beta-enabled", true)

			helper.ExecuteCommand(testCase.executeCommand)

			mockHandler.AssertCalledTimes(testCase.expectedCalledTimes)
			if testCase.expectedCalledTimes > 0 {
				mockHandler.AssertCalledWithMethod(http.MethodPost)
			}

			helper.AssertErr(testCase.expectedError)
		})
	}
}
