package Server

import (
	"LostSlot/src/Controllers"
	"LostSlot/src/Entities"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type User Entities.User

type Server struct {
}

func Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	apiRouter := chi.NewRouter()

	apiRouter.Use(middleware.RequestID)
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Recoverer)
	apiRouter.Use(middleware.URLFormat)
	apiRouter.Use(render.SetContentType(render.ContentTypeJSON))

	log.Println("routing user requests")
	apiRouter.Route("/users", func(r chi.Router) {
		//apiRouter.With(Pagination).Get("/", Controllers.GetUsersByPage)
		apiRouter.Post("/", Controllers.CreateUser) // POST /users

		apiRouter.Route("/{userId}", func(r chi.Router) {
			r.Get("/", Controllers.GetUser)
			r.Delete("/", Controllers.DeleteUser)
		})
	})

	log.Println("Mouting api router")
	r.Mount("/api", apiRouter)

	log.Println("Listening on port 3333")
	err := http.ListenAndServe(":3333", r)

	return err
}

//func Pagination(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		PageId := r.URL.Query().Get(string(PageIdKey))
//		intPageId := 0
//		var err error
//		if PageId != "" {
//			intPageId, err = strconv.Atoi(PageId)
//			if err != nil {
//				render.Render(w, r, ErrRender(err))
//				return
//			}
//		}
//		ctx := context.WithValue(r.Context(), PageIdKey, intPageId)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
