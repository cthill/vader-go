package absolutes

import (
	"github.com/iseurie/vader-go"
	"strings"
)

const (
	boostIncr float64 = 0.293
	boostDecr float64 = -0.293

	// CFactor is the empirically-determined capitalization emphasis scaling factor.
	CFactor float64 = 0.733

	// NScalar is an enigma: I'm not sure what the h*ck this number is; if you
	// find out, let me know!
	NScalar float64 = -0.74
)

// DefaultLexicon is the default lexicon, indexing emoji against their ratings and
// all known terms against their sentiment ratings, respectively.
var DefaultLexicon = Lexicon{
	Ratings: terms,
	Emotes:  emotes,
}

// SentimentLadenIdioms indexes and scales a set of sentiment-laden idioms
// which aren't accounted for by the lexicon.
var SentimentLadenIdioms = map[string]float64{
	"cut the mustard": 2, "hand to mouth": -2,
	"back handed": -2, "blow smoke": -2,
	"blowing smoke": -2, "upper hand": 1,
	"break a leg": 2, "cooking with gas": 2,
	"in the black": 2, "in the red": -2,
	"on the ball": 2, "under the weather": -2,
}

// Punctuation contains a list of common punctuation.
var Punctuation = [...]string{
	".", "!", "?", ",", ";", ":", "-", "'", "\"", "!!", "!!!", "??", "???",
	"?!?", "!?!", "?!?!", "!?!?",
}

// SpecialCaseIdioms indexes and scales a set of idioms that have antithetical
// sentiment to that set forth by the lexicon. Extending this set is a good
// candidate for further work.
var SpecialCaseIdioms = map[string]float64{
	"the shit": 3, "the bomb": 3, "bad ass": 1.5, "yeah right": -2,
	"kiss of death": -1.5,
}

var negate = map[string]struct{}{
	"aint":      struct{}{},
	"arent":     struct{}{},
	"cannot":    struct{}{},
	"cant":      struct{}{},
	"couldnt":   struct{}{},
	"darent":    struct{}{},
	"didnt":     struct{}{},
	"doesnt":    struct{}{},
	"ain't":     struct{}{},
	"aren't":    struct{}{},
	"can't":     struct{}{},
	"couldn't":  struct{}{},
	"daren't":   struct{}{},
	"didn't":    struct{}{},
	"doesn't":   struct{}{},
	"dont":      struct{}{},
	"hadnt":     struct{}{},
	"hasnt":     struct{}{},
	"havent":    struct{}{},
	"isnt":      struct{}{},
	"mightnt":   struct{}{},
	"mustnt":    struct{}{},
	"neither":   struct{}{},
	"don't":     struct{}{},
	"hadn't":    struct{}{},
	"hasn't":    struct{}{},
	"haven't":   struct{}{},
	"isn't":     struct{}{},
	"mightn't":  struct{}{},
	"mustn't":   struct{}{},
	"neednt":    struct{}{},
	"needn't":   struct{}{},
	"never":     struct{}{},
	"none":      struct{}{},
	"nope":      struct{}{},
	"nor":       struct{}{},
	"not":       struct{}{},
	"nothing":   struct{}{},
	"nowhere":   struct{}{},
	"oughtnt":   struct{}{},
	"shant":     struct{}{},
	"shouldnt":  struct{}{},
	"uhuh":      struct{}{},
	"wasnt":     struct{}{},
	"werent":    struct{}{},
	"oughtn't":  struct{}{},
	"shan't":    struct{}{},
	"shouldn't": struct{}{},
	"uh-uh":     struct{}{},
	"wasn't":    struct{}{},
	"weren't":   struct{}{},
	"without":   struct{}{},
	"wont":      struct{}{},
	"wouldnt":   struct{}{},
	"won't":     struct{}{},
	"wouldn't":  struct{}{},
	"rarely":    struct{}{},
	"seldom":    struct{}{},
	"despite":   struct{}{},
}

var boostUp = map[string]struct{}{
	"absolutely":    struct{}{},
	"amazingly":     struct{}{},
	"awfully":       struct{}{},
	"completely":    struct{}{},
	"considerably":  struct{}{},
	"decidedly":     struct{}{},
	"deeply":        struct{}{},
	"effing":        struct{}{},
	"enormously":    struct{}{},
	"entirely":      struct{}{},
	"especially":    struct{}{},
	"exceptionally": struct{}{},
	"extremely":     struct{}{},
	"fabulously":    struct{}{},
	"flipping":      struct{}{},
	"flippin":       struct{}{},
	"fricking":      struct{}{},
	"frickin":       struct{}{},
	"frigging":      struct{}{},
	"friggin":       struct{}{},
	"fully":         struct{}{},
	"fucking":       struct{}{},
	"greatly":       struct{}{},
	"hella":         struct{}{},
	"highly":        struct{}{},
	"hugely":        struct{}{},
	"incredibly":    struct{}{},
	"intensely":     struct{}{},
	"majorly":       struct{}{},
	"more":          struct{}{},
	"most":          struct{}{},
	"particularly":  struct{}{},
	"purely":        struct{}{},
	"quite":         struct{}{},
	"really":        struct{}{},
	"remarkably":    struct{}{},
	"so":            struct{}{},
	"substantially": struct{}{},
	"thoroughly":    struct{}{},
	"totally":       struct{}{},
	"tremendously":  struct{}{},
	"uber":          struct{}{},
	"unbelievably":  struct{}{},
	"unusually":     struct{}{},
	"utterly":       struct{}{},
	"very":          struct{}{},
}

var boostDn = map[string]struct{}{
	"almost":       struct{}{},
	"barely":       struct{}{},
	"hardly":       struct{}{},
	"just enough":  struct{}{},
	"kind of":      struct{}{},
	"kinda":        struct{}{},
	"kindof":       struct{}{},
	"kind-of":      struct{}{},
	"less":         struct{}{},
	"little":       struct{}{},
	"marginally":   struct{}{},
	"occasionally": struct{}{},
	"partly":       struct{}{},
	"scarcely":     struct{}{},
	"slightly":     struct{}{},
	"somewhat":     struct{}{},
	"sort of":      struct{}{},
	"sorta":        struct{}{},
	"sortof":       struct{}{},
	"sort-of":      struct{}{},
}

// Boost returns whether and how much of a boosting factor a given term receives.
func Boost(t string) float64 {
	k := strings.ToLower(t)
	if _, ok := boostUp[k]; ok {
		return boostIncr, ok
	}
	if _, ok := boostDn[k]; ok {
		return boostDecr, ok
	}
	return 0, true
}

// Returns whether a given token has a boost entry.
func IsBoosted(t string) bool {
	k := strings.ToLower(t)
	_, u := boostUp[k]
	_, d := boostDn[k]
	return u || d
}

// Negate returns whether a given term constitutes a negation.
func Negate(t string) bool {
	k := strings.ToLower(t)
	_, ok := negate[k]
	return ok
}
