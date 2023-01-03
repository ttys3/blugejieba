package blugejieba

import (
	"context"
	"fmt"
	"github.com/blugelabs/bluge/analysis"
	"github.com/yanyiwu/gojieba"
	"log"
	"os"
	"testing"

	"github.com/blugelabs/bluge"
)

func TestBluge(t *testing.T) {
	indexDir := ".gojieba.bluge"
	config := bluge.DefaultConfig(indexDir)
	// config.DefaultSearchAnalyzer = xxx
	writer, err := bluge.OpenWriter(config)
	if err != nil {
		log.Fatalf("error opening writer: %v", err)
	}
	defer writer.Close()

	doc := bluge.NewDocument("example_doc_id001").
		AddField(bluge.NewTextField("name", "bluge"))

	err = writer.Update(doc.ID(), doc)
	if err != nil {
		log.Fatalf("error updating document: %v", err)
	}

	reader, err := writer.Reader()
	if err != nil {
		log.Fatalf("error getting index reader: %v", err)
	}
	defer reader.Close()

	query := bluge.NewMatchQuery("bluge").SetField("name")
	request := bluge.NewTopNSearch(10, query).
		WithStandardAggregations()

	documentMatchIterator, err := reader.Search(context.Background(), request)
	if err != nil {
		log.Fatalf("error executing search: %v", err)
	}
	match, err := documentMatchIterator.Next()
	for err == nil && match != nil {
		err = match.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				fmt.Printf("match _id: %s\n", string(value))
			}
			return true
		})
		if err != nil {
			log.Fatalf("error loading stored fields: %v", err)
		}
		match, err = documentMatchIterator.Next()
	}
	if err != nil {
		log.Fatalf("error iterator document matches: %v", err)
	}
}

func ExampleBlugeJieba() {
	indexDir := ".gojieba.bluge"
	messages := []struct {
		Id   string
		Body string
	}{
		{
			Id:   "1",
			Body: "你好",
		},
		{
			Id:   "2",
			Body: "交代",
		},
		{
			Id:   "3",
			Body: "长江大桥",
		},
	}

	os.RemoveAll(indexDir)

	// var ana bluge.Analyzer
	ana := &analysis.Analyzer{
		Tokenizer: NewJiebaTokenizer(gojieba.DICT_PATH, gojieba.HMM_PATH, gojieba.USER_DICT_PATH, gojieba.IDF_PATH, gojieba.STOP_WORDS_PATH),
	}

	config := bluge.DefaultConfig(indexDir)
	// config.DefaultSearchAnalyzer = ana
	writer, err := bluge.OpenWriter(config)
	if err != nil {
		log.Fatalf("error opening writer: %v", err)
	}
	defer writer.Close()

	for _, msg := range messages {
		doc := bluge.NewDocument(msg.Id).
			AddField(bluge.NewTextField("name", msg.Body).WithAnalyzer(ana))
		err = writer.Update(doc.ID(), doc)
		if err != nil {
			log.Fatalf("error updating document: %v", err)
		}
	}

	// clean index when example finished
	defer os.RemoveAll(indexDir)

	querys := []string{
		"你好世界",
		"亲口交代",
		"长江",
	}

	reader, err := writer.Reader()
	if err != nil {
		log.Fatalf("error getting index reader: %v", err)
	}
	defer reader.Close()

	for _, q := range querys {
		query := bluge.NewMatchQuery(q).SetField("name").SetAnalyzer(ana)
		request := bluge.NewTopNSearch(10, query).
			WithStandardAggregations()

		documentMatchIterator, err := reader.Search(context.Background(), request)
		if err != nil {
			log.Fatalf("error executing search: %v", err)
		}
		match, err := documentMatchIterator.Next()
		for err == nil && match != nil {
			err = match.VisitStoredFields(func(field string, value []byte) bool {
				if field == "_id" {
					fmt.Printf("match _id: %s, _score: %v\n", string(value), match.Score)
				}
				return true
			})
			if err != nil {
				log.Fatalf("error loading stored fields: %v", err)
			}
			match, err = documentMatchIterator.Next()
		}
		if err != nil {
			log.Fatalf("error iterator document matches: %v", err)
		}
	}

	// ana.(*analysis.Analyzer).Tokenizer.(*JiebaTokenizer).Free()
	ana.Tokenizer.(*JiebaTokenizer).Free()

	// Output:
	// match _id: 1, _score: 0.6543330458293132
	// match _id: 2, _score: 0.6543330458293132
	// match _id: 3, _score: 0.4123194535362795
}
