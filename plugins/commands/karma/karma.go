package karma

import (
	"encoding/json"
	"log"
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/thoj/go-ircevent"
)

const bucketName = "karma"
const mainKey = "main"

// Data is the struct that contains the data about the karma intented to be stored somewhere.
type Data struct {
	Karma map[string]int
}

// Encode encodes a chain to json.
func (d Data) Encode() ([]byte, error) {
	enc, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to Chain
func (d *Data) Decode(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	return nil
}

// Save saves the Data
func (d Data) Save() error {
	return database.BotStorage.Save(bucketName, mainKey, &d)
}

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct {
	Data
	Action map[string]time.Time
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    Allows to add/remove/see karma points to/from a person.")
	ib.Privmsg(from, "    Add/Remove : !karma [+|-] <nickname>")
	ib.Privmsg(from, "    See : !karma = <nickname>")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "+", "-":
			if len(args) > 1 {
				if from != args[1] {
					if val, ok := p.Action[from]; ok {
						if time.Since(val) < 1*time.Minute {
							ib.Notice(from, "Please wait 1 minute between each karma operation.")
							return
						}
					}
					p.Action[from] = time.Now()
					c := 0
					if val, ok := p.Karma[args[1]]; ok {
						c = val
					}
					if args[0] == "+" {
						p.Karma[args[1]] = c + 1
						ib.Privmsgf(configuration.Config.Channel, "Someone gave a karma point to %v, total %v", args[1], c+1)
					} else {
						p.Karma[args[1]] = c - 1
						ib.Privmsgf(configuration.Config.Channel, "Someone took a karma point from %v, total %v", args[1], c-1)
					}
					if err := p.Data.Save(); err != nil {
						log.Println(err)
					}
				} else {
					ib.Notice(from, "Can't add or remove points to yourself.")
					return
				}
			} else {
				ib.Notice(from, "You need to give a nickname to operate on.")
			}
		case "=":
			if len(args) > 1 {
				for _, n := range args[1:] {
					if val, ok := p.Karma[n]; ok {
						ib.Privmsgf(to, "%v has %v point(s).", n, val)
					} else {
						ib.Privmsgf(to, "I don't have records on %v.", n)
					}
				}
			} else {
				ib.Notice(from, "Need at least a nickname.")
			}
		}
	}
}

// New initializes new plugin
func New() Plugin {
	d := Data{make(map[string]int)}
	if err := database.BotStorage.CreateBucket(bucketName); err != nil {
		log.Fatalf("While initializing Karma plugin : %s", err)
	} else {
		log.Println("Created bucket", bucketName)
	}
	database.BotStorage.Get(bucketName, mainKey, &d)
	return Plugin{d, make(map[string]time.Time)}
}