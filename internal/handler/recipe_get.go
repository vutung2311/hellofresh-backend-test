package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"vutung2311-golang-test/pkg/pagination"

	"vutung2311-golang-test/internal/model"
	"vutung2311-golang-test/internal/repository"
)

type RecipeResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	PrevLink string `json:"prevLink,omitempty"`
	NextLink string `json:"nextLink,omitempty"`
}

func CreateRecipeGetHandler(recipeRepo *repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err                error
			recipes            []*model.Recipe
			prevLink, nextLink string
		)

		requestedRecipes, err := getRequestedRecipes(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		recipes, err = recipeRepo.FindByIDs(r.Context(), requestedRecipes...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response RecipeResponse
		if !isGettingRecipeByID(r) {
			prevLink, nextLink, _ = getPrevNextLink(r)
			response.Pagination = &Pagination{
				NextLink: nextLink,
				PrevLink: prevLink,
			}
		} else {
			sort.Slice(recipes, func(i, j int) bool {
				return strings.Compare(recipes[i].PrepTime, recipes[j].PrepTime) < 0
			})
		}
		response.Data = recipes

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

func isGettingRecipeByID(r *http.Request) bool {
	return len(r.URL.Query().Get("id")) > 0
}

func getPrevNextLink(r *http.Request) (string, string, error) {
	if isGettingRecipeByID(r) {
		return "", "", nil
	}
	page, perPage, err := pagination.GetPageAndPerPage(r)
	if err != nil {
		return "", "", err
	}
	prevLink, nextLink := pagination.GetPrevAndNextLink(r, page, perPage)
	return prevLink, nextLink, nil
}

func getRequestedRecipes(r *http.Request) ([]string, error) {
	recipeIdList := r.URL.Query().Get("id")
	if len(recipeIdList) > 0 {
		return strings.Split(recipeIdList, ","), nil
	}
	page, perPage, err := pagination.GetPageAndPerPage(r)
	if err != nil {
		return nil, err
	}
	requestedRecipes := make([]string, 0)
	for i := (page-1)*perPage + 1; i <= page*perPage; i++ {
		requestedRecipes = append(requestedRecipes, strconv.Itoa(i))
	}
	return requestedRecipes, nil
}
