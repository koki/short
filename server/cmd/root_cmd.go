package cmd

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	//	"github.com/denizeren/dynamostore"
	"github.com/golang/glog"
	"github.com/google/go-github/github"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	"github.com/koki/short/server/logging"
	"github.com/kr/pretty"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
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

type Env struct {
	CookieAuthenticationKey string
	StripeKey               string
	GithubClientID          string
	GithubClientSecret      string
}

const (
	envCookieAuthKey      = "COOKIE_AUTH_KEY"
	envStripeKey          = "STRIPE_KEY"
	envGithubClientID     = "GITHUB_CLIENT_ID"
	envGithubClientSecret = "GITHUB_CLIENT_SECRET"
)

func envFromEnv() Env {
	env := Env{
		CookieAuthenticationKey: os.Getenv(envCookieAuthKey),
		StripeKey:               os.Getenv(envStripeKey),
		GithubClientID:          os.Getenv(envGithubClientID),
		GithubClientSecret:      os.Getenv(envGithubClientSecret),
	}

	if len(env.CookieAuthenticationKey) == 0 {
		glog.Fatalf("missing %s", envCookieAuthKey)
	}

	if len(env.StripeKey) == 0 {
		glog.Fatalf("missing %s", envStripeKey)
	}

	if len(env.GithubClientID) == 0 {
		glog.Fatalf("missing %s", envGithubClientID)
	}

	if len(env.GithubClientSecret) == 0 {
		glog.Fatalf("missing %s", envGithubClientSecret)
	}

	return env
}

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

	// Keys setup
	env := envFromEnv()
	store = sessions.NewCookieStore([]byte(env.CookieAuthenticationKey))
	stripe.Key = env.StripeKey
	oauthCfg = &oauth2.Config{
		ClientID:     env.GithubClientID,
		ClientSecret: env.GithubClientSecret,
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
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/oauth/github-callback", githubCallback)

	s := http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: cors.AllowAll().Handler(mux),
	}

	return s.ListenAndServe()
}

func mkErrorLog(err error) logging.Log {
	return logging.Log{
		Error: logging.MkError(err),
	}
}

func mkErrorMsgLog(msg string, args ...interface{}) logging.Log {
	x := fmt.Sprintf(msg, args...)
	return logging.Log{
		Error: &x,
	}
}

func mkInfoLog(msg string, args ...interface{}) logging.Log {
	fullMsg := fmt.Sprintf(msg, args...)
	return logging.Log{
		Other: &fullMsg,
	}
}

// TODO: What's the right way to associate a Stripe customer ID with our user?
//       Keeping it in the sessions DB might not work.
func getStripeCustomerID(sesh *sessions.Session) string {
	if id, ok := sesh.Values["stripe_id"]; ok {
		if id, ok := id.(string); ok {
			return id
		}

		logging.WriteLog(mkErrorLog(fmt.Errorf("%v contains non-string stripe_id", sesh.Values)))
	}

	return ""
}

// TODO: Plan from env.
func getStripeSubscription(stripeCustomerID string) *stripe.Sub {
	if len(stripeCustomerID) == 0 {
		logging.WriteLog(mkInfoLog("no stripe customer ID, so no stripe subscription"))
		return nil
	}

	i := sub.List(&stripe.SubListParams{
		Customer: stripeCustomerID,
		Plan:     "test_plan_0",
	})

	for i.Next() {
		s := i.Sub()
		logging.WriteLog(mkInfoLog("got stripe subscription %v for %s", s, stripeCustomerID))
		return s
	}

	logging.WriteLog(mkInfoLog("no stripe subscription for %s", stripeCustomerID))

	return nil
}

// Returns false if subscription already ended or does not exist.
func setExpiryFromSubscription(sesh *sessions.Session, sub *stripe.Sub) bool {
	if sub == nil {
		logging.WriteLog(mkInfoLog("user doesn't have a subscription"))
		return false
	}
	sesh.Values["subscription_end"] = sub.PeriodEnd

	return sub.PeriodEnd > time.Now().Unix()
}

func checkSubscription(sesh *sessions.Session) bool {
	if id, ok := sesh.Values["id"]; ok {
		// Do we need to check with Stripe right now?
		if subEnd, ok := sesh.Values["subscription_end"]; ok {
			if subEnd, ok := subEnd.(int64); ok {
				if subEnd > time.Now().Unix() {
					// We checked the subscription earlier, and it hasn't expired yet.
					logging.WriteLog(mkInfoLog("previously verified subscription for github ID %v", id))
					return true
				}
			}
		}

		// Check with Stripe.
		customerID := getStripeCustomerID(sesh)
		if len(customerID) == 0 {
			logging.WriteLog(mkInfoLog("no stripe customer ID for github ID %v", id))
			return false
		}

		sub := getStripeSubscription(customerID)
		if sub == nil {
			logging.WriteLog(mkInfoLog("no subscription for stripe customer ID %s, github ID %v", customerID, id))
			return false
		}

		return setExpiryFromSubscription(sesh, sub)
	}

	logging.WriteLog(mkInfoLog("session has no github ID, so can't look for stripe subscriptions"))
	return false
}

func login(rw http.ResponseWriter, r *http.Request) {
	sesh, err := store.Get(r, "user")
	if err != nil {
		logging.WriteLog(mkInfoLog("invalid cookie, setting up a new one for login"))
	}

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	sesh.Values["state"] = state
	sesh.Save(r, rw)

	url := oauthCfg.AuthCodeURL(state)
	http.Redirect(rw, r, url, 302)
}

func mkConvertLog(sessionValues map[interface{}]interface{}, unconverted string, err error) logging.Log {
	return logging.Log{
		Session: logging.MkSession(sessionValues),
		Convert: logging.MkConvert(unconverted),
		Error:   logging.MkError(err),
	}
}

func convert(rw http.ResponseWriter, r *http.Request) {
	defer context.Clear(r)
	sesh, err := store.Get(r, "user")
	if err != nil {
		http.Error(rw, "invalid cookie", http.StatusUnauthorized)
		logging.WriteLog(mkConvertLog(sesh.Values, "", fmt.Errorf("unauthorized")))
		return
	}
	glog.Error(pretty.Sprintf("%# v", sesh.Values))

	_, hasID := sesh.Values["id"]
	if sesh.IsNew || !hasID {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		logging.WriteLog(mkConvertLog(sesh.Values, "", fmt.Errorf("unauthorized")))
		return
	}

	// TODO: Don't ignore eligibility once subscription system is in place.
	_ = checkSubscription(sesh)
	sesh.Save(r, rw)
	/*
		if !eligible {
			http.Error(rw, "no eligible subscription", http.StatusUnauthorized)
			return
		}
	*/

	if r.Method != http.MethodPost {
		http.Error(rw, "", http.StatusMethodNotAllowed)
		logging.WriteLog(mkConvertLog(sesh.Values, "", fmt.Errorf("method not allowed")))
		return
	}

	headers := rw.Header()
	unconvertedBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.WriteLog(mkConvertLog(sesh.Values, "", err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(unconvertedBytes))}
	unconverted := string(unconvertedBytes)

	objs, err := parser.ParseStreams(streams)
	if err != nil {
		logging.WriteLog(mkConvertLog(sesh.Values, unconverted, err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	kokiObjs, err := converter.ConvertToKokiNative(objs)
	if err != nil {
		logging.WriteLog(mkConvertLog(sesh.Values, unconverted, err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
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
		rw.WriteHeader(http.StatusInternalServerError)
		logging.WriteLog(mkConvertLog(sesh.Values, unconverted, err))
		return
	}

	logging.WriteLog(mkConvertLog(sesh.Values, unconverted, nil))
}

func mkLoginLog(sessionValues map[interface{}]interface{}, err error) logging.Log {
	return logging.Log{
		Session: logging.MkSession(sessionValues),
		Login:   &logging.Login{},
		Error:   logging.MkError(err),
	}
}

const closingPage = `
<html>
<script type="text/javascript">
  parent.close();
  window.location = "https://docs.koki.io/short/";
</script>
</html>
`

func githubCallback(rw http.ResponseWriter, r *http.Request) {
	defer context.Clear(r)
	session, err := store.Get(r, "user")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		logging.WriteLog(mkLoginLog(nil, err))
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(rw, "", http.StatusUnauthorized)
		logging.WriteLog(mkLoginLog(nil, err))
		return
	}

	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		logging.WriteLog(mkLoginLog(nil, err))
		return
	}

	if !tkn.Valid() {
		http.Error(rw, "invalid token", http.StatusUnauthorized)
		logging.WriteLog(mkLoginLog(nil, err))
		return
	}

	client := github.NewClient(oauthCfg.Client(oauth2.NoContext, tkn))

	ctx := r.Context()
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		logging.WriteLog(mkLoginLog(nil, err))
		return
	}

	if user.ID != nil {
		session.Values["id"] = *user.ID
	}
	if user.Login != nil {
		session.Values["login"] = *user.Login
	}
	if user.Name != nil {
		session.Values["name"] = *user.Name
	}
	if user.Email != nil {
		session.Values["email"] = user.Email
	}
	session.Values["accessToken"] = tkn.AccessToken
	session.Save(r, rw)

	logging.WriteLog(mkLoginLog(session.Values, nil))

	fmt.Fprint(rw, closingPage)
}
