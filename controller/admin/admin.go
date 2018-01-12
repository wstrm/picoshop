package admin

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/willeponken/picoshop/controller/helper"
	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

var (
	// Allowed characters for random string generator @see randomString
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type rootHandler struct {
	http.Handler
}

type registerHandler struct {
	http.Handler
}

type articleHandler struct {
	http.Handler
}

type rootData struct {
	Error string
}

type registerData helper.RegisterResult

type articleData struct {
	CreateError   string
	StockError    string
	CreateSuccess bool
	StockSuccess  bool
	Name          string
	Description   string
	Price         string
	Id            int64
	Category      string
	Subcategory   string
	InStock       string
}

func (a *rootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "admin", view.Page{Title: "Admin - Picoshop", Data: rootData{}})
	}
}

func (a *registerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	userType := request.URL.Query().Get("type")

	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "admin.register", view.Page{Title: "Admin - Picoshop", Data: registerData{
			Type: userType,
		}})
	case http.MethodPost:
		if err := request.ParseForm(); err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		form := helper.ParseRegisterFormValues(request)
		_, result, code := helper.Register(userType, form)

		writer.WriteHeader(code)
		view.Render(request.Context(), writer, "admin.register", view.Page{
			Title: "Admin - Picoshop",
			Data:  result,
		})
		return
	}
}

func randomString(length int) string {
	bits := make([]rune, length)

	for char := range bits {
		bits[char] = characters[rand.Intn(len(characters))]
	}

	return string(bits)
}

func (a *articleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	switch request.FormValue("_method") {
	case http.MethodPut:
		method = http.MethodPut
	case http.MethodDelete:
		method = http.MethodDelete
	}

	switch method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "admin.article", view.Page{Title: "Admin - Picoshop", Data: articleData{}})
	case http.MethodPost:
		request.ParseMultipartForm(32 << 20)

		in, handler, err := request.FormFile("image")
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer in.Close()

		imageName := fmt.Sprintf("%s-%s", randomString(10), handler.Filename)
		out, err := os.OpenFile("./static/image/"+imageName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		name := request.PostFormValue("name")
		description := request.PostFormValue("description")
		price := request.PostFormValue("price")
		category := request.PostFormValue("category")
		subcategory := request.PostFormValue("subcategory")
		inStock := request.PostFormValue("in-stock")

		if err := helper.IsFilled(name, description, price, category, subcategory, inStock); err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		p, err := strconv.ParseFloat(price, 10)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s, err := strconv.ParseUint(inStock, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		article, err := model.PutArticle(model.NewArticle(name, description, p, imageName, category, subcategory, s))
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		view.Render(request.Context(), writer, "admin.article", view.Page{Title: "Admin - Picoshop",
			Data: articleData{
				Name:          name,
				Description:   description,
				Price:         price,
				Id:            article.Id,
				CreateSuccess: true,
			}})
	case http.MethodPut:
		data := articleData{
			InStock: request.FormValue("in-stock"),
		}

		inStock, err := strconv.ParseUint(data.InStock, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Invalid in stock data", http.StatusBadRequest)
			return
		}

		articleId, err := strconv.ParseInt(request.FormValue("id"), 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Invalid article id", http.StatusBadRequest)
			return
		}
		data.Id = articleId

		err = model.SetStockById(articleId, inStock)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data.StockSuccess = true
		view.Render(request.Context(), writer, "admin.article", view.Page{Title: "Admin - Picoshop",
			Data: data,
		})

	}
}

func newRootHandler() *rootHandler {
	return &rootHandler{}
}

func newArticleHandler() *articleHandler {
	return &articleHandler{}
}

func newRegisterHandler() *registerHandler {
	return &registerHandler{}
}

func NewMux() *http.ServeMux {
	rand.Seed(time.Now().UTC().UnixNano())

	mux := http.NewServeMux()

	mux.Handle("/", newRootHandler())
	mux.Handle("/register", newRegisterHandler())
	mux.Handle("/article", newArticleHandler())

	return mux
}
