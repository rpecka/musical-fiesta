package commands

var (
	blacklist = map[string]bool{
		"attack": true, "attack2": true, "autobuy": true, "back": true, "buy": true, "buyammo1": true, "buyammo2": true,
		"buymenu": true, "callvote": true, "cancelselect": true, "cheer": true, "compliment": true, "coverme": true,
		"drop": true, "duck": true, "enemydown": true, "enemyspot": true, "fallback": true, "followme": true,
		"forward": true, "getout": true, "go": true, "holdpos": true, "inposition": true, "invnext": true,
		"invprev": true, "jump": true, "lastinv": true, "messagemode": true, "messagemode2": true, "moveleft": true,
		"moveright": true, "mute": true, "negative": true, "quit": true, "radio1": true, "radio2": true, "radio3": true,
		"rebuy": true, "regroup": true, "reload": true, "report": true, "reportingin": true, "roger": true,
		"sectorclear": true, "showscores": true, "slot1": true, "slot10": true, "slot2": true, "slot3": true,
		"slot4": true, "slot5": true, "slot6": true, "slot7": true, "slot8": true, "slot9": true, "speed": true,
		"sticktog": true, "takepoint": true, "takingfire": true, "teammenu": true, "thanks": true,
		"toggleconsole": true, "use": true, "voicerecord": true,
	}
)

func IsIllegal(command string) bool {
	if _, value := blacklist[command]; value {
		return true
	} else {
		return false
	}
}
