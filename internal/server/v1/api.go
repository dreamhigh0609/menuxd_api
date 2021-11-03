package v1

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/super-radmir/menuxd_api/internal/storage"
)

// NewAPI returns the API V1 Handler with configuration.
func NewAPI() (http.Handler, error) {
	if err := storage.InitData(); err != nil {
		return nil, err
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)
	r := chi.NewRouter()

	um, ur := NewUserRouter(storage.UserStorage{})

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		Post("/login", ur.LoginHandler)
	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		Put("/forgot-password/", ur.ForgotPasswordHandler)
	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		Mount("/users", um)

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		Mount("/dishes", NewDishRouter(storage.DishStorage{}))
	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		Mount("/clients", NewClientRouter(storage.ClientStorage{}))

	r.Mount("/orders", NewOrderRouter(
		storage.OrderStorage{},
		storage.TableStorage{},
		storage.DishStorage{},
	))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		Mount("/categories", NewCategoryRouter(storage.CategoryStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/waiters", NewWaiterRouter(storage.WaiterStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/pictures", NewPictureRouter())

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/tables", NewTableRouter(storage.TableStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/bills", NewBillRouter(storage.BillStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/promotions", NewPromotionRouter(storage.PromotionStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/ads", NewAdRouter(storage.AdStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/ratings", NewRatingRouter(storage.RatingStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/questions", NewQuestionRouter(storage.QuestionStorage{}))

	r.With(middleware.DefaultCompress).
		With(middleware.Timeout(10*time.Second)).
		With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Mount("/stay", NewStayRouter(storage.StayStorage{}))

	return r, nil
}
