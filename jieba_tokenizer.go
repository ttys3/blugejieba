package blugejieba

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/yanyiwu/gojieba"
)

type JiebaTokenizer struct {
	handle *gojieba.Jieba
}

func (x *JiebaTokenizer) Analyze(input []byte) analysis.TokenStream {
	// TODO implement me
	panic("implement me")
}

var _ analysis.Tokenizer = &JiebaTokenizer{}

func NewJiebaTokenizer(dictpath, hmmpath, userdictpath, idf, stop_words string) *JiebaTokenizer {
	x := gojieba.NewJieba(dictpath, hmmpath, userdictpath, idf, stop_words)
	return &JiebaTokenizer{x}
}

func (x *JiebaTokenizer) Free() {
	x.handle.Free()
}

func (x *JiebaTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	words := x.handle.Tokenize(string(sentence), gojieba.SearchMode, true)
	for _, word := range words {
		token := analysis.Token{
			Term:         []byte(word.Str),
			Start:        word.Start,
			End:          word.End,
			PositionIncr: pos,
			Type:         analysis.Ideographic,
		}
		result = append(result, &token)
		pos++
	}
	return result
}
