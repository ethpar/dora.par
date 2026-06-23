package handlers

import (
	"bytes"
	"github.com/ethpandaops/dora/utils"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func PoolInline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slot := vars["slot"]

	name := slot + ".pdf"

	modtime := time.Now()

	var content, err = os.ReadFile(utils.Config.Graph.FilesPath + "/" + name)

	if err != nil {
		content, _ = os.ReadFile(utils.Config.Graph.EmptyFile)
	}
	//w.Header().Add("Content-Disposition", "Attachment")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\"yourfile.pdf\"")
	http.ServeContent(w, r, name, modtime, bytes.NewReader(content))
}
