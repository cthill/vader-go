package vader

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/iseurie/vader-go/absolutes"
)

var tricky = [...]string{
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

var stockRatings = [...]float64{
	-0.3412, 0.5672, -0.5574, 0.6476, -0.5849, 0.802, -0.2244, -0.2235, 0.2944,
	0.2263, -0.1695, 0.1139, -0.2235,
}

var L = &absolutes.DefaultLexicon

func TestTrickySentences(t *testing.T) {
	fmt.Println("##. cpd.; neg.; pos.; neu.: input")
	tolerance := 0.10
	for i, s := range tricky {
		st := NewSentiText(s, L)
		S := st.ScoreValence()
		P := S.Polarity
		err := math.Abs(S.Compound - stockRatings[i])
		if err > tolerance {
			t.Errorf("#%02d: Error lies %0.2f beyond tolerances",
				i+1, err)
		} else {
			fmt.Printf(
				"%02d. %02.3f; %02.3f; %02.3f; %02.3f: %s\n",
				i+1, S.Compound, P.Negative, P.Positive, P.Neutral, s)
		}
	}
}

func BenchmarkSTConstruction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		si := tricky[rand.Int()%len(tricky)]
		NewSentiText(si, L)
	}
}

func BenchmarkSTScoring(b *testing.B) {
	S := make([]SentiText, len(tricky), len(tricky))
	for i := 0; i < len(tricky); i++ {
		S[i] = NewSentiText(tricky[i], L)
	}

	for i := 0; i < b.N; i++ {
		st := S[rand.Int()%len(tricky)]
		st.ScoreValence()
	}
}
