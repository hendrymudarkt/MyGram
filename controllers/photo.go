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

type Photo2 struct {
	Title  		string 		`validate:"required" json:"title"`
	Caption     string 		`json:"caption"`
	Photo_url  	string 		`validate:"required" json:"photo_url"`
	User_id		int			`json:"user_id"`
	User	  	[]*User2 	`validate:"required,dive,required" json:"user"`
}

type User2 struct {
	Email     string 	`validate:"required,email" json:"email"`
	Username  string 	`validate:"required" json:"username"`
	Password  string 	`validate:"required,min=6" json:"password"`
	Age       float64   `validate:"required,gte=8" json:"age"`
}

type Photo struct {
	Id        	uint   		`json:"id"`
	Title  		string 		`json:"title"`
	Caption     string 		`json:"caption"`
	Photo_url  	string 		`json:"photo_url"`
	User_id		int			`json:"user_id"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	User		User
}

type User struct {
	Id			int 		`json:"-"`
	Email		string		`json:"email"`
	Username	string		`json:"username"`
}

func GetPhotos(c *fiber.Ctx) error {
	var photos []photo
	
    config.Database.Preload(clause.Associations).Find(&photos)

    return c.Status(200).JSON(photos)
}

func GetPhoto(c *fiber.Ctx) error {
    id := c.Params("id")
    var photo models.Photo

    result := config.Database.Preload(clause.Associations).Find(&photo, id)

    if result.RowsAffected == 0 {
        return c.SendStatus(404)
    }

    return c.Status(200).JSON(photo)
}

func CreatePhoto(c *fiber.Ctx) error {
	var validate = validator.New()

	userData := c.Locals("userData").(jwt.MapClaims)
    photo := new(models.Photo)

	userID := userData["id"].(float64)

    if err := c.BodyParser(photo); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	var letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	user := &User2{
		Email: userData["email"].(string),
		Username: userData["username"].(string),
		Password: letterRunes,
		Age: userData["age"].(float64),
	}

	photos := &Photo2{
		Title: photo.Title,
		Caption: photo.Caption,
		Photo_url: photo.Photo_url,
		User: []*User2{user},
	}

	err := validate.Struct(photos)
    if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
    }

	photo.User_id = int(userID)
	
    config.Database.Create(photo)
    return c.Status(201).JSON(fiber.Map{
		"id": photo.Id,
		"title": photo.Title,
		"caption": photo.Caption,
		"photo_url": photo.Photo_url,
		"user_id": userID,
		"created_at": photo.CreatedAt,
	})
}

func UpdatePhoto(c *fiber.Ctx) error {
    var validate = validator.New()

	userData := c.Locals("userData").(jwt.MapClaims)
    photo := new(models.Photo)

	userID := userData["id"].(float64)
	id := c.Params("id")

	check := config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).First(&photo)

    if check.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid user to update photo",
        })
    }

    if err := c.BodyParser(photo); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	var letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	user := &User2{
		Email: userData["email"].(string),
		Username: userData["username"].(string),
		Password: letterRunes,
		Age: userData["age"].(float64),
	}

	photos := &Photo2{
		Title: photo.Title,
		Caption: photo.Caption,
		Photo_url: photo.Photo_url,
		User: []*User2{user},
	}

	err := validate.Struct(photos)
    if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
    }

	photo.User_id = int(userID)

    config.Database.Where("id = ?", id).Updates(photo)
	i, err := strconv.Atoi(id)
    if err != nil {
        panic(err)
    }

    return c.Status(200).JSON(fiber.Map{
		"id": i,
		"title": photo.Title,
		"caption": photo.Caption,
		"photo_url": photo.Photo_url,
		"user_id": photo.User_id,
		"updated_at": photo.UpdatedAt,
	})
}

func DeletePhoto(c *fiber.Ctx) error {
    id := c.Params("id")
    var photo models.Photo

    userData := c.Locals("userData").(jwt.MapClaims)
	userID := userData["id"].(float64)
	check := config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).First(&photo)

    if check.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid user to delete photo",
        })
    }

	config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).Delete(&photo)

    return c.Status(200).JSON(fiber.Map{
        "message": "Your photo has been successfully deleted",
    })
}