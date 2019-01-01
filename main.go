package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"flag"
	"github.com/Mannaka/learning-go/trace"
	"os"
)

// temp1は一つのテンプレートをさす
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
	trace    trace.Tracer
}

// ServeHTTPはHTTPリクエストを処理
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	t.temp1.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈
	r := newRoom()
	// if you want traceOff, you must comment this line out. 
	r.tracer = trace.New(os.Stdout)
    http.Handle("/", &templateHandler{filename: "chat.html"})
    http.Handle("/room", r)
    // チャットルームを開始します
    go r.run()
	log.Println("Webサーバーを開始します。ポート：", *addr)
	// webサーバーを起動します
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe", err)
    }
}