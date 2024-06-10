package embedderTemplate

import (
	"context"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

// api_token - localai local key
// dblink - postgres link
// base_link - link to machine with localai running
func GetVectorStore(api_token string, db_link string, base_link string) (vectorstores.VectorStore, error) {
	llm, err := openai.New(
		openai.WithBaseURL(base_link),
		openai.WithAPIVersion("v1"),
		//openai.WithModel("wizard-uncensored-13b"),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
		openai.WithToken(api_token),
	)
	if err != nil {
		fmt.Println(err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		fmt.Println(err)
	}

	store, err := pgvector.New(
		context.Background(),
		//pgvector.WithPreDeleteCollection(true),
		pgvector.WithConnectionURL(db_link),
		pgvector.WithEmbedder(e),
	)
	if err != nil {
		fmt.Println(err)
	}
	return store, nil
}

func StringToStore(bio string, username string, store vectorstores.VectorStore) {

	now := time.Now()
	formattedDate := now.Format("02/01/2006")

	doc := []schema.Document{
		{
			PageContent: bio,
			Metadata: map[string]interface{}{
				"username": username,
				"date":     formattedDate,
			},
		},
	}
	_, err := store.AddDocuments(context.Background(), doc)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data for " + username + " successfully loaded into vector store💾")
	}
}

func LocationToStore(location string, mainLang string, bio string, username string, store vectorstores.VectorStore) {
	//update
	now := time.Now()
	formattedDate := now.Format("02/01/2006")

	doc := []schema.Document{
		{
			PageContent: bio,
			Metadata: map[string]interface{}{
				"username": username,
				"date":     formattedDate,
				"language": mainLang,
				"location": location,
			},
		},
	}
	_, err := store.AddDocuments(context.Background(), doc)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data for " + username + " successfully loaded into vector store💾")
	}
}
