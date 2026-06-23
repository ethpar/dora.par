package handlers

import (
	"bytes"
	"github.com/ethpandaops/dora/utils"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func EpochGraph(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	epoch := vars["epoch"]

	name := epoch + "/chain_tree.svg"

	modtime := time.Now()

	var content, err = os.ReadFile(utils.Config.Alert.RootDir + "/" + name)

	if err == nil {
		w.Header().Add("Content-Disposition", "Attachment")
		http.ServeContent(w, r, name, modtime, bytes.NewReader(content))
	}

}
