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
	defaultPerPage  = 25
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
		for i := (page-1)*perPage + 1; i < page*perPage; i++ {
			requestRecipeIDs = append(requestRecipeIDs, strconv.Itoa(i))
		}

		recipes, err = recipeRepo.FindByIDs(r.Context(), requestRecipeIDs...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextURL := url.URL{}
		nextURL.Host = r.URL.Host
		nextURL.Path = r.URL.Path
		nextURL.Query().Set(pageQueryKey, strconv.Itoa(page+1))
		nextURL.Query().Set(perPageQueryKey, strconv.Itoa(perPage))
		prevURL := url.URL{}
		prevURL.Host = r.URL.Host
		prevURL.Path = r.URL.Path
		prevURL.Query().Set(pageQueryKey, strconv.Itoa(page-1))
		prevURL.Query().Set(perPageQueryKey, strconv.Itoa(perPage))

		err = json.NewEncoder(w).Encode(
			RecipeResponse{
				Data: recipes,
				Pagination: Pagination{
					PrevLink: prevURL.String(),
					NextLink: nextURL.String(),
				},
			},
		)
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
