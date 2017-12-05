package register

import (
	"context"
	"log"
	"net/http"

	"github.com/willeponken/picoshop/controller/helper"
	"github.com/willeponken/picoshop/middleware/auth"
	"github.com/willeponken/picoshop/view"
)

type registerHandler struct {
	http.Handler
	authManager *auth.Manager
}

type registerData helper.RegisterResult

func renderRegister(ctx context.Context, writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = registerData{}
	}

	page := view.Page{
		Title: "Register - Picoshop",
		Data:  data,
	}

	view.Render(ctx, writer, "register", page)
}

func (r *registerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userType := request.URL.Query().Get("type")

	switch request.Method {
	case http.MethodGet: // Serve register view
		renderRegister(ctx, writer, http.StatusOK, registerData{
			Type: userType,
		})

	case http.MethodPost: // Retreive user registration
		err := request.ParseForm()
		if err != nil {
			renderRegister(ctx, writer, http.StatusBadRequest, registerData{
				Error: helper.InvalidFormDataError().Error(),
				Type:  userType,
			})
			return
		}

		form := helper.ParseRegisterFormValues(request)
		log.Println(form)
		user, result, code := helper.Register(userType, form)
		log.Println(user, result, code)

		if code != http.StatusOK {
			renderRegister(ctx, writer, code, result)
			return
		}

		err = r.authManager.Login(user, writer, request)
		if err != nil {
			log.Println(err)
			renderRegister(ctx, writer, http.StatusInternalServerError, registerData{
				Error: helper.InternalServerError().Error(),
				Type:  userType,
			})
			return
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther) // See RFC 2616 (redirect after POST)
	}
}

func NewHandler(authManager *auth.Manager) *registerHandler {
	return &registerHandler{
		authManager: authManager,
	}
}
