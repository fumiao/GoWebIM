package chat

import (
	"container/list"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jeffail/gabs"
	"github.com/labstack/echo"
	"github.com/yiqguo/GoWebIM/src/models"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"text/template"
)

func init() {
	go chatmsg()
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn
}

var (
	sub      = make(chan Subscriber, 10)
	subevent = make(chan models.Event, 10)
	sublist  = list.New()
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
	content := make(map[string]interface{})
	content["name"] = uuid
	content["users"] = getUsers(sublist, uuid)
	return c.Render(http.StatusOK, "index", content)
}

func Ws(c *echo.Context) (err error) {
	uuid := c.Query("uuid")
	ws := c.Socket()
	join(uuid, ws)
	msg := ""

	for {
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			event := models.NewEvent(models.EVENT_BROAD, uuid, "系统消息: "+uuid+"离开聊天室")
			removeUser(sublist, uuid)
			sendMessage(event)
			return nil
		}
		fmt.Println(msg)
		msgjson, err := gabs.ParseJSON([]byte(msg))
		if err != nil {
			return nil
		}
		st := msgjson.Path("st").Data().(string)
		msg := msgjson.Path("msg").Data().(string)
		if st == "all" {
			event := models.NewEvent(models.EVENT_BROAD, uuid, msg)
			sendMessage(event)
		} else {
			event := models.NewEvent(models.EVENT_UNICAST, st, msg)
			sendMessage(event)
		}
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

func join(username string, conn *websocket.Conn) {
	sub <- Subscriber{Name: username, Conn: conn}
}

func chatmsg() {
	for {
		select {
		case s := <-sub:
			if isNewUser(sublist, s.Name) {
				event := models.NewEvent(models.EVENT_BROAD, s.Name, "系统消息:"+s.Name+" 加入聊天室")
				subevent <- event
				fmt.Printf("%s is a new user\n", s.Name)
				sublist.PushBack(s)
			} else {
				fmt.Printf("%s is a old user\n", s.Name)
			}
		case se := <-subevent:
			sendMessage(se)
		}
	}
}

func isNewUser(subscribers *list.List, user string) bool {
	for subl := subscribers.Front(); subl != nil; subl = subl.Next() {
		if subl.Value.(Subscriber).Name == user {
			return false
		}
	}
	return true
}

func getUsers(subscribers *list.List, name string) []string {
	var users []string
	for subl := subscribers.Front(); subl != nil; subl = subl.Next() {
		if sname := subl.Value.(Subscriber).Name; sname != name {
			users = append(users, sname)
		}
	}
	return users
}

func getUser(subscribers *list.List, name string) *websocket.Conn {
	for subl := subscribers.Front(); subl != nil; subl = subl.Next() {
		if sname := subl.Value.(Subscriber).Name; sname == name {
			return subl.Value.(Subscriber).Conn
		}
	}
	return nil
}

func removeUser(subscribers *list.List, name string) {
	for subl := subscribers.Front(); subl != nil; subl = subl.Next() {
		if sname := subl.Value.(Subscriber).Name; sname == name {
			subscribers.Remove(subl)
		}
	}
}

func getConns(subscribers *list.List, name string) []*websocket.Conn {
	var conns []*websocket.Conn
	for subl := subscribers.Front(); subl != nil; subl = subl.Next() {
		if sname := subl.Value.(Subscriber).Name; sname != name {
			conns = append(conns, subl.Value.(Subscriber).Conn)
		}
	}
	return conns
}

func sendMessage(event models.Event) {
	if event.SendType == models.EVENT_BROAD {
		conns := getConns(sublist, event.Name)
		for _, conn := range conns {
			_ = websocket.Message.Send(conn, event.Content)
		}
	} else {
		conn := getUser(sublist, event.Name)
		_ = websocket.Message.Send(conn, event.Content)
	}
}
