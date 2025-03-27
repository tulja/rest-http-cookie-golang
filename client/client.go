package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Define the URL and request body
	url := "http://localhost:8080"
	jsonData := `{"username":"user1","password":"password1"}`

	// Send a POST request with the JSON body
	resp, err := http.Post(url+"/signin", "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the status code and response body
	fmt.Printf("Response Status: %s\n", resp.Status)
	cookie := parseSessionCookieInfo(resp.Header["Set-Cookie"][0])
	fmt.Printf("Response Body: %s\n", body)

	// Make the HTTP request
	req, err := http.NewRequest("GET", url+"/welcome", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.AddCookie(cookie)
	// Send the request
	client := &http.Client{}
	resp, err = client.Do(req)

	// Read the response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", body)

	// Make the HTTP request
	req, err = http.NewRequest("GET", url+"/logout", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.AddCookie(cookie)
	// Send the request
	client = &http.Client{}
	resp, err = client.Do(req)

	// Read the response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", body)
}

func parseSessionCookieInfo(sessionStr string) *http.Cookie {
	cookie := &http.Cookie{}
	// Define the session string

	// Split the string by semicolon to separate the components
	parts := strings.Split(sessionStr, ";")
	var sessionToken, expires string

	// Iterate through each part and extract the relevant information
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "session_token") {
			// Extract session token
			parts := strings.Split(part, "=")
			sessionToken = parts[1]
		} else if strings.HasPrefix(part, "Expires") {
			// Extract expiry date
			parts := strings.Split(part, "=")
			expires = parts[1]
		}
	}

	// Parse the expiration time
	layout := "Mon, 02 Jan 2006 15:04:05 GMT"
	expiryTime, err := time.Parse(layout, expires)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return nil
	}

	// Print extracted variables
	fmt.Println("Session Token:", sessionToken)
	fmt.Println("Expires:", expires)
	fmt.Println("Parsed Expiry Time:", expiryTime)
	cookie.Name = "session_token"
	cookie.Value = sessionToken
	cookie.Expires = expiryTime
	return cookie
}
