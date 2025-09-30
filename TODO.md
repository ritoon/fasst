# TP – Découverte du routeur Echo en Go

Ce TP vous fait découvrir les bases du framework [Echo](https://echo.labstack.com/) en 10 étapes.

---

## 1. Initialisation & Hello World

**Objectif :** démarrer un serveur Echo.

```bash
mkdir tp-echo && cd tp-echo
go mod init tp-echo
go get github.com/labstack/echo/v4
```

Créez `main.go` :

```go
package main

import (
  "net/http"

  "github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  e.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, Echo!")
  })
  e.Logger.Fatal(e.Start(":1323"))
}
```

Lancez :

```bash
go run .
```

Puis ouvrez [http://localhost:1323](http://localhost:1323).

---

## 2. Routes de base (GET / POST / PUT / DELETE)

**Objectif :** manipuler les différentes méthodes HTTP.

```go
e.POST("/users", func(c echo.Context) error { return c.String(201, "created") })
e.GET("/users/:id", func(c echo.Context) error { return c.String(200, c.Param("id")) })
e.PUT("/users/:id", func(c echo.Context) error { return c.NoContent(204) })
e.DELETE("/users/:id", func(c echo.Context) error { return c.NoContent(204) })
```

---

## 3. Params & Query

**Objectif :** récupérer des paramètres d’URL et de query string.

```go
e.GET("/show", func(c echo.Context) error {
  return c.String(200, "q="+c.QueryParam("q"))
})
```

Exemples :

- `http://localhost:1323/show?q=test` → affiche `q=test`

---

## 4. Middleware essentiels (Logger, Recover)

**Objectif :** journaliser les requêtes et éviter les crashs.

```go
import "github.com/labstack/echo/v4/middleware"

e.Use(middleware.Logger())
e.Use(middleware.Recover())
```

---

## 5. Groupes & Auth basique

**Objectif :** créer un préfixe `/admin` et protéger avec Basic Auth.

```go
admin := e.Group("/admin")
admin.Use(middleware.BasicAuth(func(u, p string, c echo.Context) (bool, error) {
  return u == "joe" && p == "secret", nil
}))
admin.GET("/dash", func(c echo.Context) error { return c.String(200, "ok") })
```

---

## 6. Binding JSON vers une struct

**Objectif :** décoder un corps JSON proprement.

```go
type UserDTO struct {
  Name  string `json:"name"`
  Email string `json:"email"`
}

e.POST("/users", func(c echo.Context) error {
  var u UserDTO
  if err := c.Bind(&u); err != nil {
    return c.String(400, "bad request")
  }
  return c.JSON(201, u)
})
```

Tester avec :

```bash
curl -X POST :1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@ex"}'
```

---

## 7. Fichiers statiques

**Objectif :** servir des fichiers statiques.

```go
e.Static("/static", "public")
```

Créer un dossier `public/` avec un fichier `hello.txt`, puis ouvrir :
[http://localhost:1323/static/hello.txt](http://localhost:1323/static/hello.txt)

---

## 8. Templates HTML

**Objectif :** rendre des pages HTML dynamiques.

1. Implémentez un `Renderer` basé sur `html/template`.
2. Attribuez-le à `e.Renderer`.
3. Utilisez-le :

```go
e.GET("/hello/:name", func(c echo.Context) error {
  return c.Render(200, "hello", map[string]any{"Name": c.Param("name")})
})
```

Exemple de template `views/hello.html` :

```html
{{define "hello"}}Hello, {{.Name}}!{{end}}
```

---

## 9. Gestion centralisée des erreurs

**Objectif :** personnaliser les réponses d’erreur.

```go
e.HTTPErrorHandler = func(err error, c echo.Context) {
  code := 500
  if he, ok := err.(*echo.HTTPError); ok {
    code = he.Code
  }
  _ = c.JSON(code, map[string]any{"error": err.Error()})
}
```

## 10) Endpoint d’upload de fichier (multipart/form-data)

**Objectif :** accepter un fichier via `POST /upload` et le sauvegarder localement.

1. Préparer le dossier :

```go
os.MkdirAll("uploads", 0o755)
```

2. (Optionnel) Limiter la taille :

```go
e.Use(middleware.BodyLimit("10M"))
```

3. Route d’upload :

```go
e.POST("/upload", func(c echo.Context) error {
  fh, err := c.FormFile("file")
  if err != nil {
    return echo.NewHTTPError(400, "paramètre 'file' manquant")
  }
  src, err := fh.Open()
  if err != nil { return err }
  defer src.Close()

  dst, err := os.Create(filepath.Join("uploads", fh.Filename))
  if err != nil { return err }
  defer dst.Close()

  if _, err := io.Copy(dst, src); err != nil { return err }

  return c.JSON(201, map[string]any{
    "filename": fh.Filename,
    "size":     fh.Size,
    "path":     "/uploads/" + fh.Filename,
  })
})
```

4. Test rapide :

```bash
curl -F 'file=@README.md' :1323/upload
```

> Import à ajouter : `io`, `os`, `path/filepath`, et `github.com/labstack/echo/v4/middleware`.

---

## 11) Générer un JWT (token de session) + route protégée

**Objectif :** créer un endpoint `/login` qui renvoie un JWT puis protéger `/api/me`.

1. Dépendance :

```bash
go get github.com/golang-jwt/jwt/v5
```

2. Claims personnalisés :

```go
type JwtCustomClaims struct {
  Name  string `json:"name"`
  Admin bool   `json:"admin"`
  jwt.RegisteredClaims
}
```

3. Endpoint de login (token HS256) :

```go
var jwtSecret = []byte("secret") // à stocker en variable d'env en prod

e.POST("/login", func(c echo.Context) error {
  var u struct {
    Email string `json:"email"`
    Name  string `json:"name"`
  }
  if err := c.Bind(&u); err != nil {
    return echo.NewHTTPError(400, "payload invalide")
  }

  // Démo : "auth" triviale
  if u.Email == "" || u.Name == "" {
    return echo.NewHTTPError(401, "identifiants invalides")
  }

  claims := &JwtCustomClaims{
    Name:  u.Name,
    Admin: true,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
      Issuer:    "tp-echo",
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signed, err := token.SignedString(jwtSecret)
  if err != nil { return err }

  return c.JSON(200, map[string]string{"token": signed})
})
```

4. Groupe protégé par JWT :

```go
api := e.Group("/api")
api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
  SigningKey: jwtSecret,
  Claims:     &JwtCustomClaims{},
}))

api.GET("/me", func(c echo.Context) error {
  userToken := c.Get("user").(*jwt.Token)
  claims := userToken.Claims.(*JwtCustomClaims)
  return c.JSON(200, map[string]any{
    "name":  claims.Name,
    "admin": claims.Admin,
    "exp":   claims.ExpiresAt,
  })
})
```

5. Test :

```bash
# 1) Obtenir un token
TOKEN=$(curl -s -X POST :1323/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"joe@ex","name":"Joe"}' | jq -r .token)

# 2) Appeler la route protégée
curl -H "Authorization: Bearer $TOKEN" :1323/api/me
```

> Imports :
> `github.com/golang-jwt/jwt/v5`, `github.com/labstack/echo/v4/middleware`, `time`.

🎯 **Fin du TP**
