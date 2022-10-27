package controllers

import (
	"MyGram/config"
	"MyGram/models"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"
)

type comments struct {
    Message string
}

type comment struct {
	Id        	uint   		`json:"id"`
	Message     string 		`json:"message"`
	Photo_id	int 		`json:"photo_id"`
	User_id	  	int 		`json:"user_id"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	CreatedAt	time.Time 	`json:"created_at"`
	User	  	user
	Photo		photo
}

type user struct {
	Id			int 		`json:"id"`
	Email		string		`json:"email"`
	Username	string		`json:"username"`
}

type photo struct {
	Id        	uint   		`json:"id"`
	Title  		string 		`json:"title"`
	Caption     string 		`json:"caption"`
	Photo_url  	string 		`json:"photo_url"`
	User_id		int			`json:"user_id"`
}

func GetComments(c *fiber.Ctx) error {
	var comments []comment
	
    config.Database.Preload(clause.Associations).Find(&comments)

    return c.Status(200).JSON(comments)
}

func GetComment(c *fiber.Ctx) error {
    id := c.Params("id")
    var comment models.Comment

    result := config.Database.Preload(clause.Associations).Find(&comment, id)

    if result.RowsAffected == 0 {
        return c.SendStatus(404)
    }

    return c.Status(200).JSON(comment)
}

func CreateComment(c *fiber.Ctx) error {
	var validate = validator.New()

	userData := c.Locals("userData").(jwt.MapClaims)
    comment := new(models.Comment)

	userID := userData["id"].(float64)

    if err := c.BodyParser(comment); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	comments := &comments{
		Message: comment.Message,
	}

	err := validate.Struct(comments)
    if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
    }

	comment.User_id = int(userID)
	
    config.Database.Create(comment)
    return c.Status(201).JSON(fiber.Map{
		"id": comment.Id,
		"message": comment.Message,
		"photo_id": comment.Photo_id,
		"user_id": comment.User_id,
		"created_at": comment.CreatedAt,
	})
}

func UpdateComment(c *fiber.Ctx) error {
    var validate = validator.New()

	userData := c.Locals("userData").(jwt.MapClaims)
    comment := new(models.Comment)

	userID := userData["id"].(float64)
	id := c.Params("id")

	check := config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).First(&comment)

    if check.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid user to update comment",
        })
    }

    if err := c.BodyParser(comment); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	comments := &comments{
		Message: comment.Message,
	}

	err := validate.Struct(comments)
    if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
    }

	comment.User_id = int(userID)

    config.Database.Where("id = ?", id).Updates(comment)
	i, err := strconv.Atoi(id)
    if err != nil {
        panic(err)
    }

    return c.Status(200).JSON(fiber.Map{
		"id": i,
		"message": comment.Message,
		"photo_id": comment.Photo_id,
		"user_id": comment.User_id,
		"updated_at": comment.UpdatedAt,
	})
}

func DeleteComment(c *fiber.Ctx) error {
    id := c.Params("id")
    var comment models.Comment

    userData := c.Locals("userData").(jwt.MapClaims)
	userID := userData["id"].(float64)
	check := config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).First(&comment)

    if check.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid user to delete comment",
        })
    }

	config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).Delete(&comment)

    return c.Status(200).JSON(fiber.Map{
        "message": "Your comment has been successfully deleted",
    })
}