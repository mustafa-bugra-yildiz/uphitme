package router

import (
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/view"
)

func (s state) docsSchedulerApiHandler(w http.ResponseWriter, r *http.Request) {
	err := view.Docs.SchedulerApi(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
