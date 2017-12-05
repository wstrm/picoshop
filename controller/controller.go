package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/controller/admin"
	"github.com/willeponken/picoshop/controller/article"
	"github.com/willeponken/picoshop/controller/cart"
	"github.com/willeponken/picoshop/controller/home"
	"github.com/willeponken/picoshop/controller/login"
	"github.com/willeponken/picoshop/controller/logout"
	"github.com/willeponken/picoshop/controller/register"
	"github.com/willeponken/picoshop/controller/search"
	"github.com/willeponken/picoshop/controller/static"
	"github.com/willeponken/picoshop/controller/user"
	"github.com/willeponken/picoshop/controller/warehouse"
	"github.com/willeponken/picoshop/middleware/auth"
	"github.com/willeponken/picoshop/model"
)

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

	mux.Handle("/", a.Middleware(
		home.NewHandler(), openPolicy))

	mux.Handle("/login", a.Middleware(
		login.NewHandler(a), openPolicy))

	mux.Handle("/logout", a.Middleware(
		logout.NewHandler(a), openPolicy))

	mux.Handle("/register", a.Middleware(
		register.NewHandler(a), openPolicy))

	mux.Handle("/user", a.Middleware(
		user.NewHandler(), userPolicy))

	mux.Handle("/article", a.Middleware(
		article.NewHandler(), openPolicy))

	mux.Handle("/cart", a.Middleware(
		cart.NewHandler(), userPolicy))

	mux.Handle("/admin/", a.Middleware(
		http.StripPrefix("/admin",
			admin.NewMux()), adminPolicy))

	mux.Handle("/warehouse", a.Middleware(
		warehouse.NewHandler(), employeePolicy))

	mux.Handle("/search", a.Middleware(
		search.NewHandler(), openPolicy))

	mux.Handle("/static/", http.StripPrefix("/static", static.NewHandler()))

	return mux
}
