package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/slham/sandbox-api/dao"
	"github.com/tamathecxder/randomail"
)

var testDB *sql.DB

func createTestUser(username, email string) (string, error) {
	url := "http://localhost:8080/users"
	method := "POST"
	body := fmt.Sprintf(`{"username": "%s", "password": "thisIsAG00dPassword!", "email": "%s"}`, username, email)
	resp, err := sendJSONHttpRequest(method, url, body, nil)
	if err != nil {
		return "", fmt.Errorf("failed to execute request. %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create test user. %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body. %w", err)
	}

	user := map[string]any{}
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body. %w", err)
	}

	userID, ok := user["id"].(string)
	if !ok {
		return "", fmt.Errorf("got invalid response. user.id:%v", userID)
	}

	return userID, nil
}

func loginTestUser(username string) (*http.Cookie, error) {
	url := "http://localhost:8080/auth/login"
	method := "POST"
	body := fmt.Sprintf(`{"username": "%s", "password": "thisIsAG00dPassword!"}`, username)
	resp, err := sendJSONHttpRequest(method, url, body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request. %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to create test user. %w", err)
	}

	fmt.Println("resp headers", resp.Header)

	var testCookie *http.Cookie
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "sandbox-cookie" {
			testCookie = cookie
			break
		}
	}

	return testCookie, nil
}

func getDB() (*sql.DB, error) {
	if testDB != nil {
		return testDB, nil
	}

	database, err := dao.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database. %w", err)
	}

	return database.DB, nil
}

func cleanUpTestUsers(suffix string) error {
	db, err := getDB()
	if err != nil {
		return fmt.Errorf("failed to get test db. %w", err)
	}

	stmt := "DELETE FROM sandbox.user WHERE username LIKE '%" + suffix + "'"

	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to clean up %s text users. %w", suffix, err)
	}

	return nil
}

func sendJSONHttpRequest(method, url, body string, testCookie *http.Cookie) (*http.Response, error) {
	bodyReader := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to build request. %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if testCookie != nil {
		req.AddCookie(testCookie)
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request. %w", err)
	}

	return resp, nil
}

func randomEmail() string {
	return randomail.GenerateRandomEmail()
}

func randomUsername(suffix string) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
