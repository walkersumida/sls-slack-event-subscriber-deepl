package ja

import "regexp"

func IsJA(text string) bool {
	re := regexp.MustCompile("[ぁ-んァ-ヶー一-龠々]")
	return re.MatchString(text)
}
