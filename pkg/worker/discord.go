package worker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	dfpb "github.com/huo-ju/dfserver/pkg/pb"
	"github.com/huo-ju/dfserver/pkg/service"
)

type ProcessDiscordWorker struct {
	ds *service.DiscordService
}

func (f *ProcessDiscordWorker) Name() string {
	return "process.discord"
}

func (f *ProcessDiscordWorker) Work(lastoutput *dfpb.Output, settingsdata []byte) (bool, error) {
	var settings map[string]interface{}
	err := json.Unmarshal(settingsdata, &settings)
	if err != nil {
		//TODO: save err log
		return true, err
	}
	r := bytes.NewReader(lastoutput.Data)
	messageid := settings["message_id"].(string)
	channelid := settings["channel_id"].(string)
	guildid := settings["guild_id"].(string)

	content := fmt.Sprintf("%s by %s", string(lastoutput.Args), *lastoutput.ProducerName)
	ref := &discordgo.MessageReference{MessageID: messageid, ChannelID: channelid, GuildID: guildid}
	msg := &discordgo.MessageSend{
		Content:   content,
		File:      &discordgo.File{Name: "output.png", Reader: r},
		Reference: ref,
	}

	f.ds.ReplyMessage(channelid, msg)

	return true, err
}