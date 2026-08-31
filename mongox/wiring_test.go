package mongox_test

import (
	"testing"

	"github.com/hugo2lee/gotox/mongox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDatabaseWiring(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	require.NoError(t, err)
	database := client.Database("test")

	wrapped := mongox.NewWithDatabase(database, nil)
	assert.Same(t, database, wrapped.DB())
}
