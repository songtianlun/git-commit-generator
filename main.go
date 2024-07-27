package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Language string
	Commit   string
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var (
	apiURL    string
	modelName string
)

func init() {
	apiURL = os.Getenv("LLM_API_URL")
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1/chat/completions"
		log.Println("Using default API URL:", apiURL)
	} else {
		log.Println("Using custom API URL:", apiURL)
	}

	modelName = os.Getenv("LLM_MODEL_NAME")
	if modelName == "" {
		modelName = "gpt-3.5-turbo"
		log.Println("Using default model name:", modelName)
	} else {
		log.Println("Using custom model name:", modelName)
	}
}

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/generate", handleGenerate)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}


func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Description string `json:"description"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commit, err := generateCommit(requestBody.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"commit": commit})
}

func generateCommit(description string) (string, error) {
	prompt := fmt.Sprintf(`Given the following description of changes, generate a Git commit message following this format: <type>(<scope>): <subject>

Types include: feat, fix, docs, style, refactor, perf, test, chore, revert

The scope should be the module or area of the codebase affected.

The subject should be a concise description written in imperative mood.

Description of changes: %s

Generate the commit message, determining the appropriate type and scope based on the description.`, description)

	request := OpenAIRequest{
		Model: modelName,
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant that generates Git commit messages."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	log.Printf("Sending request to %s\n", apiURL)
	log.Printf("Request payload: %s\n", string(jsonData))

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("LLM_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Received response. Status: %s\n", resp.Status)
	log.Printf("Response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from the language model")
	}

	return strings.TrimSpace(openAIResp.Choices[0].Message.Content), nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
