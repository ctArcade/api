package api

import (
	"errors"
	"pgdb-test/db"
	"pgdb-test/util"

	"github.com/gofiber/fiber/v2"
)

func postMessage(c *fiber.Ctx) error {
	type NewMessage struct {
		AuthorName string
		Content    string `json:"content"`
	}

	u := new(NewMessage)

	if err := c.BodyParser(u); err != nil {
		return util.Error(c, fiber.StatusBadRequest, err)
	}

	u.AuthorName = util.GetSession(c)["username"].(string)

	message, err := client.Message.CreateOne(
		db.Message.Content.Set(u.Content),
		db.Message.Author.Link(
			db.User.Username.Equals(u.AuthorName),
		),
	).With(
		db.Message.Author.Fetch(),
	).Exec(ctx)

	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(message)
}

func getMessages(c *fiber.Ctx) error {
	ms, err := client.Message.FindMany(
		db.Message.Content.Contains(c.Params("message")),
	).With(
		db.Message.Author.Fetch(),
	).Exec(ctx)

	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(ms)
}

func deleteMessage(c *fiber.Ctx) error {
	username := util.GetSession(c)["username"].(string)

	txn := client.Message.FindMany(
		db.Message.And(
			db.Message.ID.Equals(c.Params("messageID")),
			db.Message.AuthorName.Equals(username),
		),
	)

	msg, err := txn.Exec(ctx)

	if len(msg) == 0 {
		return util.Error(c, fiber.StatusBadRequest, errors.New("make sure the message exists and is owned by you"))
	} else if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	_, err = txn.Delete().Exec(ctx)

	if err != nil {
		return util.Error(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"message": "message deleted successfully"})
}
