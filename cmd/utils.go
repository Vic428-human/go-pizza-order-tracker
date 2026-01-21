package main

import (
	"encoding/json"
	"fmt"
	"html/template" // 這邊不要用成 text/template，會導致 Gin 無法正確渲染模板 (SetHTMLTemplate)
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port   string
	DBPath string
}

func loadConfig() Config {
	return Config{
		Port:   getEnv("PORT", "8080"), // 定義key 跟 value
		DBPath: getEnv("DATABASE_URL", "./data/orders.db"),
	}
}

// Example: Setting an environment variable in code
// err := os.Setenv("NEW_VAR", "GoLang Rocks!")
// Retrieve the newly set variable
// fmt.Println("NEW_VAR =", os.Getenv("NEW_VAR")) => NEW_VAR = GoLang Rocks!
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func loadTemplates(router *gin.Engine) error {
	functions := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		// any = interface{}
		// return template.JS tells engine do not escape this value
		// Create a FuncMap to add custom functions (e.g., JSON encoding)
		"toJSON": func(v interface{}) template.JS { // 參數 v 需要編碼成JSON格式的原始資料，通常是個結構。
			b, err := json.Marshal(v) // Marshal()會回傳JSON字串([]byte切片)以及error值，如果編碼失敗 error 就不為nil。
			if err != nil {
				fmt.Println("Error marshaling:", err)
			}
			// https://ithelp.ithome.com.tw/articles/10335017
			return template.JS(b)
		},
	}

	tmpl, err := template.New("").Funcs(functions).ParseGlob("templates/*.tmpl")
	if err != nil {
		return err
	}
	// Gin 的 SetHTMLTemplate 需要 *html/template.Template
	// 但你現在用的是 Go 標準庫的 text/template
	// Gin 的 router.SetHTMLTemplate 明確要求 *html/template.Template，所以你不能傳 *text/template.Template。
	// 只要把 text/template 換成 html/template
	router.SetHTMLTemplate(tmpl)
	return nil
}

// 設置 session => 工具函式風格 :要回傳 error → 不要在裡面 c.String，交給呼叫者處理。工具函式可以在不同情境下重複使用，不會綁死在某種回應方式上

/*
	// 设置 session 路由
	r.GET("/set-session", SetSession)

// 要自己處理 HTTP → 不要回傳 error，直接在函式裡面回應，適合路由直接呼叫的時候的寫法

	func SetSession(c *gin.Context) {
	    session := sessions.Default(c)
	    session.Set("user", "Fossen")
	    if err := session.Save(); err != nil {
	        c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving session: %s", err.Error()))
	        return
	    }
	    c.String(http.StatusOK, "Session set")
	}
*/
func SetSession(c *gin.Context, key string, value string) error {
	session := sessions.Default(c)
	session.Set(key, value)
	return session.Save()
}

// 獲取 session
func GetSession(c *gin.Context, key string) string {
	session := sessions.Default(c)
	str, ok := session.Get(key).(string)
	if !ok {
		return ""
	}
	return str
}
