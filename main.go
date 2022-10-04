package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var serverScrap = "https://www.freevmess.com/"
var countryServer = "unitedstates-v2ray-server"
var Exit = make(chan bool)
var number = [2]int64{0, 0}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func answerMathProblem(str string) int64 {
	var aw string
	in := 0
	str = strings.TrimSpace(str)
	for i := 0; i < len(str)-1; i++ {
		if string(str[i]) == " " {
			in++
		} else if string(str[i]) != "" && (!(string(str[i]) >= "0" && string(str[i]) <= "9")) {
			aw += string(str[i])
			in--
		} else {
			number[in] *= 10
			ans, _ := strconv.Atoi(string(str[i]))
			number[in] += int64(ans)
		}
	}
	switch aw {
	case "+":
		return number[0] + number[1]
	case "-":
		return number[0] - number[1]
	case "/":
		return number[0] / number[1]
	case "*":
		return number[0] * number[1]
	}
	return 0
}

func getFirstCharFromString(str string, str2 string) int {
	for i := 0; i < len(str)-len(str2); i++ {
		if string(str[i:i+len(str2)]) == str2 {
			return i
		}
	}
	return -1
}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func answerCapchaFromHtml(html string) string {
	var answer string
	for i := 0; i < len(html)-7; i++ {
		if html[i:i+7] == "What is" {
			a := getFirstCharFromString(html[i+8:], "?")
			if a == -1 {
				return ""
			}
			answer = strconv.FormatInt(answerMathProblem(string(html[i+8:i+8+a])), 10)
			break
		}
	}
	return answer
}

func generateVMESS() error {
	resp, _ := http.Get(serverScrap)

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	capcha := answerCapchaFromHtml(string(html))

	url := serverScrap + countryServer

	method := "POST"

	name := RandStringRunes(5) + "0"

	payload := strings.NewReader(fmt.Sprintf("name=%s&nameser=singapore&id_server=33&proov=v2ray&firstNumber=%s&secondNumber=%s&captcha=%s&submit=",
		name, strconv.FormatInt(number[0], 10), strconv.FormatInt(number[1], 10), capcha))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Origin", serverScrap)
	req.Header.Add("Referer", serverScrap+countryServer)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	ind := getFirstCharFromString(string(body), "vmess://")

	fmt.Println(string(body)[ind : ind+getFirstCharFromString(string(body)[ind:], `"`)])

	Exit <- true
	return nil
}

//soon
func generateVLESS() error {
	return nil
}

func main() {
	_, err := http.Get(serverScrap + countryServer)
	if err != nil {
		panic("we have problem . plaese try again")
	}

	go func() {
		fmt.Println("Plaese select protocol for vpn\n1.vmess\nEnter number:")

		var po int = 0
		for po > 1 || po < 1 {
			po, _ = fmt.Scanf("%d", &po)
		}
		fmt.Println("loading...")

		switch po {
		case 1:
			_ = generateVMESS()
			break
		default:
			_ = generateVLESS()
			break
		}
	}()

	<-Exit
}
