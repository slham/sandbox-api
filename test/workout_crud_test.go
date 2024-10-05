package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	matcher "github.com/panta/go-json-matcher"
	"gopkg.in/go-playground/assert.v1"
)

func TestWorkout(t *testing.T) {
	testCases := []struct {
		name    string
		req     string
		method  string
		url     string
		resp    string
		code    int
		comment string
	}{
		{
			name:   "create fail validations",
			method: "POST",
			url:    "/workouts",
			req:    `{"name":"","exercises":[{"name":"","muscles":[{"name":"","muscleGroup":"Arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			resp:   `{"errors": "failed to validate create workout request. workout must have a name. exercise must have a name. muscle must have a name. invalid muscle group. valid options: [arms back chest core heart legs shoulders]"}`,
			code:   http.StatusBadRequest,
		},
		{
			name:   "create happy path 1",
			method: "POST",
			url:    "/workouts",
			req:    `{"name":"Arms Light","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			resp:   `{"id":"#regex ^work_[a-zA-Z0-9]{27}$","name":"Arms Light","user_id":"#regex ^user_[a-zA-Z0-9]{27}$","created":"#datetime","updated":"#datetime","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create happy path 2",
			method: "POST",
			url:    "/workouts",
			req:    `{"name":"Arms Heavy","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			resp:   `{"id":"#regex ^work_[a-zA-Z0-9]{27}$","name":"Arms Heavy","user_id":"#regex ^user_[a-zA-Z0-9]{27}$","created":"#datetime","updated":"#datetime","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create fail conflict name",
			method: "POST",
			url:    "/workouts",
			req:    `{"name":"Arms Heavy","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			resp:   `{"errors": "workout name already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "update fail validations",
			method: "PATCH",
			url:    "/workouts/%s",
			req:    `{"name":"","exercises":[{"name":"","muscles":[{"name":"","muscleGroup":"Arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			resp:   `{"errors":"failed to validate update workout request. workout must have a name. exercise must have a name. muscle must have a name. invalid muscle group. valid options: [arms back chest core heart legs shoulders]"}`,
			code:   http.StatusBadRequest,
		},
		{
			name:   "update fail conflict name",
			method: "PATCH",
			url:    "/workouts/%s",
			req:    `{"name":"Arms Heavy","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			resp:   `{"errors": "workout name already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "update happy path",
			method: "PATCH",
			url:    "/workouts/%s",
			req:    `{"name":"Popeye","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			resp:   `{"id":"#regex ^work_[a-zA-Z0-9]{27}$","name":"Popeye","user_id":"#regex ^user_[a-zA-Z0-9]{27}$","created":"#datetime","updated":"#datetime","exercises":[{"name":"Curl","muscles":[{"name":"Bicep","muscleGroup":"arms"}],"sets":[{"weight":45,"reps":10},{"weight":45,"reps":10},{"weight":45,"reps":10}]}]}`,
			code:   http.StatusOK,
		},
		{
			name:   "delete happy path",
			method: "DELETE",
			url:    "/workouts/%s",
			code:   http.StatusNoContent,
		},
		{
			name:    "get all happy path",
			method:  "GET",
			url:     "/workouts",
			code:    http.StatusOK,
			comment: "clean up test data",
		},
	}

	userID := createTestUser(t)
	workoutIDs := []string{}
	for _, tc := range testCases {
		bodyReader := bytes.NewReader([]byte(tc.req))
		url := fmt.Sprintf("http://localhost:8080/users/%s", userID)
		if tc.method == "PATCH" || tc.method == "DELETE" {
			tc.url = fmt.Sprintf(tc.url, workoutIDs[0])
		}
		url = fmt.Sprintf("%s%s", url, tc.url)
		req, err := http.NewRequest(tc.method, url, bodyReader)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		assert.Equal(t, resp.StatusCode, tc.code)

		if resp.StatusCode == http.StatusNoContent {
			t.Logf("skipping test success")
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		respString := string(bodyBytes)

		if tc.method == "GET" {
			workoutsList := []map[string]any{}
			err = json.Unmarshal(bodyBytes, &workoutsList)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			for _, workout := range workoutsList {
				u := fmt.Sprintf("http://localhost:8080/users/%s/workouts/%s", userID, workout["id"])
				rq, err := http.NewRequest("DELETE", u, nil)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				rq.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(rq)
				if err != nil {
					t.Log(err)
					t.Fail()
				}

				assert.Equal(t, resp.StatusCode, http.StatusNoContent)
			}

			t.Logf("skipping test success")
			continue
		}

		matches, err := matcher.JSONStringMatches(respString, tc.resp)
		if !matches || err != nil {
			t.Log(err)
			t.Logf("resp: %s. expected: %s", respString, tc.resp)
			t.Fail()
		}

		if tc.method == "POST" && resp.StatusCode == http.StatusCreated {
			respMap := map[string]any{}
			err = json.Unmarshal(bodyBytes, &respMap)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			if id, ok := respMap["id"].(string); ok {
				workoutIDs = append(workoutIDs, id)
			}
		}
	}

}

func createTestUser(t *testing.T) string {
	url := "http://localhost:8080/users"
	reqBody := `{"username": "test_user_workout_crud", "password": "thisIsAG00dPassword!", "email": "test@workoutCrud.com"}`
	bodyReader := bytes.NewReader([]byte(reqBody))
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	respMap := map[string]string{}
	err = json.Unmarshal(bodyBytes, &respMap)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	return respMap["id"]
}

func deleteTestUser(t *testing.T, userID string) {
	url := fmt.Sprintf("http://localhost:8080/users/%s", userID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Log("failed to clean up test user")
		t.Fail()
	}
}
