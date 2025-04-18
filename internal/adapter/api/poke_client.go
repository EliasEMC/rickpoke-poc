package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sony/gobreaker"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type PokeClient struct {
	BaseURL string
	CB      *gobreaker.CircuitBreaker
	Client  *http.Client
}

// NewPokeClient es un helper opcional
func NewPokeClient(baseURL string, cb *gobreaker.CircuitBreaker, timeout time.Duration) *PokeClient {
	return &PokeClient{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		CB:      cb,
		Client:  &http.Client{Timeout: timeout},
	}
}

// GetByName implementa service.PokemonFetcher
func (c *PokeClient) GetByName(ctx context.Context, name string) (*model.Pokemon, error) {
	out, err := c.CB.Execute(func() (interface{}, error) {
		url := fmt.Sprintf("%s/pokemon/%s", c.BaseURL, strings.ToLower(name))

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("pokeapi returned %s", resp.Status)
		}

		// Solo necesitamos el nombre y el primer tipo
		var dto struct {
			Name  string `json:"name"`
			Types []struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"types"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
			return nil, err
		}

		poke := &model.Pokemon{
			Name: strings.Title(dto.Name),
		}
		if len(dto.Types) > 0 {
			poke.Type = dto.Types[0].Type.Name
		}
		return poke, nil
	})

	if err != nil {
		return nil, err
	}
	return out.(*model.Pokemon), nil
}

var _ service.PokemonFetcher = (*PokeClient)(nil) // compile‑time check

func (c *PokeClient) State() gobreaker.State {
	return c.CB.State()
}