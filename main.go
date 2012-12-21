package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "strconv"
    "regexp"
)

type Letter struct {
    letter string
    pageCount int
}

type Alphabet struct {
    letters []Letter
}

func main() {
    var TurkmenAlphabet = NewAlphabet()
    (&TurkmenAlphabet).setPageCounts()

    for _, l := range TurkmenAlphabet.letters {
        for i := 1; i <= l.pageCount; i++ {
            resp := queryByLetterAndPage(l.letter, i)
            words := parseWords(resp)
            for _, word := range words {
                fmt.Println(word)
            }
        }
    }
}

func NewAlphabet() *Alphabet {
    var alphabet = [...]string{
        "a", "b", "ç", "d", "e", "ä",
        "f", "g", "h", "i", "j", "ž",
        "k", "l", "m", "n", "ň", "o",
        "ö", "p", "r", "s", "ş", "t",
        "u", "ü", "w", "y", "ý",
    }

    TurkmenAlphabet := Alphabet{
        letters: make([]Letter, len(alphabet)),
    }

    for i, l := range alphabet {
        TurkmenAlphabet.letters[i] = Letter{l, 0}
    }

    return &TurkmenAlphabet
}

func (alph *Alphabet) setPageCounts() {
    for i, l := range alph.letters {
        alph.letters[i].pageCount =
            parsePagination(l.letter, queryByLetter(l.letter))
    }
}

func parsePagination(letter string, body []byte) int {
    s := string(body)
    pageCount := 0
    numbers := [20]int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,16,18,19,20}

    for _, page := range numbers {
        paginationUrl := "http://enedilim.com/sozluk/harp/" + letter + "/" + strconv.Itoa(page)

        if strings.Contains(s, paginationUrl) {
            pageCount = page
        }
    }
    return pageCount
}

func parseWords(body []byte) []string  {
    s := string(body)
    r, _ := regexp.Compile("http://enedilim.com/sozluk/soz/[a-zA-ZçäžşýüöÇÄŽŞÝÜÖňŇ-]+")
    results := r.FindAllString(s, -1)
    for i, result := range results {
        results[i] = strings.Replace(result, "http://enedilim.com/sozluk/soz/", "", 1)
    }
    return results
}

func queryByLetter(letter string) []byte {
    baseUrl := "http://enedilim.com/sozluk/harp"
    return query(baseUrl + "/" + letter)
}

func queryByLetterAndPage(letter string, page int) []byte {
    baseUrl := "http://enedilim.com/sozluk/harp/"
    return query(baseUrl + "/" + letter + "/" + strconv.Itoa(page))
}


func query(url string) []byte {
    resp, err := http.Get(url)
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil
    }

    if err != nil {
        fmt.Println(err)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    return body
}
