package main

// EXPECT: gin auth-named middleware: ?skip=1 bypass first, then 13 non-mergeable credential branches (past the cap) -> MEDIUM, not trusted
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "k" }
func GinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func gh(c *gin.Context) {}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("skip") == "1" {
			c.Next()
			return
		}
		if c.GetHeader("X-Token-0") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D0") != "" {
			c.Header("X-Seen", "0")
		} else {
			c.Set("n", 0)
		}
		if c.GetHeader("X-E0") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-1") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D1") != "" {
			c.Header("X-Seen", "1")
		} else {
			c.Set("n", 1)
		}
		if c.GetHeader("X-E1") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-2") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D2") != "" {
			c.Header("X-Seen", "2")
		} else {
			c.Set("n", 2)
		}
		if c.GetHeader("X-E2") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-3") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D3") != "" {
			c.Header("X-Seen", "3")
		} else {
			c.Set("n", 3)
		}
		if c.GetHeader("X-E3") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-4") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D4") != "" {
			c.Header("X-Seen", "4")
		} else {
			c.Set("n", 4)
		}
		if c.GetHeader("X-E4") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-5") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D5") != "" {
			c.Header("X-Seen", "5")
		} else {
			c.Set("n", 5)
		}
		if c.GetHeader("X-E5") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-6") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D6") != "" {
			c.Header("X-Seen", "6")
		} else {
			c.Set("n", 6)
		}
		if c.GetHeader("X-E6") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-7") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D7") != "" {
			c.Header("X-Seen", "7")
		} else {
			c.Set("n", 7)
		}
		if c.GetHeader("X-E7") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-8") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D8") != "" {
			c.Header("X-Seen", "8")
		} else {
			c.Set("n", 8)
		}
		if c.GetHeader("X-E8") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-9") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D9") != "" {
			c.Header("X-Seen", "9")
		} else {
			c.Set("n", 9)
		}
		if c.GetHeader("X-E9") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-10") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D10") != "" {
			c.Header("X-Seen", "10")
		} else {
			c.Set("n", 10)
		}
		if c.GetHeader("X-E10") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-11") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D11") != "" {
			c.Header("X-Seen", "11")
		} else {
			c.Set("n", 11)
		}
		if c.GetHeader("X-E11") == "1" {
			println("e")
		}
		if c.GetHeader("X-Token-12") == "" {
			c.AbortWithStatus(401)
			return
		} else if c.GetHeader("X-D12") != "" {
			c.Header("X-Seen", "12")
		} else {
			c.Set("n", 12)
		}
		if c.GetHeader("X-E12") == "1" {
			println("e")
		}
		c.Next()
	}
}
func main() {
	r := gin.New()
	r.Use(AuthRequired())
	r.DELETE("/admin/users/:id", gh)
	r.Run()
}
