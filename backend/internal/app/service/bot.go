package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/sacurio/jb-challenge/internal/app/model"
	"github.com/sacurio/jb-challenge/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	// TODO: The list of known commands must be from a better source, probably from database as a list.
	knownCmds = []string{"/stock="}
)

type Bot interface {
	IsValidCommand(string) bool
	GetBotUser(User) (*model.User, error)
	ProcessCommandAsync(string, chan<- string)
}

type bot struct {
	command     string
	name        string
	wildcard    string
	url         string
	channelResp chan string
	logger      *logrus.Logger
}

func NewBot(cfg config.BotConfig, logger *logrus.Logger) Bot {
	return &bot{
		name:     cfg.Name,
		wildcard: cfg.Wildcard,
		url:      cfg.URL,
		logger:   logger,
	}
}

func (b *bot) makeRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.logger.Errorf("Error at read the content: %s", err.Error())
		return "", err
	}

	stockService := NewStock(string(body))
	_, err = stockService.Parse()
	if err != nil {
		return "", err
	}

	return stockService.String(), nil
}

func (b *bot) IsValidCommand(input string) bool {
	re := regexp.MustCompile(`^/stock=.+$`)
	return re.MatchString(input)
}

func (b *bot) GetBotUser(userService User) (*model.User, error) {
	botUser, err := userService.GetByUsername(b.name)
	if err != nil {
		return nil, err
	}
	return botUser, nil
}

func (b *bot) ProcessCommandAsync(cmd string, ch chan<- string) {
	url := b.buildURL(cmd)
	req, err := b.makeRequest(url)
	if err != nil {
		b.logger.Errorf("error requesting data from service, %s", err.Error())
	}
	ch <- req
}

func (b *bot) buildURL(cmd string) string {
	var cmdValue string
	for i := 0; i < len(knownCmds); i++ {
		if strings.HasPrefix(cmd, knownCmds[i]) {
			cmdValue = strings.Replace(cmd, knownCmds[i], "", len(knownCmds[i]))
			break
		}
	}
	b.command = cmdValue
	return fmt.Sprintf("%s", strings.Replace(b.url, b.wildcard, cmdValue, len(b.wildcard)))
}
