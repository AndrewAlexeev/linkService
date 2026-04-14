package config
import(
	"os"
)
type Config struct{
		Port string
		DbHost string
    	DbPort string
    	DbUser string
		DbName string
		DbPassword string


}

func InitConfig() Config{
	config := Config{
		Port : os.Getenv("PORT"),
		DbHost : os.Getenv("DB_HOST"),
		DbPort : os.Getenv("DB_PORT"),
		DbPassword : os.Getenv("DB_PASSWORD"),
		DbUser :  os.Getenv("DB_USER"),
		DbName : os.Getenv("DB_NAME")}
	return config

}