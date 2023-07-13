package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func ExtractEmail(str string) string {
	left := strings.Index(str, "<")
	right := strings.LastIndex(str, ">")
	if left == -1 || right == -1 {
		// fmt.Println("Cadena inválida")
		return ""
	}
	return string(str[left+1 : right])
}

func ExtactSubDomain(email string) string {
	atSign := strings.Index(email, "@")
	if atSign == -1 {
		fmt.Println("Dirección de correo electrónico inválida")
		return ""
	}
	domain := email[atSign+1:]
	dot := strings.IndexByte(domain, '.')
	if dot == -1 {
		fmt.Println("Dominio de correo electrónico inválido")
		return ""
	}
	return strings.ReplaceAll(domain[:dot], "a", "")
}

func GetAgency(keyAgency string) string {
	jsonFile, err := ioutil.ReadFile("./locale/agency.json")
	if err != nil {
		return "./locale/agency.json not found"
	}
	var data map[string]interface{}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return "format error in locale file"
	}
	value, ok := data[keyAgency]
	if !ok {
		return "key not found" // Clave no encontrada
	}
	stringValue, ok := value.(string)
	if !ok {
		return "Value is not a string" // Valor no es una cadena de texto
	}
	return stringValue
}

func ConfigEnv(data string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	return os.Getenv(data)
}
