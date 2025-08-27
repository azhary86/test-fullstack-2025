package main

import (
    "context"
    "crypto/sha1"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
)

type User struct {
    RealName string `json:"realname"`
    Email string `json:"email"`
    Password string `json:"password"` 
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

var rdb *redis.Client
var ctx = context.Background()

func initRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr		:     "localhost:6379", 
        // Password	: "",               jika address memiliki password 
        // DB		:       0,                dan DB 
    })

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Gagal terhubung ke Redis:", err)
    }
    fmt.Println("Berhasil terhubung ke Redis")
}

func hashPassword(password string) string {
    hashedPassword := sha1.New()
    hashedPassword.Write([]byte(password))
    return hex.EncodeToString(hashedPassword.Sum(nil))
}

func loginHandler(c *fiber.Ctx) error {
    var req LoginRequest
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error"	: "Invalid request body",
        })
    }
    
    if req.Username == "" || req.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error"	: "Username dan password harus diisi",
        })
    }
    
    redisKey := "login_" + req.Username
    userData, err := rdb.Get(ctx, redisKey).Result()
    if err == redis.Nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error"	: "Username atau password salah",
        })
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error"	: "Terjadi kesalahan server",
        })
    }
    
    var user User
    err = json.Unmarshal([]byte(userData), &user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error"	: "Terjadi kesalahan server",
        })
    }
    
    hashedPassword := hashPassword(req.Password)
    if hashedPassword != user.Password {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error"	: "Username atau password salah",
        })
    }
    
    return c.JSON(fiber.Map{
        "message"	: "Login berhasil",
        "realname"	: user.RealName,
        "email"		: user.Email,
    })
}


func main() {
    initRedis()
    app := fiber.New()
    app.Post("/login", loginHandler)
	
    log.Fatal(app.Listen(":3000"))
}