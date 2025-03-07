package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type User struct {
	PublicRepos int `json:"public_repos"`
}

type Starz struct {
	Name            string `json:"name"`
	StargazersCount int    `json:"stargazers_count"`
}

type GitHub interface {
	GetStarzList(ctx context.Context, username string) []*Starz
}

type client struct {
	token string
}

func NewClient(token string) *client {
	return &client{token: token}
}

func (c *client) GetStarzList(ctx context.Context, username string) ([]*Starz, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s", username), nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get starz list: %s", res.Status)
	}

	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	pageCount := int(math.Ceil(float64(user.PublicRepos) / 100))

	var allStarz []*Starz
	for i := 1; i <= pageCount; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", username, i), nil)
		if err != nil {
			return nil, err
		}
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get starz list: %s", res.Status)
		}

		var starz []*Starz
		err = json.NewDecoder(res.Body).Decode(&starz)
		if err != nil {
			return nil, err
		}

		allStarz = append(allStarz, starz...)
	}

	return allStarz, nil
}
