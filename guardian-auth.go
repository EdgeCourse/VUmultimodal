/*
another approach: JWT with Gin -demo assumes MySQL
https://codewithmukesh.com/blog/jwt-authentication-in-golang/


The following demo uses Go-Guardian
*/

/*

Let's handle authentication running in cluster mode using Golang and Go-Guardian library.

Let's say we have app replicated twice(A, B) and running behind a load balancer, the user asks for token the load balancer route request to replication A where the token generated and cached in the app memory now the same user request a protected resource and the load balancer route request to replication B this will end up with an authorization error.

When building a modern application, you don’t want to implement authentication module from scratch;
you want to focus on building awesome software. go-guardian is here to help with that.
Here are a few bullet point reasons you might like to try it out:
-provides simple, clean, and idiomatic API.
-provides top trends and traditional authentication methods.
-provides a package to caches the authentication decisions, based on different mechanisms and algorithms.
-provides two-factor authentication and one-time password as defined in RFC-4226 and RFC-6238

*/

//step one: create project

//mkdir scalable-guardian-auth && cd scalable-guardian-auth && go mod init scalable-guardian-auth && touch main.go

//(this demo file is called guardian-auth.go)

/*
install gorilla mux, go-guardian, jwt-go.
go get github.com/gorilla/mux
go get github.com/shaj13/go-guardian
go get "github.com/dgrijalva/jwt-go"


Creating endpoints -demo:
package main
import (
  "github.com/gorilla/mux"
)
func main() {
  router := mux.NewRouter()
}
Establish the endpoints of our API: create all of our endpoints in the main function, every endpoint needs a function to handle the request and we will define those below the main function.

Route handlers: define the functions that will handle the requests.

Set up Go-Guardian

Use HTTP middleware to intercept the request and authenticate users before it reaches the final route.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator
var cache store.Cache

func main() {
	port := os.Getenv("PORT")
	setupGoGuardian()
	router := mux.NewRouter()
	router.HandleFunc("/v1/auth/token", middleware(http.HandlerFunc(createToken))).Methods("GET")
	router.HandleFunc("/v1/book/{id}", middleware(http.HandlerFunc(getBookAuthor))).Methods("GET")
	log.Printf("server started and listening on http://127.0.0.1:%s", port)
	http.ListenAndServe("127.0.0.1:"+port, router)
}

func createToken(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "medium",
		"aud": "any",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	jwtToken, _ := token.SignedString([]byte("secret"))
	w.Write([]byte(jwtToken))
}

func getBookAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	books := map[string]string{
		"1449311601": "Ryan Boyd",
		"148425094X": "Yvonne Wilson",
		"1484220498": "Prabath Siriwarden",
	}
	body := fmt.Sprintf("Author: %s \n", books[id])
	w.Write([]byte(body))
}

func setupGoGuardian() {
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)

	basicStrategy := basic.New(validateUser, cache)
	tokenStrategy := bearer.New(verifyToken, cache)

	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if userName == "medium" && password == "medium" {
		return auth.NewDefaultUser("medium", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("Invalid credentials")
}

func middleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.UserName())
		next.ServeHTTP(w, r)
	})
}

func verifyToken(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := auth.NewDefaultUser(claims["sub"].(string), "", nil, nil)
		return user, nil
	}

	return nil, fmt.Errorf("Invaled token")
}

//curl  -k http://127.0.0.1:8080/v1/auth/token -u medium:medium

//then with whatever token it returns:

//curl  -k http://127.0.0.1:8080/v1/book/66 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2NTc3OTE4OTQsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.Xt25tIe56Mb3ZKvV7jyRM3m_hr9qE_F1Ot0yZsHT1BY%"
