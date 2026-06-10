package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "caijian"

// 定义 JWT 的签名密钥，实际使用中应从环境变量或配置文件获取
var jwtKey = []byte(secretKey)

// Claims 定义 JWT 的声明信息
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Login 处理登录请求，验证用户信息并颁发 JWT
func Login(c *gin.Context) {
	// 模拟验证用户信息，实际使用中应从数据库或其他数据源验证
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 这里简单模拟验证，实际需要替换为真实的验证逻辑
	if username == "admin" && password == "Edu@9527" {
		// 设置 JWT 的过期时间
		expirationTime := time.Now().Add(24 * time.Hour)
		// 创建 JWT 的声明信息
		claims := &Claims{
			UserID: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		// 创建 JWT 对象
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// 生成签名字符串
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 JWT 失败"})
			return
		}

		// 返回 JWT 给客户端
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
	}
}

// JWTAuthMiddleware JWT 验证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 JWT
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 JWT"})
			c.Abort()
			return
		}

		// 解析 JWT
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 JWT"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT不正确"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
