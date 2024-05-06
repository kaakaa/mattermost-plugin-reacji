package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func (p *Plugin) initAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", p.handleInfo).Methods(http.MethodGet)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.Use(checkAuthenticity)
	apiV1.HandleFunc("/reacjis", p.handleGetReacjiList).Methods(http.MethodGet)
	return r
}

func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("New request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) handleInfo(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("Mattermost Reacji Plugin %s\n", p.PluginVersion))
}

func (p *Plugin) handleGetReacjiList(w http.ResponseWriter, r *http.Request) {
	list, err := p.Store.Reacji().Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ret []*reacji.Reacji
	channelId := r.URL.Query().Get("channel_id")
	if channelId != "" {
		for _, e := range list.Reacjis {
			if e.FromChannelID == channelId || e.FromChannelID == FromAllChannelKeyword {
				ret = append(ret, e)
			}
		}
	} else {
		ret = list.Reacjis
	}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		p.API.LogWarn("failed to write reacji list", "error", err.Error(), "channel_id", channelId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
