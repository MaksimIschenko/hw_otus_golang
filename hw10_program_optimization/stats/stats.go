package stats

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson" //nolint:depguard
)

type DomainStat map[string]int

type LightUser struct {
	Email string `json:"Email"` //nolint:tagliatelle
}

type LightUsers []LightUser

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := GetUserEmails(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func GetUserEmails(r io.Reader, domain string) (LightUsers, error) {
	scanner := bufio.NewScanner(r)
	result := make(LightUsers, 0)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var lUser LightUser
		if err := easyjson.Unmarshal(line, &lUser); err != nil {
			return result, err
		}

		if strings.Contains(lUser.Email, domain) {
			result = append(result, lUser)
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func countDomains(u LightUsers, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
