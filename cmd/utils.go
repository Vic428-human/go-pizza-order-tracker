package main

/** 用途:
1. 載入環境變數config
2. 載入模板
3. 建立基於 GORM 的 session store。
4. 對 session 的操作
**/
import (
	"encoding/json"
	"fmt"
	"html/template" // 這邊不要用成 text/template，會導致 Gin 無法正確渲染模板 (SetHTMLTemplate)
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	Port   string
	DBPath string
}

// 1. 載入環境變數config
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

// 2. 載入模板
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

// 3. 建立基於 GORM 的 session store。
func createSessionStore(db *gorm.DB, secret []byte) gormsessions.Store {
	store := gormsessions.NewStore(db, true, secret)

	// 用 Options() 方法設定 session 選項
	store.Options(sessions.Options{
		Path:     "/",   // Cookie 的作用路徑，"/" 表示整個網站都會帶上這個 Cookie
		MaxAge:   86400, // Cookie 的存活時間（秒），這裡是 86400 秒 = 1 天
		Secure:   true,  // 是否只在 HTTPS 連線中傳送 Cookie，true 表示僅限安全連線
		HttpOnly: true,  // 是否禁止 JavaScript 存取 Cookie，true 可防止 XSS 攻擊
		SameSite: 3,     // SameSite 屬性，用來限制跨站請求攜帶 Cookie
		// 0 = Default, 1 = Lax, 2 = Strict, 3 = None
		// 這裡設為 3，表示允許跨站請求攜帶 Cookie（需搭配 Secure）
	})

	return store
}

/*
情況1.设置 session 路由，要自己處理 HTTP → 不要回傳 error，直接在函式裡面回應，適合路由直接呼叫的時候的寫法

router.GET("/set-session", SetSession)

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

// 情況2.工具函式風格:要回傳 error → 不要在裡面 c.String，交給呼叫者處理。工具函式可以在不同情境下重複使用，不會綁死在某種回應方式上
func SetSession(c *gin.Context, key string, value interface{}) error {
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

func ClearAllSession(c *gin.Context) error {
	session := sessions.Default(c)
	// session.Delete("user") => 只刪除 session 中 特定 key（這裡是 "user"）的值。
	// session.Clear() =>  清空整個 session 的所有 key-value，讓 session 變成空的 map。
	session.Clear()
	return session.Save()
}
