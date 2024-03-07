package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"log"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) (<-chan User, error) {
	ch := make(chan User)

	go func() {
		defer close(ch)

		user := User{}
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
				log.Println("Error unmarshalling user:", err)
				continue
			}
			ch <- user
		}

		if err := scanner.Err(); err != nil {
			log.Println("Error reading users:", err)
		}
	}()

	return ch, nil
}

func countDomains(usersChan <-chan User, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for user := range usersChan {
		if user.Email != "" && domain != "" && strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
