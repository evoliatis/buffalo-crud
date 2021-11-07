package actions

import (
	"net/http"
	"todo/models"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

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

// UserCreate default implementation.
func UserCreate(c buffalo.Context) error {
	// Create an empty receive users
	user := models.User{}
	//send an user
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/create.html"))
}

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
}

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
