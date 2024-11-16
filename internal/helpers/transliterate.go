package helpers

import (
	"strings"
	"unicode"
)

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch", 'ы': "y",
	'э': "e", 'ю': "yu", 'я': "ya", 'ь': "", 'ъ': "",
	'А': "a", 'Б': "b", 'В': "v", 'Г': "g", 'Д': "d", 'Е': "e", 'Ё': "e",
	'Ж': "zh", 'З': "z", 'И': "i", 'Й': "y", 'К': "k", 'Л': "l", 'М': "m",
	'Н': "n", 'О': "o", 'П': "p", 'Р': "r", 'С': "s", 'Т': "t", 'У': "u",
	'Ф': "f", 'Х': "kh", 'Ц': "ts", 'Ч': "ch", 'Ш': "sh", 'Щ': "shch", 'Ы': "y",
	'Э': "e", 'Ю': "yu", 'Я': "ya", ' ': "-",
}

func Transliterate(input string) string {
	var result strings.Builder
	for _, char := range input {
		lowerChar := unicode.ToLower(char)

		if val, ok := translitMap[lowerChar]; ok {
			result.WriteString(val)
		} else if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(lowerChar)
		} else {
			result.WriteString("-")
		}
	}
	return result.String()
}
