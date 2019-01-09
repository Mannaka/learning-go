package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/Mannaka/learning-go/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// temp1は一つのテンプレートをさす
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

// ServeHTTPはHTTPリクエストを処理
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.temp1.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈
	// Gomniauthのセットアップ
	// TODO: セキュリティキー、クライアントID、秘密鍵を環境変数から読み込ませる
	gomniauth.SetSecurityKey(os.Getenv("SEQ_KEY_FOR_CHAT"))
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/github"),
		google.New(os.Getenv("CLIENT_FOR_CHAT"), os.Getenv("KEY_FOR_CHAT"), "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom(UseGravatar)
	// if you want traceOff, you must comment this line out.
	r.tracer = trace.New(os.Stdout)
	// TODO:　localのbootstrapを適用させる
	// http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()
	log.Println("Webサーバーを開始します。ポート：", *addr)
	// webサーバーを起動します
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
