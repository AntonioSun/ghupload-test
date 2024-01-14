package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: main.go <token> <repo> <issue number> <image paths>")
		os.Exit(1)
	}

	token := os.Args[1]
	repo := os.Args[2]
	issueNumber := os.Args[3]
	imagePaths := os.Args[4:]

	var uploadedImageURLs []string
	for _, imagePath := range imagePaths {
		imageURL, err := uploadImageAndGetURL(imagePath)
		if err != nil {
			fmt.Printf("Error uploading image '%s': %v\n", imagePath, err)
			os.Exit(1)
		}
		uploadedImageURLs = append(uploadedImageURLs, imageURL)
	}

	commentBody := createCommentBody(uploadedImageURLs)
	err := postComment(token, repo, issueNumber, commentBody)
	if err != nil {
		fmt.Printf("Error posting comment: %v\n", err)
		os.Exit(1)
	}

	outputUploadedURLs(uploadedImageURLs)
}

func uploadImageAndGetURL(imagePath string) (string, error) {
	// Implement your image upload logic here
	// Return the URL after successful upload
	return "https://example.com/" + imagePath, nil
}

func createCommentBody(urls []string) string {
	var sb strings.Builder
	for _, url := range urls {
		sb.WriteString(fmt.Sprintf("![Image](%s)\n", url))
	}
	return sb.String()
}

func postComment(token, repo, issueNumber, comment string) error {
	data, err := json.Marshal(map[string]string{"body": comment})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", repo, issueNumber)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("GitHub API responded with status: %d", resp.StatusCode)
	}
	return nil
}

func outputUploadedURLs(urls []string) {
	for i, url := range urls {
		fmt.Printf("::set-output name=image_url_%d::%s\n", i+1, url)
	}
}
