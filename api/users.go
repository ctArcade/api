package api

import (
	"errors"
	"pgdb-test/db"
	"time"

	"pgdb-test/util"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var token *jwt.Token

func registerUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	u := new(NewUser)

	if err := c.BodyParser(u); err != nil {
		return util.Error(c, fiber.StatusBadRequest, err)
	}

	_, err := client.User.FindUnique(
		db.User.Username.Equals(u.Username),
	).Exec(ctx)

	if err == db.ErrNotFound {
	} else if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	} else {
		return util.Error(c, fiber.StatusBadRequest, errors.New("user with that username already exists"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	user, err := client.User.CreateOne(
		db.User.Username.Set(u.Username),
		db.User.Password.Set(hashedPassword),
	).Exec(ctx)
	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(user)
}

func loginUser(c *fiber.Ctx) error {
	type Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	a := new(Auth)

	if err := c.BodyParser(a); err != nil {
		return util.Error(c, fiber.StatusBadRequest, err)
	}

	user, err := client.User.FindUnique(
		db.User.Username.Equals(a.Username),
	).Exec(ctx)

	if err == db.ErrNotFound {
		return util.Error(c, fiber.StatusBadRequest, errors.New("user with that username does not exist"))
	} else if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(a.Password))
	if err != nil {
		return util.Error(c, fiber.StatusUnauthorized, err)
	}

	token = jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["isAdmin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"token": t})
}

func getUsers(c *fiber.Ctx) error {
	users, err := client.User.FindMany(
		db.User.Username.Contains(c.Params("username")),
	).With(
		db.User.Messages.Fetch(),
	).Exec(ctx)

	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(users)
}

func whoami(c *fiber.Ctx) error {
	username := util.GetSession(c)["username"].(string)

	user, err := client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Messages.Fetch(),
	).Exec(ctx)

	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(user)
}

// deleting users messed up relation with messages

// func delU(u string) error {
// 	_, err := client.User.FindUnique(
// 		db.User.Username.Equals(u),
// 	).Delete().Exec(ctx)

// 	if err != nil {
// 		return err
// 	}

// 	token.Valid = false

// 	return nil
// }

// func deleteUser(c *fiber.Ctx) error {
// 	var (
// 		targetUser = c.Params("username")
// 		ss         = util.GetSession(c)
// 		isAdmin    = ss["isAdmin"].(bool)
// 		username   = ss["username"].(string)
// 	)

// 	if targetUser == "" {
// 		targetUser = username
// 	}

// 	if targetUser != username && !isAdmin {
// 		return util.Error(c, fiber.StatusUnauthorized, errors.New("only admins can delete users other then their own"))
// 	}

// 	err := delU(targetUser)
// 	if err != nil {
// 		return util.Error(c, fiber.StatusInternalServerError, err)
// 	}

// 	return c.JSON(fiber.Map{"message": "user deleted successfully"})
// }
