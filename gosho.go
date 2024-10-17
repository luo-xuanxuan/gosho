package gosho

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type Word struct {
	Text     string
	Furigana string
	PoS      string
}

type Sentence struct {
	Words []Word
}

var format_url string = "https://jisho.org/search/\"%s\""

func Request(text string) ([]Sentence, error) {
	url := fmt.Sprintf(format_url, text)

	var sentences []Sentence

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, err
	}

	sentencesParent := parseHTMLNode(doc, "section", "zen_bar", "")

	if len(sentencesParent) == 0 {
		return nil, fmt.Errorf("No Zen Bar")
	}

	sentenceNodes := parseHTMLNode(sentencesParent[0], "ul", "", "")

	for _, sentenceNode := range sentenceNodes {
		var sentence Sentence
		wordNodes := parseHTMLNode(sentenceNode, "li", "", "")
		for _, wordNode := range wordNodes {
			var word Word

			word.PoS = getAttr(wordNode, "data-pos")

			furiganaWrapperNode := parseHTMLNode(wordNode, "span", "", "japanese_word__furigana_wrapper")
			if len(furiganaWrapperNode) > 0 {
				word.Furigana = extractText(furiganaWrapperNode[0])
			}

			textWrapperNode := parseHTMLNode(wordNode, "span", "", "japanese_word__text_wrapper")
			if len(textWrapperNode) > 0 {
				word.Text = extractText(textWrapperNode[0])
			}

			sentence.Words = append(sentence.Words, word)
		}
		sentences = append(sentences, sentence)
	}

	return sentences, nil
}
