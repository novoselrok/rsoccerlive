package highlights

import (
	"testing"

	"github.com/novoselrok/rsoccerlive/src/redditclient"
)

func TestMatchingHighlightTitles(t *testing.T) {
	titles := []string{
		"Barcelona 2-[7] Bayern Munich: Philippe Coutinho goal 85'",
		"Atalanta 1-[2] Paris Saint-Germain: Choupo-Moting goal 90+3'",
		"Orlando City 2-0 Minnesota United - Nani 42’ (great goal!)",
		"Barcelona 0-1 Bayern Munich: Thomas Muller goal 4'",
		"RB Leipzig [2]-1 Atletico Madrid: Tyler Adams goal 88'",
		"Benfica (1) - 1 Braga - Vinicius 56'",
	}

	for _, title := range titles {
		if !isHighlightTitle(title) {
			t.Errorf("Title [%s] should be a highlight title.", title)
		}
	}
}

func TestNonMatchingHighlightTitles(t *testing.T) {
	titles := []string{
		"Harry Maguire 'fought outside Greek bar after his sister was stabbed in the arm'",
		"PSG qualify to the Champions League semi-final for the first time since 1994",
		"Atlético Madrid | Oblak: \"If you don’t concede you only have to score one\"",
		"Fulham will play Brentford in the Championship playoff final",
	}

	for _, title := range titles {
		if isHighlightTitle(title) {
			t.Errorf("Title [%s] should not be a highlight title.", title)
		}
	}
}

func TestMatchingHighlightSubmission(t *testing.T) {
	submissions := []redditclient.Submission{
		{"id1", "https://streamja.com/5rJzl", "Atalanta 1-0 PSG - Mario Pasalic 27'", "/r/soccer/comments/i8kebn/atalanta_10_psg_mario_pasalic_27/", "author1", false, 1},
		{"id2", "https://streamable.com/pvyeij", "Barcelona 2-[8] Bayern Munich: Philippe Coutinho goal 88'", "/r/soccer/comments/i9txqv/barcelona_28_bayern_munich_philippe_coutinho_goal/", "author2", false, 1},
	}

	for _, submission := range submissions {
		if !isHighlightSubmission(submission) {
			t.Errorf("Submission [%+v] should be a highlight submission", submission)
		}
	}
}

func TestNonMatchingHighlightSubmission(t *testing.T) {
	submissions := []redditclient.Submission{
		{"id1", "https://twitter.com/HarrogateTown/status/1289952540751097856?s=19", "Harrogate Town have been promoted to the EFL for the first time in their 106 year history", "/r/soccer/comments/i2eafw/harrogate_town_have_been_promoted_to_the_efl_for/", "author1", false, 1},
		{"id2", "https://www.quotidiano.net/sport/calcio/moratti-inter-1.5414427", "Historic ex Inter owner Moratti asked about Juve: \"Its one thing to get eliminated by PSG, another to get eliminated by the 7th placed team in Ligue 1. Honestly, there is no comparison.\"", "/r/soccer/comments/i9jbwy/historic_ex_inter_owner_moratti_asked_about_juve/", "author2", false, 1},
	}

	for _, submission := range submissions {
		if isHighlightSubmission(submission) {
			t.Errorf("Submission [%+v] should not be a highlight submission", submission)
		}
	}
}
