package services
import (
	"crypto/rand"
	"log"
)

func CreateShortLink(link string, size int) string{
    randomBytes := make([]byte, size)
    _, err := rand.Read(randomBytes)
	// TODO ЛОГИРУЕМ ОИШИБКУ, ЭТО ПЛОХО ИСПРАВИТЬ
    if err != nil {
        log.Fatal(err)
    }
	return string(randomBytes)
}