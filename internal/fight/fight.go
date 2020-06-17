package fight

import "github.com/WSBenson/goku/internal"

// Fighter ... is a struct used by other internal packages to represent a dbz warrior
type Fighter struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

// Fighters ... is a struct used by other internal packages to hold a list of dbz warriors
type Fighters struct {
	Fighters []Fighter `json:"fighters"`
}

// NewFighter ... creates a new dbz warrior
func NewFighter(name string, power int) Fighter {
	f := Fighter{Name: name, Power: power}
	return f
}

// MessageCases ... handles each case for the number of fighters passed
// to the server through a JSON POST request or to elasticsearch with a command
func MessageCases(fs Fighters) string {

	switch length := len(fs.Fighters); {
	case length == 1:
		return evaluatePowerLvl(fs.Fighters[0])

	case length == 2:
		return compareTwoPowers(fs.Fighters[0], fs.Fighters[1])

	case length > 2:
		return compareOverTwoPowers(fs)

	default:
		internal.Logger.Error().Msg("Improperly formated JSON message")
		return "You did not format your JSON message properly, check the console log and the README."
	}
}

// The evaluatePowerLvl function handles the case where one fighter is passed to the server.
// It will return a string that evaluates whether this fighter's power level is over 9000
func evaluatePowerLvl(f Fighter) string {
	if f.Power > 9000 {
		return "The scouter says " + f.Name + "'s power level is over 9000! You better start running."
	} else if f.Power < 206 {
		return f.Name + " is weaker than Krillin, cmon bruh."
	} else {
		return f.Name + "'s power level isn't even over 9000, they're a straight side character."
	}
}

// The compareTwoPowers funcion handles the case where two fighters are passed to the server.
// It will return a string that compares these two fighters' power levels (who is stronger).
func compareTwoPowers(f Fighter, f1 Fighter) string {
	if f.Power == f1.Power && f.Power < 206 && f1.Power < 206 {
		return f.Name + " and " + f1.Name + " are equally trash, they better fuse or something."
	} else if f.Power < 206 && f1.Power < 206 {
		return f.Name + " and " + f1.Name + " are both weaker than Krillin, im done."
	} else if f.Power > f1.Power {
		return f.Name + "'s power level is superior to " + f1.Name + "'s"
	} else if f.Power == f1.Power {
		return f.Name + "'s power level is equal to " + f1.Name + "'s"
	} else {
		return f1.Name + "'s power level is superior to " + f.Name + "'s"
	}
}

// The compareOverTwoPowers functions handles the case where more than two fighters are
// passed to the server. It will return a string stating who has the highest power level.
func compareOverTwoPowers(fs Fighters) string {
	maxPower := 0
	maxFighter := fs.Fighters[0]
	// loops through each fighter to compare their powers
	for _, fighter := range fs.Fighters {
		if fighter.Power > maxPower {
			// Sets the highest power and the fighter who has the highest power
			maxPower = fighter.Power
			maxFighter = fighter
		}
	}

	return maxFighter.Name + " is the strongest of all the Z fighters."
}
