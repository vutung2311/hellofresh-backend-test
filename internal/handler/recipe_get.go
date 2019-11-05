package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"vutung2311-golang-test/internal/model"
	"vutung2311-golang-test/internal/repository"
)

const (
	defaultPerPage  = 5
	pageQueryKey    = "page"
	perPageQueryKey = "perPage"
)

func CreateRecipeGetHandler(recipeRepo *repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			page    int
			perPage = defaultPerPage
			recipes []model.Recipe
		)

		pageParam := r.URL.Query().Get("page")
		perPageParam := r.URL.Query().Get("perPage")
		if len(pageParam) > 0 {
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				http.Error(w, fmt.Sprintf("bad page param: %s", pageParam), http.StatusUnprocessableEntity)
				return
			}
		}
		if len(perPageParam) > 0 {
			perPage, err = strconv.Atoi(perPageParam)
			if err != nil {
				http.Error(w, fmt.Sprintf("bad per page param: %s", perPageParam), http.StatusUnprocessableEntity)
				return
			}
		}

		if page < 1 {
			page = 1
		}
		if perPage < 0 {
			perPage = defaultPerPage
		}

		requestRecipeIDs := make([]string, 0)
		for i := (page-1)*perPage + 1; i <= page*perPage; i++ {
			requestRecipeIDs = append(requestRecipeIDs, strconv.Itoa(i))
		}

		recipes, err = recipeRepo.FindByIDs(r.Context(), requestRecipeIDs...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextURL := new(url.URL)
		if r.TLS == nil {
			nextURL.Scheme = "http"
		} else {
			nextURL.Scheme = "https"
		}
		nextURL.Host = r.Host
		nextURL.Path = r.URL.Path
		nextURLQuery := nextURL.Query()
		nextURLQuery.Set(pageQueryKey, strconv.Itoa(page+1))
		nextURLQuery.Set(perPageQueryKey, strconv.Itoa(perPage))
		nextURL.RawQuery = nextURLQuery.Encode()

		prevURL := new(url.URL)
		if r.TLS == nil {
			prevURL.Scheme = "http"
		} else {
			prevURL.Scheme = "https"
		}
		prevURL.Host = r.Host
		prevURL.Path = r.URL.Path
		prevURLQuery := prevURL.Query()
		prevURLQuery.Set(pageQueryKey, strconv.Itoa(page-1))
		prevURLQuery.Set(perPageQueryKey, strconv.Itoa(perPage))
		prevURL.RawQuery = prevURLQuery.Encode()

		recipeResponse := RecipeResponse{
			Data: recipes,
			Pagination: Pagination{
				NextLink: nextURL.String(),
			},
		}
		if page > 1 {
			recipeResponse.Pagination.PrevLink = prevURL.String()
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(recipeResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

type RecipeResponse struct {
	Data       interface{} `json:"data"`
	Pagination `json:"pagination"`
}

type Pagination struct {
	PrevLink string `json:"prevLink,omitempty"`
	NextLink string `json:"nextLink,omitempty"`
}
