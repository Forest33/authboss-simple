package main

import (
	"context"
	"fmt"
	"github.com/aarondl/tpl"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	"github.com/volatiletech/authboss"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	_ "github.com/volatiletech/authboss/auth"
	"github.com/volatiletech/authboss/defaults"
	_ "github.com/volatiletech/authboss/logout"
	"github.com/volatiletech/authboss/remember"
	"net/http"
	"time"
)

type WebServer struct {
	config       *Config
	ab           *authboss.Authboss
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
	templates    tpl.Templates
}

func StartWebServer(config *Config) *WebServer {
	ws := &WebServer{
		config: config,
		ab:     authboss.New(),
	}
	done := make(chan bool)
	go func() {
		ws.start(done)
	}()
	<-done
	return ws
}

func (ws *WebServer) initAuthboss() {
	ws.ab.Config.Paths.RootURL = fmt.Sprintf("http://%s:%d", ws.config.WebServerHost, ws.config.WebServerPort)
	ws.ab.Config.Paths.Mount = "/"
	ws.ab.Config.Modules.LogoutMethod = "GET"

	ws.ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/login", "views")
	ws.ab.Config.Core.MailRenderer = abrenderer.NewEmail("/", "")

	ws.ab.Config.Storage.Server = NewStorer(ws.config)
	ws.ab.Config.Storage.SessionState = ws.sessionStore
	ws.ab.Config.Storage.CookieState = ws.cookieStore

	defaults.SetCore(&ws.ab.Config, false, true)

	ws.ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON:    false,
		UseUsername: true,
	}

	if err := ws.ab.Init(); err != nil {
		panic(err)
	}
}

func (ws *WebServer) start(done chan bool) {
	ws.templates = tpl.Must(tpl.Load("views", "", "layout.html.tpl", nil))

	ws.cookieStore = abclientstate.NewCookieStorer([]byte(ws.config.CookieStoreKey), nil)
	ws.cookieStore.HTTPOnly = false
	ws.cookieStore.Secure = false

	ws.sessionStore = abclientstate.NewSessionStorer(ws.config.SessionCookieName, []byte(ws.config.SessionStoreKey), nil)
	cstore := ws.sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((time.Duration(ws.config.SessionMaxAge) * 24 * time.Hour) / time.Second))

	ws.initAuthboss()

	router := mux.NewRouter()
	router.Use(ws.logger, ws.nosurfing, ws.ab.LoadClientStateMiddleware, ws.auth, remember.Middleware(ws.ab), ws.dataInjector)
	router.HandleFunc("/foo", ws.fooEndpoint).Methods("GET")
	router.HandleFunc("/bar", ws.barEndpoint).Methods("GET")
	router.HandleFunc("/sigma", ws.sigmaEndpoint).Methods("GET")
	router.HandleFunc("/", ws.indexEndpoint).Methods("GET")
	router.Handle("/login", ws.ab.Config.Core.Router)
	router.Handle("/logout", ws.ab.Config.Core.Router)

	done <- true

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ws.config.WebServerHost, ws.config.WebServerPort), router)
	if err != nil {
		panic(err)
	}
}

func (ws *WebServer) dataInjector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := ws.layoutData(w, &r)
		r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
		handler.ServeHTTP(w, r)
	})
}

func (ws *WebServer) layoutData(w http.ResponseWriter, r **http.Request) authboss.HTMLData {
	currentUserName := ""
	userRole := ROLE_USER
	userInter, err := ws.ab.LoadCurrentUser(r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
		userRole = userInter.(*User).Role
	}

	return authboss.HTMLData{
		"logged":            userInter != nil,
		"current_user_name": currentUserName,
		"csrf_token":        nosurf.Token(*r),
		"admin_role":        userRole == ROLE_ADMIN,
	}
}

func (ws *WebServer) fooEndpoint(w http.ResponseWriter, r *http.Request) {
	ws.renderPage(w, r, "foo", authboss.HTMLData{})
}

func (ws *WebServer) barEndpoint(w http.ResponseWriter, r *http.Request) {
	ws.renderPage(w, r, "bar", authboss.HTMLData{})
}

func (ws *WebServer) sigmaEndpoint(w http.ResponseWriter, r *http.Request) {
	ws.renderPage(w, r, "sigma", authboss.HTMLData{})
}

func (ws *WebServer) indexEndpoint(w http.ResponseWriter, r *http.Request) {
	ws.renderPage(w, r, "index", authboss.HTMLData{})
}

func (ws *WebServer) renderPage(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	var current authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)
	if dataIntf == nil {
		current = authboss.HTMLData{}
	} else {
		current = dataIntf.(authboss.HTMLData)
	}

	current.MergeKV("csrf_token", nosurf.Token(r))
	current.Merge(data)

	err := ws.templates.Render(w, name, current)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintln(w, "Error occurred rendering template:", err)
}
