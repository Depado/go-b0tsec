package duckduckgo

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const apiURL = "http://api.duckduckgo.com/?q=%s&format=json%s"

type message struct {
	Definition       string
	DefinitionSource string
	Heading          string
	AbstractText     string
	Abstract         string
	AbstractSource   string
	Image            string
	Type             string
	AnswerType       string
	Redirect         string
	DefinitionURL    string
	Answer           string
	AbstractURL      string
	Results          []relatedTopic
	RelatedTopics    []relatedTopic
}

type relatedTopic struct {
	Result string
	Icon   struct {
		URL    string
		Height interface{}
		Width  interface{}
	}
	FirstURL string
	Text     string
}

// Plugin is the duckduckgo plugin.
type Plugin struct{}

// Help provides some help on the plugin
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    Search directly on DuckDuckGo.")
	ib.Privmsg(from, "    Example : !command Who is James Cameron ?")
}

// Get actually sends the data to the channel
func (p Plugin) Get(ircbot *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		res, err := p.fetch(strings.Join(args, " "))
		if err != nil || res == "" {
			return
		}
		ircbot.Privmsg(configuration.Config.Channel, res)
	}
}

func (p Plugin) fetch(query string) (string, error) {
	var t message
	url := utils.EncodeURL(apiURL, query)
	err := utils.FetchURL(url, &t)
	if err != nil {
		return "", err
	}
	return t.Abstract, nil
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}
