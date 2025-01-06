package integration

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	matcher "github.com/panta/go-json-matcher"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateWorkout(t *testing.T) {
	testCases := []struct {
		name string
		req  string
		resp string
		code int
	}{
		{
			name: "create fail validations",
			req:  `{"name":"","exercises":[{"name":"","muscles":[{"name":"","muscleGroup":"Arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			resp: `{"errors": "failed to validate create workout request. workout must have a name. exercise must have a name. muscle must have a name. invalid muscle group. valid options: [arms back chest core heart legs shoulders]"}`,
			code: http.StatusBadRequest,
		},
		{
			name: "create happy path 1",
			req:  `{"name":"Arms Light","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			resp: `{"id":"#regex ^work_[a-zA-Z0-9]{27}$","name":"Arms Light","user_id":"#regex ^user_[a-zA-Z0-9]{27}$","created":"#datetime","updated":"#datetime","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			code: http.StatusCreated,
		},
		{
			name: "create happy path 2",
			req:  `{"name":"Arms Heavy","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			resp: `{"id":"#regex ^work_[a-zA-Z0-9]{27}$","name":"Arms Heavy","user_id":"#regex ^user_[a-zA-Z0-9]{27}$","created":"#datetime","updated":"#datetime","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			code: http.StatusCreated,
		},
		{
			name: "create fail conflict name",
			req:  `{"name":"Arms Heavy","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			resp: `{"errors": "workout name already exists"}`,
			code: http.StatusConflict,
		},
	}

	suffix := "_test_create_workout"
	email := randomEmail()
	username := randomUsername(suffix)
	userID, err := createTestUser(username, email)
	if err != nil {
		t.Log("err", err)
		t.Fail()
	}

	testCookie, err := loginTestUser(username)
	if err != nil {
		t.Log("err", err)
		t.Fail()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/users/%s/workouts", userID)
			resp, err := sendJSONHttpRequest("POST", url, tc.req, testCookie)
			if err != nil {
				t.Log("err", err)
				t.Fail()
			}

			assert.Equal(t, resp.StatusCode, tc.code)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			respString := string(bodyBytes)

			matches, err := matcher.JSONStringMatches(respString, tc.resp)
			if !matches || err != nil {
				t.Log("err", err)
				t.Log("got", respString)
				t.Log("wanted", tc.resp)
				t.Fail()
			}
		})
	}
	cleanUpTestUsers(suffix)
}
