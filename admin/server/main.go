package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	cookieName = "dairycart"
	staticDir  = "assets"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func index(res http.ResponseWriter, req *http.Request) {
	homePage := `
	<!DOCTYPE html>
	<html>
		<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Dairycart</title>
		<link rel="stylesheet" href="/assets/vendor/css/bulma.css">
		<link rel="stylesheet" href="/assets/css/app.css">

		<!-- external dependencies -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css">
		<!--  -->

		</head>
		<body>
		<div class="columns">
			<aside class="column is-2 aside hero is-fullheight is-hidden-mobile">
				<div>
					<div class="main">
					<div class="title">Menu</div>
					<a href="#" class="item"><span class="icon"><i class="fa fa-home"></i></span><span class="name">Dashboard</span></a>
					<a href="#products" class="item"><span class="icon"><i class="fa fa-briefcase"></i></span><span class="name">Products</span></a>
					<a href="#" class="item"><span class="icon"><i class="fa fa-th-list"></i></span><span class="name">Orders</span></a>
					</div>
				</div>
			</aside>
			<div class="column is-10 admin-panel">
				<nav class="nav has-shadow" id="top">
					<div class="container">
						<div class="nav-left">
							<a class="nav-item" href="../index.html">
							<img src="/assets/images/logo.png" alt="Description">Dairycart</a>
						</div>
						<!--
							I don't know what this section accomplishes, but I'm too afraid to delete it
						-->
						<div class="nav-right nav-menu is-hidden-tablet">
							<a href="#" class="nav-item is-active">Dashboard</a>
							<a href="#" class="nav-item">Products</a>
							<a href="#" class="nav-item">Orders</a>
						</div>
					</div>
				</nav>
				<div class="scooted">
					<div id="elm-app"></div>
				</div>
			</div>
		</div>
		</div>
		<script src="/assets/js/elm.js"></script>
		<script>Elm.Dairycart.embed(document.getElementById('elm-app'));</script>
		</body>
	</html>
	`
	fmt.Fprintf(res, homePage)
}

func serveLogin(res http.ResponseWriter, req *http.Request) {
	loginPage := `
	<!DOCTYPE html>
	<html>
		<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Dairycart Login</title>
		<link rel="stylesheet" href="/assets/vendor/css/bulma.css">
		</head>
		<body>
			<section class="hero is-fullheight is-dark is-bold">
				<div class="hero-body">
					<div class="container">
						<div class="columns is-vcentered">
							<div class="column is-4 is-offset-4">
								<h1 class="title">Login</h1>
								<div class="box">
									<label class="label">Username</label>
									<p class="control">
										<input class="input" type="text" placeholder="username">
									</p>
									<label class="label">Password</label>
									<p class="control">
										<input class="input" type="password" placeholder="••••••••••••••••••••••••••••••••••••••••••••••••">
									</p>
									<hr>
									<p class="control">
										<button class="button is-primary">Login</button>
									</p>
								</div>
							</div>
						</div>
					</div>
				</div>
			</section>
		</body>
	</html>
	`
	fmt.Fprintf(res, loginPage)
}

// HTTP middleware setting a value on the request context
func cookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookies := req.Cookies()
		if len(cookies) == 0 {
			http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		}

		for _, c := range cookies {
			if c.Name == cookieName {
				next.ServeHTTP(res, req)
				return
			}
		}
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
	})
}

func main() {
	// debug = strings.ToLower(os.Getenv("DEBUG")) == "true"

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))

	FileServer(r, "/assets/", http.Dir(staticDir))
	r.Get("/login", serveLogin)
	r.Route("/", func(r chi.Router) {
		// r.Use(cookieMiddleware)
		r.Get("/", index)
	})

	port := ":1234"
	log.Printf("server is listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
