package main

import (
	"fmt"
	"goweb2/capybara"
	"reflect"

	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	// 初始化数据库
	//	InitDatabase()
	cap := capybara.New()

	cap.GET("/user/:id/post/:post_id", func(ctx capybara.Context) {
		id := ctx.Param("id")           // 获取路径参数 "id"
		post_id := ctx.Param("post_id") // 获取路径参数 "id"
		funcType := reflect.TypeOf(ctx.Handler())
		ctx.JSON(200, map[string]interface{}{
			"id":      id,
			"post_id": post_id,
			"path":    ctx.Path(),
			"handler": funcType.Name()})
	})
	// 路由组
	// authGroup := cap.Group("/auth")
	// authGroup.POST("/login", Login, Logging)
	// authGroup.POST("/register", Register, Logging)

	// profileGroup := cap.Group("/profile")
	// profileGroup.Use(JWTAuth2("capybara"))
	// profileGroup.POST("/viewUser", ViewUserInformation, Logging)

	// cap.GET("/html", HtmlTest)
	cap.Run(":8080")
}

func HtmlTest(ctx capybara.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Complex HTML</title>
</head>
<body>
    <h1>Welcome to My Website</h1>
    <p>This is a complex HTML page.</p>
    <ul>
        <li>Item 1</li>
        <li>Item 2</li>
        <li>Item 3</li>
    </ul>
</body>
</html>
`
	ctx.HTML(200, html)
}

func InitDatabase() {
	dsn := "root:0220059cyCY@tcp(127.0.0.1:3306)/myweb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Println("failed connection db")
	}
	DB = db
	DB.AutoMigrate(&User{})
	fmt.Println("创建表成功")
}

type User struct {
	Name     string `json:"user"`
	Password string `json:"password"`
}

func TestDBCreat() {
	// 写入单条数据
	user := User{Name: "Alice", Password: "0220059cyCY"}
	result := DB.Create(&user)
	if result.Error != nil {
		panic("failed to create user")
	}
}

// 实现登陆
func Login(ctx capybara.Context) {
	currentUser := User{}
	if err := ctx.Bind(&currentUser); err != nil {
		ctx.JSON(500, map[string]string{"message": err.Error()})
	}
	fmt.Println(currentUser.Name, currentUser.Password)
	// 存储从数据库中查询的结果
	saveCurrentUser := User{}
	result := DB.Where("name = ?", currentUser.Name).First(&saveCurrentUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("用户名 '%s' 不存在\n", currentUser.Name)
		} else {
			fmt.Printf("查询出错： %v\n", result.Error)
		}
	} else {
		fmt.Println(saveCurrentUser)
		if saveCurrentUser.Password == currentUser.Password {
			// 表示登陆成功，然后向客户端传输JWT Token
			token, err := generateToken(saveCurrentUser)
			if err != nil {
				ctx.JSON(500, map[string]string{"message": fmt.Sprintf("生成 Token 失败：%v", err)})
				return
			}
			ctx.JSON(200, map[string]string{
				"message": "登陆成功",
				"token":   token})
		}
		return
	}
	ctx.JSON(500, map[string]string{"message": "密码错误"})
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func generateToken(user User) (string, error) {
	claims := Claims{
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token 有效期 24 小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名 Token
	secretKey := []byte("capybara")
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func Register(ctx capybara.Context) {
	currentUser := User{}
	if err := ctx.Bind(&currentUser); err != nil {
		ctx.JSON(500, map[string]string{"message": err.Error()})
	}
	fmt.Println(currentUser.Name, currentUser.Password)

	// 检查姓名是否已经存在
	result := DB.Where("username = ?", currentUser.Name).First(&currentUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("用户名 '%s' 不存在\n", currentUser.Name)
		} else {
			fmt.Printf("查询出错： %v\n", result.Error)
		}
	} else {
		fmt.Printf("用户名 '%s' 已存在\n", currentUser.Name)
		ctx.JSON(500, map[string]string{"message": "用户名已经存在"})
		return
	}

	// 存储到数据库中
	result = DB.Create(&currentUser)
	if result.Error != nil {
		panic("failed to create user")
	}
	ctx.JSON(200, map[string]string{"message:": "用户注册成功"})
}

func ViewUserInformation(ctx capybara.Context) {
	ctx.JSON(200, map[string]string{"message": "ViewUserInformation"})
}

// 日志中间件
func Logging(next capybara.HandlerFunc) capybara.HandlerFunc {
	return func(ctx capybara.Context) {
		ctx.JSON(200, map[string]string{"message": "Logging"})
		log.Printf("Resquest: %s, %s", ctx.Request().URL, ctx.Request().Method)
		next(ctx)
	}
}

func ScanClient() capybara.Middlewares {
	return func(next capybara.HandlerFunc) capybara.HandlerFunc {
		return func(ctx capybara.Context) {
			ctx.JSON(200, map[string]string{"message": "ScanClient"})
			next(ctx)
		}
	}
}

func JWTAuth2(secert string) capybara.Middlewares {
	return func(next capybara.HandlerFunc) capybara.HandlerFunc {
		return func(ctx capybara.Context) {
			authHeader := ctx.GetHeader("Authorization")
			if authHeader == "" {
				ctx.JSON(404, map[string]string{"error": "Authorization header missing"})
				return
			}
			// 提取 Token 字符串
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				ctx.JSON(404, map[string]string{"error": "Invalid token format"})
				return
			}
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secert), nil
			})
			if err != nil || !token.Valid {
				ctx.JSON(404, map[string]string{"error": "Invalid token"})
				return
			}
			// 存储 Claims 到上下文
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx.Set("user", claims)
			} else {
				ctx.JSON(404, map[string]string{"error": "Failed to parse token claims"})
				return
			}
			ctx.JSON(200, map[string]string{"message": "auth middlewares ok"})
			next(ctx)
		}
	}
}
