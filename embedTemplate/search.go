package embedderTemplate

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

func Rag(question string, numOfResults int, api_token string, store vectorstores.VectorStore, baseLink string) (result string, err error) {

	// Create an embeddings client using the. Requires environment variable OPENAI_API_KEY to be set.
	llm, err := openai.New(
		openai.WithBaseURL(baseLink),
		openai.WithAPIVersion("v1"),
		openai.WithToken(api_token),
		openai.WithModel("wizard-uncensored-13b"),
		//openai.WithEmbeddingModel("text-embedding-ada-002"),
	)
	if err != nil {
		return "", err
	}

	result, err = chains.Run(
		context.Background(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, numOfResults),
		),
		question,
		chains.WithMaxTokens(2048),
	)
	if err != nil {
		return "", err
	}

	fmt.Println("====final answer====\n", result)

	return result, nil

}

func SemanticSearch(searchQuery string, maxResults int, store vectorstores.VectorStore) (searchResults []schema.Document, err error) {

	searchResults, err = store.SimilaritySearch(context.Background(), searchQuery, maxResults)

	if err != nil {
		return nil, err
	}

	fmt.Println("============== similarity search results ==============")

	for _, doc := range searchResults {
		fmt.Println("similarity search info -", doc.PageContent)
		fmt.Println("similarity search score -", doc.Score)
		fmt.Println("============================")

	}

	return searchResults, nil

}
