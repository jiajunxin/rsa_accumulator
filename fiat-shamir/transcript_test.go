package fiatshamir

import "testing"

func TestTranscript(t *testing.T) {
	testStrings := []string{"111", "222", "333"}
	trans1 := InitTranscript(testStrings)

	var trans2 Transcript
	trans2.AppendSlice(testStrings)

	challenge1 := trans1.GetChallengeAndAppendTranscript()
	challenge2 := trans2.GetChallengeAndAppendTranscript()

	if challenge1.Cmp(challenge2) != 0 {
		t.Errorf("Different ways to init transcript leads to different results")
		trans1.Print()
		trans2.Print()
	}
}
