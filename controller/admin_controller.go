package controller

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/view"
)

type AdminController struct {
	postDao dao.PostDao
	viewSet *view.ViewSet
	config  *config.Config
}

func NewAdminController(pd dao.PostDao, vs *view.ViewSet, c *config.Config) *AdminController {
	return &AdminController{
		postDao: pd,
		viewSet: vs,
		config:  c,
	}
}

func (ac *AdminController) GetIndex(w http.ResponseWriter, r *http.Request) {
	ps, err := ac.postDao.SelectByQuery(&dao.PostQuery{
		Start:   1,
		Results: 200,
	})
	if err != nil {
		ac.fatalResponse(w, err)
		return
	}
	if err := ac.viewSet.AdminView(ps).Render(w); err != nil {
		ac.fatalResponse(w, err)
		return
	}
}

func (c *AdminController) fatalResponse(w http.ResponseWriter, err error) {
	logrus.Error(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
