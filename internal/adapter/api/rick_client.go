package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type RickClient struct {
	BaseURL string
	CB      *gobreaker.CircuitBreaker
	Client  *http.Client
}

func NewRickClient(baseURL string, cb *gobreaker.CircuitBreaker, timeout time.Duration) *RickClient {
	return &RickClient{
		BaseURL: baseURL,
		CB:      cb,
		Client:  &http.Client{Timeout: timeout},
	}
}

func (c *RickClient) GetByID(ctx context.Context, id int) (*model.Rick, error) {
	out, err := c.CB.Execute(func() (interface{}, error) {
		url := fmt.Sprintf("%s/character/%d", c.BaseURL, id)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("rickapi returned %s", resp.Status)
		}

		var dto struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Origin struct {
				Name string `json:"name"`
			} `json:"origin"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
			return nil, err
		}
		return &model.Rick{
			ID:     dto.ID,
			Name:   dto.Name,
			//Origin: dto.Origin.Name,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*model.Rick), nil
}

var _ service.CharacterFetcher = (*RickClient)(nil)

func (c *RickClient) State() gobreaker.State {
	return c.CB.State()
}
