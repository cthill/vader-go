package vader

import (
	"github.com/iseurie/vader-go/absolutes"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
	"unicode/utf8"
)

// Lexicon represents a mapping from emoji to their respective descriptions,
// and terms to their respective ratings.
type Lexicon struct {
	Emotes  map[rune]string
	Ratings map[string]float64
}

// EmGet retrieves the description of an emoji.
func (L Lexicon) EmGet(e rune) string {
	if v, ok := L.Emotes[e]; ok {
		return strings.Join(tstrip(t), " ")
	}
	return ""
}

// Rates checks whether the lexicon has a rating for a given token.
func (L Lexicon) Rates(s string) bool {
	k := transform.String(s, normalize)
	_, ok := L.Ratings[k]
	return ok
}

// Rating retrieves the rating of a token, returning zero where none is found.
func (L Lexicon) Rating(s string) float64 {
	k := transform.String(s, normalize)
	if v, ok := L.Ratings[k]; ok {
		return v
	}
	return 0
}

// LoadLexicon loads a lexicon into memory. Where either emoji or dict are nil,
// the default resources are restored and used in their stead. Streams provided
// to LoadLexicon must be tab-delimited with one row per line. The emoji
// lexicon, if provided, must map UTF-8 emojis to their descriptions, and the
// lexicon must map English terms to their sentiment ratings as expressed in
// scientific notation.
func LoadLexicon(emoji, lexicon io.Reader) (Lexicon, error) {
	rtn := Lexicon{}
	if lexicon != nil {
		sc := bufio.NewScanner(lexicon)
		for sc.Scan() {
			var t string
			var r float64
			fmt.Fscanf(strings.Trim(sc.Text()), "%s\t%f", &t, &r)
			Lexicon.Ratings[t] = r
		}
		if err := sc.Err(); err != io.EOF {
			return err
		}
	} else {
		Lexicon.Ratings = absolutes.DefaultLexicon.Ratings
	}
	if emoji != nil {
		sc := bufio.NewScanner(emoji)
		for sc.Scan() {
			var emote rune
			var desc string
			fmt.Fscanf(strings.Trim(sc.Text()), "%c\t%s", &emote, &desc)
			Lexicon.Emotes[emote] = desc
		}
		if err := sc.Err(); err != io.EOF {
			return err
		}
	} else {
		Lexicon.Emotes = absolutes.DefaultLexicon.Emotes
	}
	return rtn, nil
}
