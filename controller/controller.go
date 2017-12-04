package controller

import (
	"errors"
	"net/http"

	"github.com/willeponken/picoshop/middleware/auth"
	"github.com/willeponken/picoshop/model"
)

func internalServerError() error {
	return errors.New("Something internal went wrong!")
}

func invalidFormDataError() error {
	return errors.New("Invalid form data")
}

func getEmployeePolicy() auth.Policy {
	policy := auth.NewPolicy()

	policy.SetProtected(true)

	policy.SetUser(model.Admin{}, true)
	policy.SetUser(model.Warehouse{}, true)
	policy.SetUser(model.Customer{}, false)

	return policy
}

func getAdminPolicy() auth.Policy {
	policy := auth.NewPolicy()

	policy.SetProtected(true)

	policy.SetUser(model.Admin{}, true)
	policy.SetUser(model.Warehouse{}, false)
	policy.SetUser(model.Customer{}, false)

	return policy
}

func getUserPolicy() auth.Policy {
	policy := auth.NewPolicy()

	policy.SetProtected(true)

	policy.SetUser(model.Admin{}, true)
	policy.SetUser(model.Warehouse{}, true)
	policy.SetUser(model.Customer{}, true)

	return policy
}

func getOpenPolicy() auth.Policy {
	policy := auth.NewPolicy()

	policy.SetProtected(false)

	return policy
}

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// register model.Admin, model.Customer and model.Warehouse types for authentication
	a := auth.NewManager("auth", model.Admin{}, model.Customer{}, model.Warehouse{})

	employeePolicy := getEmployeePolicy() // A, W
	openPolicy := getOpenPolicy()         // A, W, C, *
	adminPolicy := getAdminPolicy()       // A
	userPolicy := getUserPolicy()         // A, W, C

	mux.Handle("/", a.Middleware(newHomeHandler(), openPolicy))
	mux.Handle("/login", a.Middleware(newLoginHandler(a), openPolicy))
	mux.Handle("/logout", a.Middleware(newLogoutHandler(a), openPolicy))
	mux.Handle("/register", a.Middleware(newRegisterHandler(a), openPolicy))
	mux.Handle("/user", a.Middleware(newUserHandler(), userPolicy))
	mux.Handle("/article", a.Middleware(newArticleHandler(), openPolicy))
	mux.Handle("/cart", a.Middleware(newCartHandler(), userPolicy))
	mux.Handle("/static/", newStaticHandler()) // static does not need to be intercepted with user information
	mux.Handle("/admin", a.Middleware(newAdminHandler(), adminPolicy))
	mux.Handle("/warehouse", a.Middleware(newWarehouseHandler(), employeePolicy))
	mux.Handle("/search", a.Middleware(newSearchHandler(), openPolicy))

	return mux
}
