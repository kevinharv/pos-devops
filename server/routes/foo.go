package routes

import (
	"net/http"

	"github.com/kevinharv/pos-devops/server/utils"
)

type FooStruct struct {
	Foo string `json:"foo"`
}

func FooHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := utils.Encode(w, r, http.StatusOK, &FooStruct{"test"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return;
		}
	}
}