package vader

import (
	"fmt"
	"github.com/iseurie/vader-go"
	"github.com/iseurie/vader-go/absolutes"
	"testing"
)

func TestTrickySentences(t *testing.T) {
	tricky := [...]string{
		"Sentiment analysis has never been good.",
		"Sentiment analysis has never been this good!",
		"Most automated sentiment analysis tools are shit.",
		"With VADER, sentiment analysis is the shit!",
		"Other sentiment analysis tools can be quite bad.",
		"On the other hand, VADER is quite bad ass",
		"VADER is such a badass!", // slang with punctuation emphasis
		"Without a doubt, excellent idea.",
		"Roger Dodger is one of the most compelling variations on this theme.",
		"Roger Dodger is at least compelling as a variation on the theme.",
		"Roger Dodger is one of the least compelling variations on this theme.",
		"Not such a badass after all.",        // Capitalized negation with slang
		"Without a doubt, an excellent idea.", // "without {any} doubt" as negation
	}
	perceived := [...]int{-1, 1, -1, 1, 1, 1, 1, 1, 1, -1, -1, 1}
	measured := make([]float64, cap(perceived), cap(perceived))
	L := absolutes.DefaultLexicon
	fmt.Println("#. cmpd; neg.; pos.; neu.: input")
	for i, s := range tricky {
		st := vader.NewSentiText(s, &L)
		S := st.ScoreValence()
		P := S.Polarity
		fmt.Printf(
			"%d. %.3f; %.3f; %.3f; %.3f: %s\n",
			i+1, S.Compound, P.Negative, P.Positive, P.Neutral, s)
		// signimismatch
		if S.Compound*float64(measured[i]) < 0 {
			t.Errorf("Negation test received a mismatch on #%d: ``%s''.",
				i+1, tricky[i])
		}
	}

}
