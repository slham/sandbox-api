package main

/*
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestWorkout(t *testing.T) {
	testCases := []struct {
		name string
		req  string
		resp map[string]string
		code int
	}{
		{
			name:   "create fail validations",
			method: "POST",
			url:    "/users/%s/workouts",
			req:    `{"name":"","exercises":[{"name":"","muscles":[{"name":"","muscleGroup":"Arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
			resp:   map[string]string{"errors": "failed to validate create workout request. workout must have a name. exercise must have a name. muscle must have a name"},
			code:   http.StatusBadRequest,
		},
			{
				name:   "create happy path 1",
				method: "POST",
				url:    "/users/%s/workouts",
				req:    `{"name":"workout1","exercises":[{"name":"curl","muscles":[{"name":"bicep","muscleGroup":"arms"}],"sets":[{"weight":25,"reps":10},{"weight":25,"reps":10},{"weight":25,"reps":10}]}]}`,
				resp:   map[string]string{"username": "test_user", "email": "a@b.c"},
				code:   http.StatusCreated,
			},
	}

	for _, tc := range testCases {
		bodyReader := bytes.NewReader([]byte(tc.req))
		url := "http://localhost:8080/workouts"
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

		for k, v := range tc.resp {
			assert.Equal(t, respMap[k], v)
		}

	}
}
*/
