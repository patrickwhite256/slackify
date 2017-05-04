package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/patrickwhite256/slackify/config"
	"github.com/patrickwhite256/slackify/external/slack"
	"github.com/patrickwhite256/slackify/external/spotify"
	"github.com/patrickwhite256/slackify/model"
)

func main() {
	var confFilename string
	flag.StringVar(&confFilename, "config", "config.json", "Config Filename")

	flag.Parse()
	conf, err := config.FromFile(confFilename)
	if err != nil {
		log.Fatal(err)
	}

	db, err := connectDB(conf)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: support more than one user :^)
	user, err := model.LoadUser(1, db)
	if err != nil {
		log.Fatal(err)
	}

	slackClient := slack.NewClient(user.SlackAccessToken)
	spotifyClient := spotify.NewClient(conf, user.SpotifyRefreshToken)

	//infinite loop for PoC
	for {
		song, err := spotifyClient.GetCurrentlyPlaying()
		if err != nil {
			log.Fatal(err)
		}
		var emoji, status string
		if song != nil {
			emoji = user.SpotifyEmoji
			status = fmt.Sprintf("%s - %s", song.Name, song.Artists)
		} else {
			emoji = user.PreviousEmoji
			status = user.PreviousStatus
		}
		err = slackClient.SetStatus(status, emoji)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updated - %s %s\n", emoji, status)
		time.Sleep(time.Second * time.Duration(conf.UpdateTime))
	}
}

func connectDB(conf *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", conf.DBUser, conf.DBPassword, conf.DBAddress, conf.DBDatabase))
	if err != nil {
		return nil, err
	}
	// test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}
