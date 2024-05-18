package tastybot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Bot struct {
	Id         int
	TastyToken string
	BaseUrl    string // prod or dev stands
	Status     string // waiting | case open | techies | stoped by error
	CasesCount int    //Amount of opened cases
}

type Storage interface {
	Add(ctx context.Context, bot *Bot) (int, error)
	PickById(ctx context.Context, id int) (*Bot, error)
	PickAll(ctx context.Context) ([]Bot, error)
	ChangeStatusById(ctx context.Context, id int, s string) error
	IncreaseCaseCountById(ctx context.Context, id int) error
}

type OpenCaseResBody struct {
	Status bool `json:"status"`
	Data   struct {
		Code string ` json:"code"`
	}
}

func New(tastyToken, baseUrl string, stg Storage) error {
	newBot := &Bot{TastyToken: tastyToken, Status: "READY", BaseUrl: baseUrl}
	id, err := stg.Add(context.Background(), newBot)
	if err != nil {
		return fmt.Errorf("error during adding bot %w", err)
	}
	fmt.Printf("Created bot with id %v \n", id)
	return nil
}

func FindBotById(id int, stg Storage) (*Bot, error) {
	bot, err := stg.PickById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func (bot *Bot) openCase(caseName string) string {
	url := bot.BaseUrl + "/api/v1/cases/" + caseName + "/open"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+bot.TastyToken)
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
			return "STOPPED BY: authentication_exception"
		case "not_enough_balance":
			return "STOPPED BY: not_enough_balance"
		case "unpredicted_exception":
			return "STOPPED BY: unpredicted_exception"
		}
	}

	fmt.Printf("bot: %v openned case %v at %v \n", bot.Id, caseName, time.Now().Format(time.ANSIC))
	return "ok"
}

func RunCases(id int, caseName string, stg Storage) error {
	if id == 0 {
		return fmt.Errorf("please provide bot id")
	}
	status := "WORKING CASE OPPENING"

	err := stg.ChangeStatusById(context.Background(), id, status)
	if err != nil {
		return err
	}
	bot, err := stg.PickById(context.Background(), id)
	if err != nil {
		return err
	}

	//case open cooldown
	timer, _ := time.ParseDuration("10s")

	var openStatus string
	for {
		if bot.Status == "WORKING CASE OPPENING" {
			openStatus = bot.openCase(caseName)

			if openStatus == "ok" {
				stg.IncreaseCaseCountById(context.Background(), bot.Id)
			} else {
				return fmt.Errorf("bot status: %v", openStatus)
				break
			}
			time.Sleep(timer)
		}
	}

	return fmt.Errorf("bot status: %v", openStatus)
}

func (u *Bot) RunTechies() {
	u.Status = "RUNNING"
	techiesCooldown := 3*time.Hour + 5*time.Minute
	for {
		u.PlayTechies()
		time.Sleep(techiesCooldown)
	}

}

func (b *Bot) GetStatus() {
	fmt.Printf("bot: <id: %v, baseURl: %v, status: %v> \n", b.Id, b.BaseUrl, b.Status)
}

func StatusAll(stg Storage) ([]Bot, error) {
	bots, err := stg.PickAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error during get status All %w", err)
	}

	for _, bot := range bots {
		bot.GetStatus()
	}

	return bots, nil

}

func (u *Bot) PlayTechies() {
	ttCookie := &proto.NetworkCookie{
		Name:     "tastyToken",
		Value:    u.TastyToken,
		Domain:   "tastydrop.in",
		Path:     "/",
		Expires:  0,
		HTTPOnly: false,
		Secure:   true,
		Priority: proto.NetworkCookiePriorityMedium,

		SameParty: false,

		SourceScheme: proto.NetworkCookieSourceSchemeSecure,
	}

	b := rod.New()
	page := b.MustConnect().
		MustSetCookies(ttCookie).
		MustPage("https://tastydrop.in/techies")

	if page == nil {
		fmt.Println("Ошибка: Не удалось подключиться к странице")
		return
	}
	timer := time.NewTimer(60 * time.Second)
	isGame := true

	go func() {
		<-timer.C
		fmt.Println("Время вышло, выход из программы")
		isGame = false
	}()

	page.MustWaitDOMStable()
	startGame(page)
	for i := 1; i <= 6; i++ {
		if !isGame {
			fmt.Println("Игра окончена на ходу: ", i)
			break
		}
		isGame = makeTurn(page, i)
		fmt.Println("Ход", i)

	}
}

func startGame(page *rod.Page) {
	page.MustWaitDOMStable()
	el := page.MustElement("div.button-block__default")
	startBtn := el.MustElement("a")
	startBtn.MustClick()
	time.Sleep(2 * time.Second)
}

func makeTurn(page *rod.Page, number int) bool {
	page.MustWaitDOMStable()
	gameField := page.MustElement("div.game-map-controller__content")
	colls := gameField.MustElements("div.game-map-controller__content-column")
	if len(colls) == 0 {
		return false
	}
	coll := colls[number]
	cells := coll.MustElements("div.game-field")
	if len(cells) == 0 {
		return false
	}
	cell := cells[1]
	cell.MustClick()
	return true
}
