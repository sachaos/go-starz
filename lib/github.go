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
}

func NewClient() *client {
	return &client{}
}

func (c *client) GetStarzList(ctx context.Context, username string) ([]*Starz, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s", username))
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
		res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", username, i))
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
