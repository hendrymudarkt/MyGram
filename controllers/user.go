package controllers

import (
	"MyGram/config"
	"MyGram/helpers"
	"MyGram/models"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetUsers(c *fiber.Ctx) error {
    var users []models.User

    config.Database.Find(&users)
    return c.Status(200).JSON(users)
}

func GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User

    result := config.Database.Find(&user, id)

    if result.RowsAffected == 0 {
        return c.SendStatus(404)
    }

    return c.Status(200).JSON(&user)
}

func CreateUser(c *fiber.Ctx) error {
    var validate = validator.New()
    user := new(models.User)

	user.Password = helpers.HashPass(user.Password)

    if err := c.BodyParser(user); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	errors := validate.Struct(user)
    if errors != nil {
       return c.Status(fiber.StatusBadRequest).JSON(errors)
        
    }
	
    config.Database.Create(user)
    return c.Status(201).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
    userData := c.Locals("userData").(jwt.MapClaims)

    user := new(models.User)
    
    userID := uint(userData["id"].(float64))
    id := c.Params("id")

    if err := c.BodyParser(user); err != nil {
        return c.Status(503).SendString(err.Error())
    }

    user.Id = userID

    config.Database.Where("id = ?", id).Updates(user)

    return c.Status(200).JSON(fiber.Map{
        "id": user.Id,
        "email": user.Email,
        "username": user.Username,
        "age": user.Age,
        "updated_at": user.UpdatedAt,
    })
}

func DeleteUser(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User

    result := config.Database.Delete(&user, id)

    if result.RowsAffected == 0 {
        return c.SendStatus(404)
    }

    return c.Status(200).JSON(fiber.Map{
        "message": "Your account has been successfully deleted",
    })
}

func LoginUser(c *fiber.Ctx) error {
    user := new(models.User)
    password := user.Password

    if err := c.BodyParser(user); err != nil {
        return c.Status(503).SendString(err.Error())
    }

    err := config.Database.Debug().Where("email = ?", user.Email).Take(user).Error

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid username/email",
        })
    }

    comparePass := helpers.ComparePass([]byte(user.Password), []byte(password))

    if !comparePass {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid password",
        })
    }

    token := helpers.GenerateToken(int(user.Id), int(user.Age), user.Email, user.Username)

    return c.JSON(fiber.Map{
        "token": token,
    })
}