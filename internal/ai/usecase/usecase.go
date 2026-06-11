package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Input struct {
	Description string
	Categories  []Category
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Result struct {
	SkillID   *int   `json:"skill_id"`
	Reasoning string `json:"reasoning"`
}

type FindWorkersInput struct {
	Description string
	Categories  []Category
}

type FindWorkersResult struct {
	SkillIDs  []int    `json:"skill_ids"`
	Reasoning string   `json:"reasoning"`
}

type AIService struct {
	apiURL   string
	apiKey   string
	model    string
	client   *http.Client
}

func NewAIService(apiURL, apiKey, model string) *AIService {
	return &AIService{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *AIService) DescribeSkill(ctx context.Context, input Input) (*Result, error) {
	var sb strings.Builder
	sb.WriteString("Anda adalah asisten cerdas platform KerjaDekat. Tugas Anda adalah memetakan deskripsi kebutuhan pengguna ke salah satu kategori jasa yang tersedia.\n\nKategori jasa yang tersedia:\n")
	for _, c := range input.Categories {
		sb.WriteString(fmt.Sprintf("- ID %d: %s\n", c.ID, c.Name))
	}
	sb.WriteString(fmt.Sprintf("\nDeskripsi pengguna: \"%s\"\n", input.Description))
	sb.WriteString("\nBerikan respons dalam format JSON murni (langsung JSON, tanpa markdown):\n")
	sb.WriteString(`{"skill_id": <nomor ID kategori yang paling relevan, atau null jika tidak ada>, "reasoning": "<penjelasan singkat dalam Bahasa Indonesia mengapa kategori ini dipilih>"}`)

	payload := map[string]any{
		"model": s.model,
		"messages": []map[string]string{
			{"role": "user", "content": sb.String()},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.1,
		"max_tokens":      512,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ai api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai api error %d: %s", resp.StatusCode, string(respBody))
	}

	raw := string(respBody)

	// 9router may return SSE streaming even without stream=true.
	// Strip trailing SSE data: [DONE] or any data: lines after the JSON.
	if idx := strings.Index(raw, "data: "); idx != -1 {
		raw = raw[:idx]
	}

	// Trim whitespace / BOM / leftover
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return nil, fmt.Errorf("ai returned empty body")
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(raw), &chatResp); err != nil {
		log.Printf("[AI] raw response: %s", raw)
		return nil, fmt.Errorf("parse chat response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("ai returned no choices")
	}

	content := chatResp.Choices[0].Message.Content
	var result Result
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("parse ai content json: %w (content: %s)", err, content)
	}

	return &result, nil
}

func (s *AIService) FindWorkers(ctx context.Context, input FindWorkersInput) (*FindWorkersResult, error) {
	var sb strings.Builder
	sb.WriteString("Anda adalah asisten cerdas platform KerjaDekat. Tugas Anda adalah memetakan deskripsi kebutuhan pengguna ke satu atau lebih kategori jasa yang tersedia.\n\nKategori jasa yang tersedia:\n")
	for _, c := range input.Categories {
		sb.WriteString(fmt.Sprintf("- ID %d: %s\n", c.ID, c.Name))
	}
	sb.WriteString(fmt.Sprintf("\nDeskripsi pengguna: \"%s\"\n", input.Description))
	sb.WriteString("\nBerikan respons dalam format JSON murni (langsung JSON, tanpa markdown):\n")
	sb.WriteString(`{"skill_ids": [<array nomor ID kategori yang relevan, kosongkan jika tidak ada>], "reasoning": "<penjelasan singkat dalam Bahasa Indonesia mengapa kategori ini dipilih>"}`)

	payload := map[string]any{
		"model": s.model,
		"messages": []map[string]string{
			{"role": "user", "content": sb.String()},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.1,
		"max_tokens":      512,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ai api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai api error %d: %s", resp.StatusCode, string(respBody))
	}

	raw := string(respBody)

	if idx := strings.Index(raw, "data: "); idx != -1 {
		raw = raw[:idx]
	}

	raw = strings.TrimSpace(raw)

	if raw == "" {
		return nil, fmt.Errorf("ai returned empty body")
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(raw), &chatResp); err != nil {
		log.Printf("[AI] raw response: %s", raw)
		return nil, fmt.Errorf("parse chat response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("ai returned no choices")
	}

	content := chatResp.Choices[0].Message.Content
	var result FindWorkersResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("parse ai content json: %w (content: %s)", err, content)
	}

	return &result, nil
}
