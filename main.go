package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
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
	//Tags        []string `json:"tags"`
}

type dataFinal struct {
	Achievements []achieveNodeFinal         `json:"achievements"`
	SteamID      string                     `json:"steamid"`
	GameName     string                     `json:"gamename"`
	Heists       map[string]map[string]bool `json:"heists"`
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

/*func getTags(desc string) []string {
	var tags []string

	for _, exp := range []string{
		`(?<=the )(.*?)(?=\s* job)`,
		`four man crew`,
		`(?<=the )(.*?)(?=\s* difficulty)`,
	} {
		res := regexp.MustCompile(exp).FindAllString(desc, -1)
		tags = append(tags, res...)
	}

	return tags
}*/

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
		Heists:       make(map[string]map[string]bool),
	}
	for i := range g.Game.AvailableGameStats.Achievements {
		if g.Game.AvailableGameStats.Achievements[i].Name != p1.Playerstats.Achievements[i].ApiName {
			return nil, err
		}
		a := g.Game.AvailableGameStats.Achievements[i]
		data.Achievements[i] = achieveNodeFinal{
			Name:        a.Name,
			DisplayName: a.DisplayName,
			Description: a.Description,
			Hidden:      a.Hidden > 0,
			Icon:        a.Icon,
			IconGray:    a.IconGray,
			Achieved:    p1.Playerstats.Achievements[i].Achieved > 0,
			Unlocktime:  p1.Playerstats.Achievements[i].Unlocktime,
			//Tags:        getTags(a.Description),
		}

		res := regexp.MustCompile(`Complete the (.*?)(?: job)*? on the (.*?) difficulty or above.`).FindStringSubmatch(a.Description)
		if res != nil {
			if data.Heists[res[1]] == nil {
				data.Heists[res[1]] = make(map[string]bool)
			}
			data.Heists[res[1]][res[2]] = data.Achievements[i].Achieved
		}
	}

	newbody, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	return newbody, nil
}

//https://regoio.herokuapp.com/
