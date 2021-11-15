package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "encoding/json"
)

type uR struct {
	Log string `json:"login"`
	Pas string `json:"password"`
}

type bOut struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

func testSign(t string, log string, pas string) {
	fmt.Println("Тест:", t)
	fmt.Println("Адрес", "http://localhost:8080/api/user/"+t)
	u := uR{
		Log: log,
		Pas: pas,
	}
	fmt.Println("Данные:", u)
	reqBody, err := json.Marshal(&u)
	if err != nil {
		print(err)
	}
	a := "http://localhost:8080/api/user/" + t
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "application/json", "", "POST", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "application/json", "", "POST", reqBody)
	fmt.Println("Окончание теста")
}

func makeRequest(address, ctype, key, rtype string, b []byte) {
	client := &http.Client{}
	req, _ := http.NewRequest(rtype, address, bytes.NewReader(b))
	req.Header.Add("Content-Type", ctype)
	req.Header.Add("Authorization", key)
	// r, err := http.Post(a, t, bytes.NewBuffer(b)) //bytes.NewBuffer(reqBody))
	r, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		print(err)
	}
	printResult(r.Body, r)
}

func makeZipPostRequest(address, ctype, key, rtype string, reqBody []byte) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	gz.Write(reqBody)
	gz.Flush()
	gz.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(rtype, address, bytes.NewReader(b.Bytes()))
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", ctype)
	req.Header.Add("Authorization", key)

	r, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	gzr, err := gzip.NewReader(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	printResult(gzr, r)
}

func printResult(body io.Reader, r *http.Response) {
	text, err := io.ReadAll(body)
	if err != nil {
		print(err)
	}
	defer r.Body.Close()
	fmt.Printf("%s\n", r.Header)
	fmt.Printf("%s\n", text)
	fmt.Printf("%d\n", r.StatusCode)
}

func main() {
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200(на новой базе) 409(на старой), 409", "успешно, уже есть")
	testSign("register", "Aleha", "123123213")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200(на новой базе) 409(на старой), 409", "успешно, уже есть")
	testSign("register", "Kartoha", "457457457457")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200, 200", "Успешно")
	testSign("login", "Aleha", "123123213")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200, 200", "Успешно")
	testSign("login", "Kartoha", "457457457457")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 401, 401", "Неверная пара логи, пароль")
	testSign("login", "Karas", "457457457457")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	// testAccrualPost("1230")

}
