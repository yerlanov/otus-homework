package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mailru/easyjson"
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
	u := getUsers(r)

	return countDomains(u, domain)
}

func getUsers(r io.Reader) <-chan User {
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

	return ch
}

func countDomains(usersChan <-chan User, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}

	result := make(DomainStat)

	for user := range usersChan {
		if user.Email == "" {
			continue
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
