package CLI

//
//import (
//	_ "tasty-bots/internal/storage"
//)
//
//const (
//	errorMessage = "Unknown command. Type 'tb new --token=<tasty token>' to create a new bot"
//)

//func Run() {
//	reader := bufio.NewReader(os.Stdin)
//	botStorage := storage.New()
//
//	fmt.Println("welcome to tasty-bots CLI!")
//	fmt.Println("type 'tb new --name=<bot name> --token=<tasty token>' to create a new case open bot")
//	fmt.Println("type 'tb status --name=<bot name>' to check status of the bot")
//	fmt.Println("use --name=all' to check status of all bots")
//
//	for {
//		fmt.Print("> ")
//		input, _ := reader.ReadString('\n')
//		input = strings.TrimSpace(input)
//
//		parts := make([]string, 0)
//		if strings.Contains(input, " ") {
//			parts = strings.Fields(input)
//		} else {
//			parts = append(parts, "error", "error")
//		}
//		command := parts[1]
//
//		switch command {
//		case "new":
//			go cmdNew(input, botStorage)
//		case "status":
//			go cmdStatus(input, botStorage)
//		default:
//			fmt.Println(errorMessage)
//		}
//	}
//}

//func cmdNew(input string, storage tastybot.Storage) {
//	token := getFlagValue("token", input)
//	name := getFlagValue("name", input)
//
//	if token == "" || name == "" {
//		fmt.Println("to create a new bot please specify a token with --token=<tasty token> and a name with --name=<bot name>")
//		return
//	}
//
//	fmt.Println("Creating a new bot with token: " + token)
//	bot := tastybot.New(token, name, storage)
//	//bot.Run()
//	bot.RunTechies()
//}

//func cmdStatus(input string, storage tastybot.Storage) {
//	name := getFlagValue("name", input)
//	switch name {
//	case "":
//		fmt.Println("to check status of the bot please specify a token with --name=<bot name>")
//	case "all":
//		tastybot.StatusAll(storage)
//	default:
//		bot, ok := storage.GetByField(name)
//		if !ok {
//			fmt.Println("no bot with token " + name)
//			return
//		}
//		bot.GetStatus()
//	}
//
//}

//func getFlagValue(flagName string, input string) string {
//	key := "--" + flagName
//	re := regexp.MustCompile(key + `=([^ ]+)`)
//	match := re.FindStringSubmatch(input)
//	if len(match) > 1 {
//		return match[1]
//	}
//	return ""
//}
