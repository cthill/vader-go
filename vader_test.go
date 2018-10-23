package vader

import (
	"fmt"
	"github.com/iseurie/vader-go"
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
	fmt.Println("#. cmpd; neg.; pos.; neu.: input")
	for _, s := range tricky {
		st := vader.NewSentiText(s)
		S := st.ScoreValence()
		P := S.Polarity
		fmt.Printf(
			"%d. %.3f; %.3f; %.3f; %.3f: %s\n",
			i+1, P.Compound, S.Negative, S.Positive, S.Neutral, s)
		// sigm mismatch
		if P.Compound*float64(measured) < 0 {
			t.Errorf("Negation test received a mismatch on #%d: ``%s''.",
				i+1, tricky[i])
		}
	}

}
