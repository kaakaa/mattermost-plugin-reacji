package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
)

func (p *Plugin) initAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", p.handleInfo).Methods(http.MethodGet)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.Use(checkAuthenticity)
	apiV1.HandleFunc("/reacjis", p.handleGetReacjiList).Methods(http.MethodGet)
	apiV1.HandleFunc("/reacjis/{key:[a-z0-9]+}/confirm", p.handleDeleteReacji).Methods(http.MethodPost)
	return r
}

func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("New request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) handleInfo(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, fmt.Sprintf("Mattermost Reacji Plugin %s\n", p.PluginVersion)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Plugin) handleGetReacjiList(w http.ResponseWriter, r *http.Request) {
	list, err := p.Store.Reacji().Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ret []*reacji.Reacji
	channelID := r.URL.Query().Get("channel_id")
	if channelID != "" {
		for _, e := range list.Reacjis {
			if e.FromChannelID == channelID || e.FromChannelID == FromAllChannelKeyword {
				ret = append(ret, e)
			}
		}
	} else {
		ret = list.Reacjis
	}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		p.API.LogWarn("failed to write reacji list", "error", err.Error(), "channel_id", channelID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Plugin) handleDeleteReacji(w http.ResponseWriter, r *http.Request) {
	var submission model.SubmitDialogRequest
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		p.API.LogError("Failed to decode SubmitDialogRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// FIXME: fix dedicated method for command
	cmdResp, _ := p.remove(r.Header.Get("Mattermost-User-ID"), []string{submission.CallbackId})
	if !strings.HasPrefix(cmdResp.Text, "Reacjis are removed") {
		http.Error(w, cmdResp.Text, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func checkAuthenticity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Mattermost-User-ID") == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
