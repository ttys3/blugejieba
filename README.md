# GoJieba Bluge support

[![GoDoc](https://godoc.org/github.com/ttys3/blugejieba?status.svg)](https://godoc.org/github.com/ttys3/blugejieba)
[![Go Report Card](https://goreportcard.com/badge/ttys3/blugejieba)](https://goreportcard.com/report/ttys3/blugejieba)

## Intro

GoJieba [Bluge](https://github.com/blugelabs/bluge) support mod

Bluge is an indexing library for Go. created by the author of the [Bleve](https://github.com/blevesearch/bleve) project

## Get the mod

```bash
go get github.com/ttys3/blugejieba
```

## Usage

```go
	// var ana bluge.Analyzer
    jieba := &analysis.Analyzer{
		Tokenizer: NewJiebaTokenizer(gojieba.DICT_PATH, gojieba.HMM_PATH, gojieba.USER_DICT_PATH, gojieba.IDF_PATH, gojieba.STOP_WORDS_PATH),
	}
	
	// for write
	doc := bluge.NewDocument(msg.Id).
    AddField(bluge.NewTextField("name", "hello bluge").WithAnalyzer(jieba))
    err = writer.Update(doc.ID(), doc)
	
	
	// for read (query)
    query := bluge.NewMatchQuery(q).SetField("name").SetAnalyzer(jieba)
    request := bluge.NewTopNSearch(10, query).
    WithStandardAggregations()
    
    documentMatchIterator, err := reader.Search(context.Background(), request)
    if err != nil {
        log.Fatalf("error executing search: %v", err)
    }
```

please see [jieba_tokenizer_test.go](jieba_tokenizer_test.go)

## docs

https://blugelabs.com/bluge/migration/

https://blugelabs.com/blog/introducing-bluge/

## related project

https://github.com/ttys3/gojieba-bleve