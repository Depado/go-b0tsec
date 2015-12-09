package seen

import (
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/plugins/afk"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "seen"
)

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct {
	Started bool
}

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "Displays if someone is afk or the time since their last message.")
	ib.Privmsg(from, "Example : !seen nickname")
}

// Get is the actual call to your plugin.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
		return
	}
	if len(args) > 0 {
		if to == configuration.Config.BotName {
			to = from
		}
		for _, v := range args {
			if d, ok := afk.Map[v]; ok {
				if d.Reason != "" {
					ib.Privmsgf(to, "%v has been afk for %v : %v", v, time.Since(d.Since).String(), d.Reason)
				} else {
					ib.Privmsgf(to, "%v has been afk for %v", v, time.Since(d.Since).String())
				}
			} else {
				if d, ok := Map[v]; ok {
					ib.Privmsgf(to, "Last message from %v : %v ago.", v, time.Since(d).String())
				}
			}
		}
	}
}

// Start starts the plugin and returns any occured error, nil otherwise
func (p *Plugin) Start() error {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		p.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occured error, nil otherwise
func (p *Plugin) Stop() error {
	p.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (p *Plugin) IsStarted() bool {
	return p.Started
}
