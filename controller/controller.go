package controller

import (
	"net/http"

	"github.com/wstrm/picoshop/controller/admin"
	"github.com/wstrm/picoshop/controller/article"
	"github.com/wstrm/picoshop/controller/cart"
	"github.com/wstrm/picoshop/controller/category"
	"github.com/wstrm/picoshop/controller/home"
	"github.com/wstrm/picoshop/controller/login"
	"github.com/wstrm/picoshop/controller/logout"
	"github.com/wstrm/picoshop/controller/order"
	"github.com/wstrm/picoshop/controller/register"
	"github.com/wstrm/picoshop/controller/search"
	"github.com/wstrm/picoshop/controller/static"
	"github.com/wstrm/picoshop/controller/user"
	"github.com/wstrm/picoshop/controller/warehouse"
	"github.com/wstrm/picoshop/middleware/auth"
	c "github.com/wstrm/picoshop/middleware/category"
	"github.com/wstrm/picoshop/model"
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

func getCustomerPolicy() auth.Policy {
	policy := auth.NewPolicy()

	policy.SetProtected(true)

	policy.SetUser(model.Admin{}, false)
	policy.SetUser(model.Warehouse{}, false)
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
	openPolicy := getOpenPolicy()         // *
	adminPolicy := getAdminPolicy()       // A
	userPolicy := getUserPolicy()         // A, W, C
	customerPolicy := getCustomerPolicy() // C

	mux.Handle("/", c.Middleware(a.Middleware(
		home.NewHandler(), openPolicy)))

	mux.Handle("/login", c.Middleware(a.Middleware(
		login.NewHandler(a), openPolicy)))

	mux.Handle("/logout", c.Middleware(a.Middleware(
		logout.NewHandler(a), openPolicy)))

	mux.Handle("/register", c.Middleware(a.Middleware(
		register.NewHandler(a), openPolicy)))

	mux.Handle("/user", c.Middleware(a.Middleware(
		user.NewHandler(a), userPolicy)))

	mux.Handle("/order", c.Middleware(a.Middleware(
		order.NewHandler(), customerPolicy)))

	mux.Handle("/article", c.Middleware(a.Middleware(
		article.NewHandler(), openPolicy)))

	mux.Handle("/cart", c.Middleware(a.Middleware(
		cart.NewHandler(), customerPolicy)))

	mux.Handle("/admin/", c.Middleware(a.Middleware(
		http.StripPrefix("/admin",
			admin.NewMux()), adminPolicy)))

	mux.Handle("/warehouse/", c.Middleware(a.Middleware(
		http.StripPrefix("/warehouse",
			warehouse.NewMux()), employeePolicy)))

	mux.Handle("/search", c.Middleware(a.Middleware(
		search.NewHandler(), openPolicy)))

	mux.Handle("/category/", c.Middleware(a.Middleware(
		http.StripPrefix("/category",
			category.NewHandler()), openPolicy)))

	mux.Handle("/static/", http.StripPrefix("/static", static.NewHandler()))

	return mux
}
