package main

import (
	"regexp"
	"testing"
)

var appid string = "218620"
var steamid string = "76561198049201876"

func TestCallSteam(t *testing.T) {
	var p1 player
	err := callSteam("GetPlayerAchievements/v1", appid, steamid, &p1)
	if err != nil || regexp.MustCompile(`\b`+steamid+`\b`).MatchString(p1.Playerstats.SteamID) {
		t.Fatalf("Failed to call steam: %s", err)
	}
}

func TestGetTags(t *testing.T) {

}

func TestGetPlayerProgress(t *testing.T) {

}
