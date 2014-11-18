package trello_markov

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	trelloAPI = "https://api.trello.com"
)

type TrelloCorpus struct {
	apiKey   string
	apiToken string
}

func NewTrelloCorpus(apiKey string, apiToken string) *TrelloCorpus {
	return &TrelloCorpus{
		apiKey:   apiKey,
		apiToken: apiToken,
	}
}

func (t *TrelloCorpus) TestCorpus() (<-chan []byte, error) {
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	c := make(chan []byte, 32)
	go func() {
		defer close(c)
		for _, f := range dir {
			if f.IsDir() || !strings.HasPrefix(f.Name(), "trello_") {
				continue
			}
			b, err := ioutil.ReadFile(f.Name())
			if err != nil {
				return
			}
			c <- b
		}
	}()
	return c, nil
}

func (t *TrelloCorpus) TextCorpus(boardID string) (<-chan []byte, error) {
	cards, err := t.FetchCards(boardID)
	if err != nil {
		return nil, err
	}
	c := make(chan []byte, len(cards))
	go func() {
		defer close(c)
		for _, card := range cards {
			card.Corpus(c)
			comments, err := t.FetchCardComments(card.ID)
			if err != nil {
				return
			}
			for _, comment := range comments {
				comment.Corpus(c)
			}
			//time.Sleep(250 * time.Millisecond)
		}
	}()
	//fmt.Println(buf.String())
	return c, nil
}

type Card struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Checklists []struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"checkItems"`
	} `json:"checklists"`
}

func (c Card) Corpus(b chan<- []byte) {
	b <- []byte(strings.TrimSpace(c.Name))
	b <- []byte(strings.TrimSpace(c.Desc))
	for _, l := range c.Checklists {
		for _, i := range l.Items {
			b <- []byte(strings.TrimSpace(i.Name))
		}
	}
}

func (t *TrelloCorpus) FetchCards(boardID string) ([]Card, error) {
	u, _ := url.Parse(trelloAPI)
	u.Path = fmt.Sprintf("/1/boards/%s/cards", boardID)
	q := u.Query()
	q.Set("key", t.apiKey)
	q.Set("token", t.apiToken)
	q.Set("fields", "desc,name")
	q.Set("filter", "all")
	q.Set("checklists", "all")
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s -> %s", u.String(), resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	var cards []Card
	if err = d.Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

type CardComment struct {
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

func (c CardComment) Corpus(b chan<- []byte) {
	if !strings.HasPrefix(c.Data.Text, ":octocat:") {
		b <- []byte(strings.TrimSpace(c.Data.Text))
	}
}

func (t *TrelloCorpus) FetchCardComments(cardID string) ([]CardComment, error) {
	u, _ := url.Parse(trelloAPI)
	u.Path = fmt.Sprintf("/1/cards/%s/actions", cardID)
	q := u.Query()
	q.Set("key", t.apiKey)
	q.Set("token", t.apiToken)
	q.Set("filter", "commentCard")
	q.Set("fields", "data")
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s -> %s", u.String(), resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	var comments []CardComment
	if err = d.Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}
