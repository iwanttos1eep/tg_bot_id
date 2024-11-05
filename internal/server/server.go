package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
)

var (
	currUserID   string
	counter      int
	counterMutex sync.Mutex
	commandChan  = make(chan string)
	userIDSent   = make(map[int64]bool)
)

func StartWebServer(botAPI *tgbotapi.BotAPI) {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.GET("/get-user-id", func(c echo.Context) error {
		return getUserIDHandler(c, botAPI)
	})
	e.GET("/increment-counter", incrementCounterHandler)
	e.GET("/command-work", commandWorkHandler)
	e.GET("/set-user-id", setUserIDHandler)

	e.Logger.Fatal(e.Start(":1323"))

	log.Printf("Web server started on :1323")
	log.Fatal(http.ListenAndServe(":1323", nil))
}

func StartBot(botAPI *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID

		if !userIDSent[userID] {
			go sendUserIDToServer(userID)
			userIDSent[userID] = true
		}

		if update.Message.IsCommand() {
			HandleCommand(update, botAPI)
		}
	}
}

func HandleCommand(update tgbotapi.Update, botAPI *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	userID := update.Message.From.ID

	switch update.Message.Command() {
	case "start":
		msg.Text = "Бот запущен"
	case "check_id":
		msg.Text = fmt.Sprintf("Ваш ID: %v", userID)
		log.Printf("UserID message: %v\n\n", msg)
	case "work":
		currentTime := time.Now().Format("15:04:05")
		commandChan <- fmt.Sprintf("Введена команда в %s", currentTime)
		log.Printf("This is commandChan: %v\n\n", commandChan)
	default:
		msg.Text = "Неизвестная команда"
	}

	if _, err := botAPI.Send(msg); err != nil {
		log.Printf("[ERROR] failed to send message: %v", err)
	}
}

func sendUserIDToServer(userID int64) {
	url := fmt.Sprintf("http://localhost:1323/set-user-id?userID=%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] failed to send user ID: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Println("User ID sent to server:", userID)
}

func setUserIDHandler(c echo.Context) error {
	userID := c.QueryParam("userID")
	currUserID = userID
	log.Printf("User ID received and stored: %v\n", userID)
	return c.String(http.StatusOK, "User ID set")
}

func getUserIDHandler(c echo.Context, botAPI *tgbotapi.BotAPI) error {

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	userIDInt, err := strconv.ParseInt(currUserID, 10, 64)
	if err != nil {
		log.Printf("[ERROR] failed to parse user ID: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to parse user ID")
	}

	msg := tgbotapi.NewMessage(userIDInt, fmt.Sprintf("Ваш ID: %v", currUserID))
	if _, err := botAPI.Send(msg); err != nil {
		log.Printf("[ERROR] failed to send message: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to send message")
	}

	response := map[string]interface{}{
		"userId": currUserID,
	}
	return c.JSON(http.StatusOK, response)
}

func incrementCounterHandler(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	counterMutex.Lock()
	counter++
	count := counter
	counterMutex.Unlock()

	response := map[string]interface{}{
		"count": count,
	}

	return c.JSON(http.StatusOK, response)
}

func commandWorkHandler(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Error 500")
	}

	go func() {
		for {
			select {
			case status := <-commandChan:
				fmt.Fprintf(c.Response().Writer, "data: %s\n\n", status)
				if flusher != nil {
					flusher.Flush()
				}
			case <-c.Request().Context().Done():
				log.Println("Client disconnected")
				return
			}
		}
	}()

	<-c.Request().Context().Done()
	return nil
}
