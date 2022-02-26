package telegram // import "heytobi.dev/fuse/telegram"

// PollOption contains information about one answer option in a poll.
// See https://core.telegram.org/bots/api#polloption
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
// See https://core.telegram.org/bots/api#pollanswer
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	User      *User  `json:"user"`
	OptionIDs []int  `json:"option_ids"`
}

// Poll contains information about a poll.
// See https://core.telegram.org/bots/api#poll
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`
	CorrectOptionID       int             `json:"correct_option_id"`
	Explanation           string          `json:"explanation"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities"`
	OpenPeriod            int             `json:"open_period"`
	CloseDate             int             `json:"close_date"`
}
