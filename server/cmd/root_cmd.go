package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"

	//	"github.com/denizeren/dynamostore"
	"github.com/golang/glog"
	"github.com/google/go-github/github"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	RootCmd = &cobra.Command{
		Use:   "short-server",
		Short: "Manageable Kubernetes manifests using koki/short through a REST API",
		Long: `Short converts the api-friendly kubernetes manifests into ops-friendly syntax.

Full documentation available at https://docs.koki.io/short
`,
		RunE: func(c *cobra.Command, args []string) error {
			return serve(c, args)
		},
		SilenceUsage: true,
	}

	// the port on which the server should listen
	port int

	// the interface on which the server should listen
	ip string

	store *sessions.CookieStore

	oauthCfg *oauth2.Config
)

func init() {
	// local flags to root command
	RootCmd.Flags().IntVarP(&port, "port", "p", 8080, "the port on which the server should listen")
	RootCmd.Flags().StringVarP(&ip, "ip", "i", "0.0.0.0", "the interface on which the server should listen")

	// parse the go default flagset to get flags for glog and other packages in future
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// defaulting this to true so that logs are printed to console
	flag.Set("logtostderr", "true")

	//suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})

	RootCmd.AddCommand(versionCmd)

	store = sessions.NewCookieStore([]byte("dummy-secret"))

	oauthCfg = &oauth2.Config{
		ClientID:     "1daa51bf202f482a44b2",
		ClientSecret: "8af958d8c1c846f8270ee1c0a7e43e3881f6022d",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "",
		Scopes:      []string{"read:user"},
	}

}

func serve(c *cobra.Command, args []string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convert)
	//mux.HandleFunc("/oauth/github-callback", githubCallback)

	s := http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: mux,
	}

	return s.ListenAndServe()
}

func login(rw http.ResponseWriter, r *http.Request) {
	sesh, err := store.Get(r, "user")
	if err != nil {
		http.Error(rw, "invalid cookie", http.StatusUnauthorized)
		return
	}

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	sesh.Values["state"] = state
	sesh.Save(r, rw)

	url := oauthCfg.AuthCodeURL(state)
	http.Redirect(rw, r, url, 302)
	return
}

func convert(rw http.ResponseWriter, r *http.Request) {
	defer context.Clear(r)
	/*
			sesh, err := store.Get(r, "user")
			if err != nil {
				http.Error(rw, "invalid cookie", http.StatusUnauthorized)
				return
			}

			id, ok := sesh.Values["id"]

			if sesh.IsNew || !ok {
				http.Error(rw, "unauthorized", http.StatusUnauthorized)
				return
			}

		eligible := checkEligibility(id.(int))
		if !eligible {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		}
	*/

	if r.Method != http.MethodPost {
		http.Error(rw, "", http.StatusMethodNotAllowed)
		return
	}

	headers := rw.Header()
	streams := []io.ReadCloser{r.Body}
	objs, err := parser.ParseStreams(streams)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	kokiObjs, err := converter.ConvertToKokiNative(objs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()
	output := query.Get("output")
	if output == "json" {
		headers.Set("Content-Type", "application/json")
		err = client.WriteObjsToJSONStream(kokiObjs, rw)
	} else {
		headers.Set("Content-Type", "application/yaml")
		err = client.WriteObjsToYamlStream(kokiObjs, rw)
	}
	if err != nil {
		glog.Errorf("Error responding to request: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func githubCallback(rw http.ResponseWriter, r *http.Request) {
	defer context.Clear(r)
	session, err := store.Get(r, "user")
	if err != nil {
		glog.Errorf("Error getting user from store %+v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		glog.Errorf("Error matching state %v", err)
		http.Error(rw, "", http.StatusUnauthorized)
		return
	}

	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		glog.Errorf("")
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !tkn.Valid() {
		http.Error(rw, "invalid token", http.StatusUnauthorized)
		return
	}

	client := github.NewClient(oauthCfg.Client(oauth2.NoContext, tkn))

	ctx := r.Context()
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["name"] = user.Name
	session.Values["id"] = user.ID
	session.Values["accessToken"] = tkn.AccessToken
	session.Save(r, rw)

	http.Redirect(rw, r, "/convert", 302)
}

func checkEligibility(id int) bool {
	return true
}
