package SteamAchievementProgressGolang

import (
	"testing"
)

var appid string = "218620"
var steamid string = "76561198049201876"

func TestCallSteam(t *testing.T) {
	var p1 player
	err := callSteam("GetPlayerAchievements/v1", appid, steamid, &p1)
	if err != nil {
		t.Fatalf("callSteam returned error")
	}
	if p1.Playerstats.SteamID != steamid {
		t.Fatalf("callSteam did not populate p1")
	}

	var p2 player
	err = callSteam("GetPlayerAchievements/v1", appid, "fdsa", &p2)
	if err == nil {
		t.Fatalf("callSteam failed to fail on bad steamid")
	}
	var p3 player
	err = callSteam("bad address", appid, steamid, &p3)
	if err == nil {
		t.Fatalf("callSteam failed to fail on bad address")
	}
}

func TestGetPlayerProgress(t *testing.T) {
	data, err := GetPlayerProgress(appid, steamid)

	if err != nil || len(data) <= 0 {
		t.Fatalf("Failed to get player progress")
	}

	_, err = GetPlayerProgress(appid, "bad steamid")
	if err == nil {
		t.Fatalf("Failed to fail on bad steamid")
	}
}

func TestCheckHeistComplete(t *testing.T) {
	heists := make(map[string][]bool)

	testString := "Complete the Dragon Heist job in 6 minutes or less on the OVERKILL difficulty or above."
	checkHeistComplete(heists, testString, false)
	if heists["Dragon Heist"] != nil {
		t.Fatalf("Failed to HIDE edge case (non-completion achievement)")
	}

	testString = "Complete the Dragon Heist job on the OVERKILL difficulty or above."
	checkHeistComplete(heists, testString, false)
	if heists["Dragon Heist"] == nil || heists["Dragon Heist"][3] != false {
		t.Fatalf("Failed to find basic heist")
	}

	testString = "Complete the Ukrainian Job on the Death Sentence difficulty with the One Down mechanic activated."
	checkHeistComplete(heists, testString, false)
	if heists["Ukrainian Job"] == nil {
		t.Fatalf("Failed to find edge case (Ukrainian Job)")
	}
	if heists["Ukrainian Job"][3] != false {
		t.Fatalf("Failed to find One Down difficulty")
	}

	testString = "Complete The Alesso Heist on the Death Wish difficulty or above."
	checkHeistComplete(heists, testString, false)
	if heists["Alesso Heist"] == nil {
		t.Fatalf("Failed to find edge case (Alesso Heist)")
	}
}
