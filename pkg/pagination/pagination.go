package pagination

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultPerPage  = 5
	pageQueryKey    = "page"
	perPageQueryKey = "perPage"
)

func GetPrevAndNextLink(r *http.Request, currentPage, currentPerPage int) (prevLink, nextLink string) {
	nextURL := new(url.URL)
	if r.TLS == nil {
		nextURL.Scheme = "http"
	} else {
		nextURL.Scheme = "https"
	}
	nextURL.Host = r.Host
	nextURL.Path = r.URL.Path
	nextURLQuery := nextURL.Query()
	nextURLQuery.Set(pageQueryKey, strconv.Itoa(currentPage+1))
	nextURLQuery.Set(perPageQueryKey, strconv.Itoa(currentPerPage))
	nextURL.RawQuery = nextURLQuery.Encode()
	nextLink = nextURL.String()

	if currentPage > 1 {
		prevURL := new(url.URL)
		if r.TLS == nil {
			prevURL.Scheme = "http"
		} else {
			prevURL.Scheme = "https"
		}
		prevURL.Host = r.Host
		prevURL.Path = r.URL.Path
		prevURLQuery := prevURL.Query()
		prevURLQuery.Set(pageQueryKey, strconv.Itoa(currentPage-1))
		prevURLQuery.Set(perPageQueryKey, strconv.Itoa(currentPerPage))
		prevURL.RawQuery = prevURLQuery.Encode()
		prevLink = prevURL.String()
	}

	return
}

func GetPageAndPerPage(r *http.Request) (page, perPage int, err error) {
	pageParam := r.URL.Query().Get("page")
	perPageParam := r.URL.Query().Get("perPage")
	if len(pageParam) > 0 {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			return 0, 0, fmt.Errorf("bad page param: %v", err)
		}
	}
	if len(perPageParam) > 0 {
		perPage, err = strconv.Atoi(perPageParam)
		if err != nil {
			return 0, 0, fmt.Errorf("bad per page param: %v", err)
		}
	}

	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = defaultPerPage
	}

	return
}
