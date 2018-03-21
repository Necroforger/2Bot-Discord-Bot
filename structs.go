package main

import (
	"sync"
	"time"

	"github.com/jonas747/dca"

	"github.com/bwmarrin/discordgo"
)

type ibStruct struct {
	Path   string `json:"path"`
	Server string `json:"server"`
}

type rule34 struct {
	PostCount int `xml:"count,attr"`

	Posts []struct {
		URL string `xml:"file_url,attr"`
	} `xml:"post"`
}

type config struct {
	Game   string `json:"game"`
	Prefix string `json:"prefix"`
	Token  string `json:"token"`
	Port   string `json:"port"`

	InDev bool `json:"indev"`

	DiscordPWKey string `json:"discord.pw_key"`

	CurrImg int `json:"curr_img_id"`
	MaxProc int `json:"maxproc"`

	Blacklist []string `json:"blacklist"`
}

type servers struct {
	Count      int
	VoiceInsts int

	Mutex sync.RWMutex `json:"-"`

	Server map[string]*server
}

type server struct {
	LogChannel string `json:"log_channel"`
	Prefix     string `json:"server_prefix"`

	Log    bool `json:"log_active"`
	Kicked bool `json:"kicked"`
	Nsfw   bool `json:"nsfw"`

	//Enabled, Message, Channel
	JoinMessage [3]string `json:"join"`

	VoiceInst *voiceInst `json:"-"`

	Playlists map[string][]song `json:"playlists"`
}

type voiceInst struct {
	ChannelID string

	Queue []song

	Playing bool

	Done chan error

	*sync.Mutex

	StreamingSession *dca.StreamingSession

	VoiceCon *discordgo.VoiceConnection
}

type song struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Image string `json:"image"`

	Duration time.Duration `json:"duration"`
}

type imageQueue struct {
	QueuedMsgs map[string]*queuedImage
}

type queuedImage struct {
	ReviewMsgID   string `json:"reviewMsgID"`
	AuthorID      string `json:"author_id"`
	AuthorDiscrim string `json:"author_discrim"`
	AuthorName    string `json:"author_name"`
	ImageName     string `json:"image_name"`
	ImageURL      string `json:"image_url"`

	FileSize int `json:"file_size"`
}

type command struct {
	Name string
	Help string

	NoahOnly      bool
	RequiresPerms bool

	PermsRequired int

	Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

type users struct {
	User map[string]*user
}

type user struct {
	Images map[string]string `json:"images"`

	DiskQuota    int `json:"quota"`
	CurrDiskUsed int `json:"curr_used"`
	QueueSize    int `json:"queue_size"`

	TempImages []string `json:"temp_images"`
}
