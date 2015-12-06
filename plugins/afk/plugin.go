package afk

import (
	"strings"
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/pluginsinit"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "choice"
)

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct{}

func init() {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		pluginsinit.Plugins[pluginCommand] = new(Plugin)
	}
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Tell the world you're afk, for a reason. Or not.")
	ib.Privmsg(from, "Example : !afk reason.")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	reason := ""
	if len(args) > 0 {
		reason = strings.Join(args, " ")
	}
	Map[from] = Data{time.Now(), reason}
	if reason != "" {
		ib.Privmsgf(configuration.Config.Channel, "%v is afk : %v", from, reason)
	} else {
		ib.Privmsgf(configuration.Config.Channel, "%v is afk.", from)
	}
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}
