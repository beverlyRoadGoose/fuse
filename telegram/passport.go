package telegram // import "heytobi.dev/fuse/telegram"

// PassportData contains information about Bot Passport data shared with the bot by the user.
// See https://core.telegram.org/bots/api#passportfile
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials *EncryptedCredentials      `json:"credentials"`
}

// PassportFile represents a file uploaded to Bot Passport.
// See https://core.telegram.org/bots/api#passportdata
type PassportFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileDate     int    `json:"file_date"`
}

// EncryptedPassportElement contains information about documents or other Bot Passport elements shared with the bot
// by the user.
// See https://core.telegram.org/bots/api#encryptedpassportelement
type EncryptedPassportElement struct {
	Type        string         `json:"type"`
	Data        string         `json:"data"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Files       []PassportFile `json:"files"`
	FrontSide   *PassportFile  `json:"front_side"`
	ReverseSide *PassportFile  `json:"reverse_side"`
	Selfie      *PassportFile  `json:"selfie"`
	Translation []PassportFile `json:"translation"`
	Hash        string         `json:"hash"`
}

// EncryptedCredentials contains data required for decrypting and authenticating EncryptedPassportElement.
// See https://core.telegram.org/bots/api#encryptedcredentials
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}
