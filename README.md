# Buffalo CRUD exemple

## Introduction

Site du projet : https://gobuffalo.io/fr

* Documentation générale : https://gobuffalo.io/fr/docs/overview/
* Documentation de fizz (model) :  https://gobuffalo.io/en/docs/db/fizz
* Document de plush (vue) : https://gobuffalo.io/fr/docs/rendering/

## Création d'un projet

`brew install buffalo`

`buffalo version`

> INFO[0000] Buffalo version is: v0.17.5

`buffalo new projet --db-type mariadb`

Options possibles :

* `--api` : pas de front
* `--db-type string` : type de base de données

## Configuration d'un projet 

* Modifier le port de buffalo serveur depuis `.env`

~~~bash
PORT=8080
~~~

* Définir les 3 bases de données dans `database.yml`

~~~yaml
development:
  dialect: mariadb
  database: essai_dev
  user: essai
  password: essai
  host: 127.0.0.1

test:
  dialect: mariadb
  database: essai_test
  user: essai
  password: essai
  host: 127.0.0.1

production:
  dialect: mariadb
  database: essai_production
  user: essai
  password: essai
  host: 127.0.0.1
~~~

* Créer les bases de données avec : `buffalo pop create` 

* Lancer le serveur de test avec `buffalo dev`

## Mise en place d'une table

### Génération model /migration

`buffalo pop generate model user` 

Ou avec plus d'options pour générer les champs dans la foulée : `buffalo pop generate model user login password name email age:int`

Création du modèle **user** (pour la la table **users**) et des migrations de la table correspondante. 

Mise en place des champs (j'utilise un entier pour les clefs primaires) dans `migrations/XXX_create_users.up.fizz`

Le système ajoute automatiquement t.Timestamps() (**created_at** et **updated_at**)

~~~sql
create_table("users") {
	t.Column("id", "int", {primary: true, auto_increment: true})
	t.Column("login", "string", {})
	t.Column("password", "string", {})
	t.Column("name", "string", {})
	t.Column("email", "string", {})
	t.Column("age", "integer", {})
	t.Timestamps()
}
~~~

Structure modèle dans `models/user.go`

~~~go
type User struct {
	ID        int       `json:"id" db:"id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password" db:"password"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Age       int       `json:"age" db:"age"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
~~~

> **Attention** : Ajout de l'auto-incrément sur l'id en int dans cet exemple, plutôt qu'un UUID

### Génération dans la base

`buffalo pop migrate` ou `buffalo pop migrate up`

* Pour la suppression des migrations : `buffalo pop migrate down`
* Suivre le status de ses migrations : `buffalo pop migrate status`

## Mise en place des actions

### Génération

Je reprends ici le même fonctionnement que pour Laravel. 

Exemple avec me modèle **user** :

| action | URL  | Method  |Rôle
| :------ | :--------------- | :----- | :-------------------------------- | 
| index   | /users           | GET    | Affiche la liste des utilisateurs |
| create  | /users/create    | GET    | affiche le formulaire de création |
| store   | /users           | POST   | Sauvegarde un nouvel enregistrement |
| show    | /users/{id}      | GET    | Affiche l'utilisateur choisi        |
| edit    | /users/{id}/edit | GET    | Affiche le formulaire de modification |
| update  | /users/{id}      | PUT    | Sauvegarde les modifications |
| destroy | /users/{id}      | DELETE | Supprime l'utilisateur choisi |

`buffalo g actions users index create store show edit update destroy`

Documentation : https://gobuffalo.io/fr/docs/resources

### Gestion des routes

Les routes sont ajoutées dans le fichier `actions/app.go`

~~~go
  // Route pour users
  app.GET("/users", UserIndex)
  app.GET("/users/create", UserCreate)
  app.POST("/users", UserStore)
  app.GET("/users/{id}", UserShow)
  app.GET("/users/{id}/edit/", UserEdit)
  app.PUT("/users/{id}", UserUpdate)
  app.DELETE("/users/{id}", UserDestroy)
~~~

Les actions du controlleur correspondants sont dans : `/actions/user.go`

### Les validateurs

On les retrouve sous la forme de la méthode `Validate` du modèle

~~~go
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Login, Name: "Login"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.IntIsPresent{Field: u.Age, Name: "Age"},
		&validators.IntIsLessThan{Field: u.Age, Name: "Age", Compared: 99, Message: "Age trop grand !"},
	), nil
}
~~~

A compléter

### Methode INDEX

La première méthode est index qui permet de lister l'ensemble des utilisateurs

~~~go
// UserIndex : Affichage des utilisateurs
func UserIndex(c buffalo.Context) error {
	// Create an array to receive users
	users := []models.User{}
	//get all the users from database
	err := models.DB.All(&users)
	// handle any error
	if err != nil {
		c.Flash().Add("error", "users errors !")
		return c.Redirect(301, "/")
	}
	//return list of todos as json
	c.Set("users", users)
	return c.Render(http.StatusOK, r.HTML("users/index.html"))
}
~~~

Utilisation du template plush correspondant dans `templates/user/index.plush.html` avec des liens pour les actions CRUD sur un élément donné

~~~html
<table class="table">
<thead>
  <tr>
    <th scope="col">#</th>
    <th scope="col">Login</th>
    <th scope="col">Name</th>
    <th scope="col">Email</th>
    <th scope="col">Age</th>
    <th scope="col"></th>
  </tr>
  </thead>
  <tbody>
  <%= for (user) in users { %>
  <tr>
    <th scope="row"><%= user.ID %></th>
    <td><%= user.Login %></td>
    <td><%= user.Name %></td>
    <td><%= user.Email %></td>
    <td><%= user.Age %></td>
    <td>
      <%= linkTo(user, {class: "btn btn-sm btn-warning"}) { %>Show<% } %>
      <%= linkTo([user,"edit"], {class: "btn btn-sm btn-primary"}) { %>Edit<% } %>
      <%= form_for( user, {action: userPath({id: user.ID}), method: "DELETE"}) { %>
        <button type="submit" class="btn btn-danger btn-sm" data-confirm="Are you sure?">Delete</button>
      <% } %>
    </td>            
  </tr>
  <% } %>
  </tbody>
</table>
<%= linkTo(["users","create"], {class: "btn btn-sm btn-success active"}) { %>Add User<% } %>
~~~

### Méthode CREATE

L'objectif de ce formulaire est de proposer la création d'un nouvel utilisateur
On configure la méthode `UserCreate()` de l'action 

~~~go
// UserCreate default implementation.
func UserCreate(c buffalo.Context) error {
	// Create an empty receive users
	user := models.User{}
	//send an user
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/create.html"))
}
~~~

On configure la vue create en la séparant en 2 :

Vue principale `templates/users/create.plush.html`

~~~html
<div class="card uper">
    <div class="card-header">
        <i class="fas fa-user"></i> Add a new User
    </div>
    <div class="card-body">
        <!-- Contenu -->
        <%= form_for( user, {method: "POST"}) { %>
            <%= partial("users/form.html") %>
            <hr />
            <a href="<%= usersPath() %>" class="btn btn-dark btn-sm" data-confirm="Are you sure?">Cancel</a>
            <button type="submit" class="btn btn-primary btn-sm">Add</button>
        <% } %>
        <!-- Fin du contenu -->
    </div>
</div>
~~~

Formulaire externalisé dans `templates/users/_form.plush.html` qui contient les différents champs réutilisés pour l'édition

~~~html
<div class="form-group">
    <%= f.InputTag("Login", {"required": true, "type":"text"}) %>
    <%= f.InputTag("Name", {"required": true, "type":"text"}) %>
    <%= f.InputTag("Email", {"required": true, "type":"email"}) %>
    <%= f.InputTag("Age", {"required": true, "type":"number"}) %>
</div>
~~~

La formulaire post ses données à destination de l'action `UserStore`

### Methode STORE

Cette méthode appelée par **CREATE** vérifie que les champs sont correctement remplis puis enregistre l'élément. Une balise Flash de succès est envoyée

 Si l'élément contient une erreur, on revient sur le formulaire précédent avec un message d'erreur Flash sur le formulaire **CREATE**.

Cette action ne nécessite pas de rendu HTML.

~~~go
// UserStore default implementation.
func UserStore(c buffalo.Context) error {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return err
	}
	// Validate the data from the html form
	verrs, err := models.DB.ValidateAndCreate(user)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("user", user)
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		return c.Render(422, r.HTML("users/create.html"))
	}
	c.Flash().Add("success", "User was created successfully")
	return c.Redirect(302, "/users/")
~~~

### Méthode SHOW

Cette méthode d'afficher un enregistrement donné

~~~go
// UserShow default implementation.
func UserShow(c buffalo.Context) error {
	// grab the id url parameter defined in app.go
	id := c.Param("id")
	// create a variable to receive the user
	user := models.User{}
	// grab the user from the database
	err := models.DB.Find(&user, id)
	// handle possible error
	if err != nil {
		c.Flash().Add("warning", "User not found !")
		return c.Redirect(301, "/users")
	}
	//return the data as json
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/show.html"))
}
~~~

L'affichage est généré par un template utilisant la même méchanisme que les formulaires mais en lecture seule : `templates/users/show.plush.html`

~~~html
<div class="card uper">
    <div class="card-header">
        <i class="fas fa-user"></i> Show an User
    </div>
    <div class="card-body">
        <!-- Contenu -->
        <div class="form-group">
        <%= form_for( user, {action: userPath({id: user.ID}), method: "PUT"}) { %>
            <%= f.InputTag("Login", {"readonly": true, "type":"text"}) %>
            <%= f.InputTag("Name", {"readonly": true, "type":"text"}) %>
            <%= f.InputTag("Email", {"readonly": true, "type":"email"}) %>
            <%= f.InputTag("Age", {"readonly": true, "type":"number"}) %>
        <% } %>
        </div>
        <hr />
        <a href="<%= usersPath() %>" class="btn btn-dark btn-sm">Back</a>        
        <!-- Fin du contenu -->
    </div>
</div>
~~~

### la méthode EDIT

Cette méthode réutilise le formulaire externalisé de la méthode **CREATE** mais précharge en amont un enregistrement qu'elle propose à la modification

~~~go
// UserEdit default implementation.
func UserEdit(c buffalo.Context) error {
	// grab the id url parameter defined in app.go
	id := c.Param("id")
	// create a variable to receive the user
	user := models.User{}
	// grab the todo from the database
	err := models.DB.Find(&user, id)
	// handle possible error
	if err != nil {
		c.Flash().Add("warning", "User not found !")
		return c.Redirect(301, "/users")
	}
	//return the data as json
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/edit.html"))
}
~~~

Le template est donc très proche de celui de **CREATE** : `templates/users/edit.plush.html` mais le formulaire appelle une méthode **PUT**

~~~html
<div class="card uper">
    <div class="card-header">
        <i class="fas fa-user"></i> Edit an User
    </div>
    <div class="card-body">
        <!-- Contenu -->
        <%= form_for( user, {action: userPath({id: user.ID}), method: "PUT"}) { %>
            <%= partial("users/form.html") %>
            <hr />
            <a href="<%= usersPath() %>" class="btn btn-dark btn-sm" data-confirm="Are you sure?">Cancel</a>
            
            <button type="submit" class="btn btn-primary btn-sm">Update</button>
        <% } %>
        <!-- Fin du contenu -->
    </div>
</div>
~~~

### La méthode UPDATE

Cette méthode permet d'enregister l'enregistrement modifié après avoir vérifier que les champs sont correctement remplis.

Cette méthode n'utilise pas de template HTML

~~~go
// UserUpdate default implementation.
func UserUpdate(c buffalo.Context) error {

	user := &models.User{}

	if err := models.DB.Find(user, c.Param("id")); err != nil {
		c.Flash().Add("warning", "User not found !")
		return c.Redirect(301, "/users")
	}

	if err := c.Bind(user); err != nil {
		return err
	}
	// Validate the data from the html form
	verrs, err := models.DB.ValidateAndSave(user)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("user", user)
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		return c.Render(422, r.HTML("users/edit.html"))
	}
	c.Flash().Add("success", "User was updated successfully")
	return c.Redirect(302, "/users/")
}
~~~

### La méthode DELETE

Cette méthode procède à la suppression d'un enregistrement. Elle ne nécessite pas de template HTML. Elle est appellée depuis le tableau après avoir demandé une confirmation à l'utilisateur

~~~go
// UserDestroy default implementation.
func UserDestroy(c buffalo.Context) error {

	user := &models.User{}

	if err := models.DB.Find(user, c.Param("id")); err != nil {
		c.Flash().Add("warning", "User not found !")
		return c.Redirect(301, "/users")
	}

	if err := models.DB.Destroy(user); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "User was destroyed successfully")
	return c.Render(200, r.Auto(c, user))
}
~~~

## Compilation de production

Une fois le projet terminé, il est alors nécessaire de le compiler 

`buffalo build`

Notez que le binaire contient l'ensemble des fichiers !

## Autre utilisation

On peut utiliser `buffalo` pour faire un serveur API seul : https://dev.to/alexmercedcoder/api-with-go-buffalo-in-2021-from-zero-to-deploy-5642

`buffalo new projet --db-type mariadb --api`

On fait alors en sorte de ne renvoyer que des contenus JSON :

~~~go
// TodoIndex default implementation.
func TodoIndex(c buffalo.Context) error {
    // Create an array to receive todos
    todos := []models.Todo{}
    //get all the todos from database
    err := models.DB.All(&todos)
    // handle any error
    if err != nil {
        return c.Render(http.StatusOK, r.JSON(err))
    }
    //return list of todos as json
    return c.Render(http.StatusOK, r.JSON(todos))
}

// TodoShow default implementation.
func TodoShow(c buffalo.Context) error {
    // grab the id url parameter defined in app.go
    id := c.Param("id")
    // create a variable to receive the todo
    todo := models.Todo{}
    // grab the todo from the database
    err := models.DB.Find(&todo, id)
    // handle possible error
    if err != nil {
        return c.Render(http.StatusOK, r.JSON(err))
    }
    //return the data as json
    return c.Render(http.StatusOK, r.JSON(&todo))
}


// TodoAdd default implementation.
func TodoAdd(c buffalo.Context) error {
    //get item from url query
    item := c.Param("item")
    //create new instance of todo
    todo := models.Todo{Item: item}
    // Create a fruit without running validations
    err := models.DB.Create(&todo)
    // handle error
    if err != nil {
        return c.Render(http.StatusOK, r.JSON(err))
    }
    //return new todo as json
    return c.Render(http.StatusOK, r.JSON(todo))
}
~~~