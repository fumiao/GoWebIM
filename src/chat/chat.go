package chat

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"text/template"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func Index(c *echo.Context) error {
	uuid := c.Query("uuid")
	if len(uuid) == 0 {
		uuid, _ = GenUUID()
		return c.Redirect(http.StatusSeeOther, "/index?uuid="+uuid)
	}
	fmt.Printf("login username: %s", uuid)
	content := make(map[string]string)
	content["name"] = uuid
	return c.Render(http.StatusOK, "index", content)
}

func Ws(c *echo.Context) (err error) {
	ws := c.Socket()
	msg := ""

	for {
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			return nil
		}
		fmt.Println(msg)
	}

}

func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}
