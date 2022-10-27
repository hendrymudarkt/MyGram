package controllers

import (
	"MyGram/config"
	"MyGram/models"
	"encoding/json"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"
)

type socialmedias struct {
    Name 			string
	SocialMediaUrl 	string
}

type social_media struct {
	Id        		uint   		`json:"id"`
	Name     		string 		`json:"name"`
	SocialMediaUrl  string 		`json:"social_media_url"`
	User_id			int			`json:"user_id"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
	Users	  		Users		`gorm:"foreignKey:User_id" json:"User"`
}

type Users struct {
	Id			int 		`json:"id"`
	Username	string		`json:"username"`
	Photos		*Photos		`gorm:"foreignKey:Id;references:User_id" json:"profile_image_url"`
}

type Photos struct {
	User_id		int 		`json:"-"`
	Photo_url	string		`json:"photo_url"`
}

type Response struct {
    Data []social_media `json:"social_medias"`
}

func GetSocialMedias(c *fiber.Ctx) error {
	var socialmedias []social_media
	
    config.Database.Preload("Users.Photos").Find(&socialmedias)

	resp := Response{Data: socialmedias}
	return c.Status(200).JSON(resp)
}

func (u *Users) MarshalJSON() ([]byte, error) {
    intermediate := map[string]interface{}{
        "id":   u.Id,
        "username": u.Username,
        "profile_image_url": u.Photos.Photo_url,
    }
    return json.Marshal(intermediate)
}

func GetSocialMedia(c *fiber.Ctx) error {
    id := c.Params("id")
    var photo models.Photo

    result := config.Database.Preload(clause.Associations).Find(&photo, id)

    if result.RowsAffected == 0 {
        return c.SendStatus(404)
    }

    return c.Status(200).JSON(photo)
}

func CreateSocialMedia(c *fiber.Ctx) error {
	var validate = validator.New()

	userData := c.Locals("userData").(jwt.MapClaims)
    socialmedia := new(models.SocialMedia)

	userID := userData["id"].(float64)

    if err := c.BodyParser(socialmedia); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	socialmedias := &socialmedias{
		Name: socialmedia.Name,
		SocialMediaUrl: socialmedia.SocialMediaUrl,
	}

	err := validate.Struct(socialmedias)
    if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
    }

	socialmedia.User_id = int(userID)
	
    config.Database.Create(socialmedia)
    return c.Status(201).JSON(fiber.Map{
		"id": socialmedia.Id,
		"name": socialmedia.Name,
		"social_media_url": socialmedia.SocialMediaUrl,
		"user_id": userID,
		"created_at": socialmedia.CreatedAt,
	})
}

func UpdateSocialMedia(c *fiber.Ctx) error {
    var validate = validator.New()

	userData := c.Locals("userData").(jwt.MapClaims)
    socialmedia := new(models.SocialMedia)

	userID := userData["id"].(float64)
	id := c.Params("id")

	check := config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).First(&socialmedia)

    if check.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid user to update social media",
        })
    }

    if err := c.BodyParser(socialmedia); err != nil {
        return c.Status(503).SendString(err.Error())
    }

	socialmedias := &socialmedias{
		Name: socialmedia.Name,
		SocialMediaUrl: socialmedia.SocialMediaUrl,
	}

	err := validate.Struct(socialmedias)
    if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
    }

	socialmedia.User_id = int(userID)

    config.Database.Where("id = ?", id).Updates(socialmedia)
    if err != nil {
        panic(err)
    }

    return c.Status(200).JSON(fiber.Map{
		"id": socialmedia.Id,
		"name": socialmedia.Name,
		"social_media_url": socialmedia.SocialMediaUrl,
		"user_id": userID,
		"created_at": socialmedia.UpdatedAt,
	})
}

func DeleteSocialMedia(c *fiber.Ctx) error {
    id := c.Params("id")
    var socialmedia models.SocialMedia

    userData := c.Locals("userData").(jwt.MapClaims)
	userID := userData["id"].(float64)
	check := config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).First(&socialmedia)

    if check.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{
            "error": "Unauthorized",
            "msg":   "invalid user to delete social media",
        })
    }

	config.Database.Debug().Where("id = ?", id).Where("user_id = ?", userID).Delete(&socialmedia)

    return c.Status(200).JSON(fiber.Map{
        "message": "Your social media has been successfully deleted",
    })
}