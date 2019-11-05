package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"vutung2311-golang-test/internal/model"
)

type httpClient interface {
	Get(ctx context.Context, url string) (resp *http.Response, err error)
}

type workerPool interface {
	AddJob(ctx context.Context, fn func() error) error
}

type RecipeRepository struct {
	baseUrl          string
	recipeJsonGetter httpClient
	pool             workerPool
}

func NewRecipeRepository(baseUrl string, recipeJsonGetter httpClient, pool workerPool) *RecipeRepository {
	return &RecipeRepository{
		baseUrl:          baseUrl,
		recipeJsonGetter: recipeJsonGetter,
		pool:             pool,
	}
}

type RecipeFetch struct {
	recipe model.Recipe
	err    error
}

func (r *RecipeRepository) FindByIDs(ctx context.Context, ids ...string) ([]model.Recipe, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(len(ids))
	fetchingChan := make(chan RecipeFetch, len(ids))
	for i := range ids {
		id := ids[i]
		err := r.pool.AddJob(cancelCtx, func() error {
			defer wg.Done()

			resp, err := r.recipeJsonGetter.Get(cancelCtx, r.baseUrl+id)
			if err != nil {
				cancelFunc()
				fetchingChan <- RecipeFetch{err: err}
				return nil
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
				fetchingChan <- RecipeFetch{err: errors.New("bad http status code")}
				return nil
			}
			if resp.StatusCode == http.StatusForbidden {
				return fmt.Errorf("got %d for URL %s", resp.StatusCode, r.baseUrl+id)
			}
			recipeFetch := RecipeFetch{recipe: model.Recipe{}, err: nil}
			recipeFetch.err = json.NewDecoder(resp.Body).Decode(&recipeFetch.recipe)
			fetchingChan <- recipeFetch
			return nil
		})
		if err != nil {
			cancelFunc()
			return nil, err
		}
	}
	wg.Wait()
	close(fetchingChan)

	result := make([]model.Recipe, 0)
	for recipeFetch := range fetchingChan {
		if recipeFetch.err != nil {
			return nil, recipeFetch.err
		}
		result = append(result, recipeFetch.recipe)
	}

	return result, nil
}
