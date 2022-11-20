package router

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	"github.com/iTukaev/news_service/internal/client"
	errorsPkg "github.com/iTukaev/news_service/internal/customerrors"
)

type param string

func (p param) str() string {
	return string(p)
}

const (
	searchParam param = "search"
	sizeParam   param = "size"
	pageParam   param = "page"
	orderParam  param = "order"
)

func (c *Core) getNews() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		search := values.Get(searchParam.str())
		if search == "" {
			c.fail(w, http.StatusBadRequest, "search parameter expected")
			return
		}

		news, err := c.app.NewsGet(r.Context(), search)
		if err != nil {
			if errors.Is(err, errorsPkg.ErrNewsNotFound) {
				c.fail(w, http.StatusBadRequest, err.Error())
				return
			}
			c.fail(w, http.StatusInternalServerError, err.Error())
			return
		}

		body, err := json.Marshal(news)
		if err != nil {
			c.fail(w, http.StatusInternalServerError, "internal error")
			return
		}

		c.success(w, http.StatusOK, body)
	}
}

func (c *Core) getNewsList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		params := getListParameters(values)

		offset := (params.page - 1) * params.size
		news, err := c.app.NewsList(r.Context(), params.size, offset, params.order)
		if err != nil {
			if errors.Is(err, errorsPkg.ErrNewsNotFound) {
				c.fail(w, http.StatusBadRequest, err.Error())
				return
			}
			c.fail(w, http.StatusInternalServerError, err.Error())
			return
		}

		list := client.NewsList{
			List: news,
			Page: client.Page{
				Size: uint64(len(news)),
				Page: params.page,
			},
		}
		body, err := json.Marshal(&list)
		if err != nil {
			c.fail(w, http.StatusInternalServerError, "internal error")
			return
		}
		c.success(w, http.StatusOK, body)

	}
}

func (c *Core) fail(w http.ResponseWriter, status int, description string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(description))
	if err != nil {
		c.logger.Errorln("response write", err)
	}
}

func (c *Core) success(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		c.logger.Errorln("response write", err)
	}
}

func getListParameters(values url.Values) listParams {
	sizeStr := values.Get(sizeParam.str())
	size, _ := strconv.ParseUint(sizeStr, 10, 64)

	pageStr := values.Get(pageParam.str())
	page, _ := strconv.ParseUint(pageStr, 10, 64)

	orderStr := values.Get(orderParam.str())
	order, _ := strconv.ParseBool(orderStr)

	return listParams{
		size:  size,
		page:  page,
		order: order,
	}
}
