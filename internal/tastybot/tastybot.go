package tastybot

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	Name       string
	Status     string // READY | RUNNING | STOPPED BY
	TastyToken string
	CasesCount int64
}

type OpenCaseResBody struct {
	Status bool `json:"status"`
	Data   struct {
		Code string ` json:"code"`
	}
}

type Storage interface {
	Add(bot *User)
	GetByField(field string) (*User, bool)
	GetAll() []*User
}

func New(tastyToken, name string, stg Storage) *User {
	newUser := &User{TastyToken: tastyToken, Status: "READY", Name: name}
	stg.Add(newUser)
	return newUser
}

func (u *User) openCase(baseURL, caseName string) {
	url := baseURL + "/api/v1/cases/" + caseName + "/open"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+u.TastyToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var resBody OpenCaseResBody
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		fmt.Println(err)
	}
	if !resBody.Status {
		switch resBody.Data.Code {
		case "authentication_exception":
			u.Status = "STOPPED BY: authentication_exception"
		case "not_enough_balance":
			u.Status = "STOPPED BY: not_enough_balance"
		case "unpredicted_exception":
			u.Status = "STOPPED BY: unpredicted_exception"
		}
		return
	}

	u.CasesCount++
}

func (u *User) Run() {
	u.Status = "RUNNING"
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	baseURL := os.Getenv("BASE_URL")
	openCase := os.Getenv("CASE")
	timer := os.Getenv("TIMER")

	if baseURL == "" {
		baseURL = "https://dota2.radiant.dev.tastyteam.cc"
		fmt.Println("You doest set BASE_URL env variable. Using default value: " + baseURL)
	}

	if openCase == "" {
		openCase = "spring24_warm-inv"
		fmt.Println("You doest set CASE env variable. Using default value: " + openCase)
	}

	if timer == "" {
		timer = "10s"
		fmt.Println("You doest set TIMER env variable. Using default value: " + timer)
	}
	timerDuration, err := time.ParseDuration(timer)
	if err != nil {
		fmt.Println(err)
	}

	for {
		if u.Status == "RUNNING" {
			u.openCase(baseURL, openCase)
			time.Sleep(timerDuration)
		}
	}
}

func (u *User) GetStatus() {
	fmt.Println("name: " + u.Name + " status: " + u.Status + " opened cases: " + strconv.FormatInt(u.CasesCount, 10))
}

func StatusAll(stg Storage) {
	bots := stg.GetAll()
	for _, bot := range bots {
		bot.GetStatus()
	}
}
