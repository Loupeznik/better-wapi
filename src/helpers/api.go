package helpers

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func GetApiToken(username string, password string) string {
	passwordHash := sha1.New()
	passwordHash.Write([]byte(password))
	location, _ := time.LoadLocation("Europe/Prague")
	hour := formatHour(time.Now().In(location).Hour())

	passwordHashString := fmt.Sprintf("%x", passwordHash.Sum(nil))

	token := fmt.Sprintf("%s%s%s", username, passwordHashString, hour)

	tokenHash := sha1.New()
	tokenHash.Write([]byte(token))

	return fmt.Sprintf("%x", tokenHash.Sum(nil))
}

func formatHour(hour int) string {
	var formattedHour string

	if hour < 10 {
		formattedHour = fmt.Sprintf("0%d", hour)
	} else {
		formattedHour = fmt.Sprintf("%d", hour)
	}

	return formattedHour
}
