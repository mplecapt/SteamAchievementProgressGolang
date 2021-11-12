package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const API_KEY string = "1D533C7025222E1B18145EB064E64FB9"

/**
 * Player json structs
 */
type achieveNodePlayer struct {
	ApiName    string `json:"apiname"`
	Achieved   uint8  `json:"achieved"`
	Unlocktime int64  `json:"unlocktime"`
}

func (a achieveNodePlayer) IsAchieved() bool {
	return a.Achieved > 0
}

type playerStats struct {
	Achievements []achieveNodePlayer `json:"achievements"`
	GameName     string              `json:"gameName"`
	SteamID      string              `json:"steamID"`
	Success      bool                `json:"success"`
}

type player struct {
	Playerstats playerStats `json:"playerstats"`
}

/**
 * Game Schema json structs
 */
type achieveNodeGame struct {
	DefaultValue int    `json:"defaultvalue"`
	Description  string `json:"description"`
	DisplayName  string `json:"displayName"`
	Hidden       uint8  `json:"hidden"`
	Icon         string `json:"icon"`
	IconGray     string `json:"icongray"`
	Name         string `json:"name"`
}

type statsNodeGame struct {
	DefaultValue int    `json:"defaultvalue"`
	DisplayName  string `json:"displayName"`
	Name         string `json:"name"`
}

type gameStats struct {
	Achievements []achieveNodeGame `json:"achievements"`
	Stats        []statsNodeGame   `json:"stats"`
}

type gameSchema struct {
	AvailableGameStats gameStats `json:"availableGameStats"`
	GameName           string    `json:"gameName"`
	GameVersion        string    `json:"gameVersion"`
}

type game struct {
	Game gameSchema `json:"game"`
}

/**
 * Final output data
 */
type achieveNodeFinal struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	Icon        string `json:"icon"`
	IconGray    string `json:"icongray"`
	Achieved    bool   `json:"achieved"`
	Unlocktime  int64  `json:"unlocktime"`
}

type dataFinal struct {
	Achievements []achieveNodeFinal `json:"achievements"`
	SteamID      string             `json"steamid"`
	GameName     string             `json"gamename"`
}

/**
 * Steam API calls, storing json into structs
 */
func callSteam(function string, appid string, steamid string, output interface{}) error {
	data, err := http.Get("https://api.steampowered.com/ISteamUserStats/" + function + "/?key=" + API_KEY + "&appid=" + appid + "&steamid=" + steamid)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}
	return nil
}

func GetPlayerProgress(appid, steamid string) ([]byte, error) {
	var p1 player
	err := callSteam("GetPlayerAchievements/v1", appid, steamid, &p1)
	if err != nil {
		return nil, err
	}

	var g game
	err = callSteam("GetSchemaForGame/v2", appid, "", &g)
	if err != nil {
		return nil, err
	}

	data := dataFinal{
		GameName:     p1.Playerstats.GameName,
		SteamID:      p1.Playerstats.SteamID,
		Achievements: make([]achieveNodeFinal, len(g.Game.AvailableGameStats.Achievements)),
	}
	for i := range g.Game.AvailableGameStats.Achievements {
		if g.Game.AvailableGameStats.Achievements[i].Name != p1.Playerstats.Achievements[i].ApiName {
			return nil, err
		}
		data.Achievements[i] = achieveNodeFinal{
			Name:        g.Game.AvailableGameStats.Achievements[i].Name,
			DisplayName: g.Game.AvailableGameStats.Achievements[i].DisplayName,
			Description: g.Game.AvailableGameStats.Achievements[i].Description,
			Hidden:      g.Game.AvailableGameStats.Achievements[i].Hidden > 0,
			Icon:        g.Game.AvailableGameStats.Achievements[i].Icon,
			IconGray:    g.Game.AvailableGameStats.Achievements[i].IconGray,
			Achieved:    p1.Playerstats.Achievements[i].Achieved > 0,
			Unlocktime:  p1.Playerstats.Achievements[i].Unlocktime,
		}
	}

	newbody, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	return newbody, nil
}
