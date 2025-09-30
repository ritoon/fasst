### Refacto les handlers de User

Créer un dossier `handler` dans le projet.

Créer un fichier `handler.go` dans le dossier `handler`, nommer le `package handler` et l'ouvrir.

```go
type Handler struct{
    router *echo.Echo
}
```

Dans le fichier main.go supprimer handlers et routes créés sur les uris **"/users"**

Dans le fichier handler/user.go créer les méthodes suivantes:

```go
func (h *Handler)CreateUser(c echo.Context) error{
    var u User
    if err := c.Bind(&u); err != nil {
        return c.String(400, "bad request")
    }
    return c.JSON(201, u)
}

func (h *Handler)UpdateUserByIDUpdateUser(c echo.Context) error{
    return c.NoContent(204)
}

func (h *Handler)GetUserByID(c echo.Context) error{
    return c.String(200, c.Param("id"))
}

func (h *Handler)DeleteUserByID(c echo.Context) error{
    return c.NoContent(204)
}
```

Dans le fichier `handler/handler.go` Créer une fonction qui initialise les routes pour User.

```go
func (h *Handler)initRoutes(){
    h.route.POST("/users", h.CreateUser)
    h.route.GET("/users/:id", h.GetUserByID)
    h.route.PUT("/users/:id", h.UpdateUserByID)
    h.route.DELETE("/users/:id", h.DeleteUserByID)
}
```

Toujours dans le même fichier, créer une fonction constructeur de Handler, pour lui passer le routeur et initialiser les routes de user.

```go
func New (r *echo.Echo) *Handler{
    h := &Handler{router: r}
    h.initRoutes()
    return h
}
```
