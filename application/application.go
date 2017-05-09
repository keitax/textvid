package application

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/gorilla/mux"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/controller"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/util"
	"github.com/keitax/textvid/view"
)

func New(config *config.Config) (http.Handler, error) {
	conn, err := dbr.Open("mysql", config.DataSourceName, nil)
	if err != nil {
		return nil, err
	}

	d := dao.NewPostDao(conn, config)

	ub := util.NewUrlBuilder(config)
	vs := view.NewViewSet(ub, config)
	pc := controller.NewPostController(d, vs, ub, config)
	ac := controller.NewAdminController(d, vs, config)

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir))))
	router.HandleFunc("/{year:[0-9]{4}}/{month:0[1-9]|1[0-2]}/{name}.html", pc.GetSingle)
	router.HandleFunc("/posts/{id:[0-9]+}/edit", pc.GetEditor)
	router.HandleFunc("/posts/", pc.GetList)
	router.HandleFunc("/posts/{id:[0-9]+}", pc.EditPost)
	router.HandleFunc("/", pc.GetIndex)
	router.HandleFunc("/admin", ac.GetIndex)

	return router, nil
}
