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

var ErrRecordNotFound = errors.New("record not fund")

type httpClient interface {
	Get(ctx context.Context, url string) (resp *http.Response, err error)
}

type workerPool interface {
	AddJob(ctx context.Context, fn func() error) error
	IsRetryableError(err error) bool
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
	recipe *model.Recipe
	err    error
}

func (r *RecipeRepository) FindByID(ctx context.Context, id string) (*model.Recipe, error) {
	resp, err := r.recipeJsonGetter.Get(ctx, r.baseUrl+id)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		return nil, errors.New("bad http status code")
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrRecordNotFound
	}

	var recipe model.Recipe
	err = json.NewDecoder(resp.Body).Decode(&recipe)
	_ = resp.Body.Close()
	return &recipe, err
}

func (r *RecipeRepository) FindByIDs(ctx context.Context, ids ...string) ([]*model.Recipe, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(len(ids))
	fetchingChan := make(chan RecipeFetch, len(ids))
	for i := range ids {
		id := ids[i]
		var err error
		for {
			err = r.pool.AddJob(cancelCtx, func() error {
				defer wg.Done()

				recipe, err := r.FindByID(cancelCtx, id)
				if err != nil && !errors.Is(err, ErrRecordNotFound) {
					cancelFunc()
					fetchingChan <- RecipeFetch{recipe: nil, err: err}
					return nil
				}
				if errors.Is(err, ErrRecordNotFound) {
					return fmt.Errorf("record in URL %s is not found", r.baseUrl+id)
				}
				fetchingChan <- RecipeFetch{recipe: recipe, err: nil}
				return nil
			})
			if !r.pool.IsRetryableError(err) {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}
	wg.Wait()
	close(fetchingChan)

	result := make([]*model.Recipe, 0)
	for recipeFetch := range fetchingChan {
		if recipeFetch.err != nil {
			return nil, recipeFetch.err
		}
		result = append(result, recipeFetch.recipe)
	}

	return result, nil
}
