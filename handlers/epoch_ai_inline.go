package handlers

import (
	"bytes"
	"github.com/ethpandaops/dora/utils"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func EpochAiInline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	epoch := vars["epoch"]

	name := epoch + "/report.md"

	modtime := time.Now()

	var content, err = os.ReadFile(utils.Config.Alert.RootDir + "/" + name)

	disp := "inline; filename=\"" + epoch + ".svg\""
	if err == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", disp)
		http.ServeContent(w, r, name, modtime, bytes.NewReader(content))
	}
}
