package controller

import (
	"blogbackend/database"
	"blogbackend/models"
	"blogbackend/util"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	//parse body to data
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("unable to parse body")
	}

	//check passowrd greater than 6
	if len(data["password"].(string)) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "password must be greater than 6 char",
		})
	}

	//check email validation
	if !ValidateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "email incorrect",
		})
	}

	//check email already exist
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "email already exist",
		})
	}

	//insert data to user struct
	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Phone:     data["phone"].(string),
	}
	user.SetPassword(data["password"].(string))

	//insert user to DB
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account successfully created!",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	//parse body to data
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("unable to parse body")
	}

	var user models.User

	database.DB.Where("email  = ?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "email address doesnt exist!",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "wrong password",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "succesfully log in",
		"user":    user,
	})

}

type Claims struct {
	jwt.StandardClaims
}
