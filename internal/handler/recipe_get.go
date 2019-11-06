package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vutung2311-golang-test/pkg/pagination"

	"vutung2311-golang-test/internal/model"
	"vutung2311-golang-test/internal/repository"
)

func CreateRecipeGetHandler(recipeRepo *repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err                error
			page, perPage      int
			recipes            []*model.Recipe
			prevLink, nextLink string
		)

		page, perPage, err = pagination.GetPageAndPerPage(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		prevLink, nextLink = pagination.GetPrevAndNextLink(r, page, perPage)

		requestRecipeIDs := make([]string, 0)
		for i := (page-1)*perPage + 1; i <= page*perPage; i++ {
			requestRecipeIDs = append(requestRecipeIDs, strconv.Itoa(i))
		}

		recipes, err = recipeRepo.FindByIDs(r.Context(), requestRecipeIDs...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		recipeResponse := RecipeResponse{
			Data: recipes,
			Pagination: Pagination{
				NextLink: nextLink,
				PrevLink: prevLink,
			},
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
