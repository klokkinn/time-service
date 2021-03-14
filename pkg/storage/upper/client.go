package upper

import (
	"fmt"

	"github.com/klokkinn/time-service/pkg/core"
	"github.com/klokkinn/time-service/pkg/core/models"

	"github.com/imdario/mergo"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

func (c *client) Add(item models.Entry) (err error) {
	collection := c.Session.Collection(entryTable)

	_, err = collection.Insert(&item)
	if err != nil {
		return fmt.Errorf("error inserting item: %w", err)
	}

	return nil
}

func (c *client) Get(id string) (result models.Entry, err error) {
	collection := c.Session.Collection(entryTable)

	condition := db.Cond{"id": id}

	results := collection.Find(condition)

	exists, err := results.Exists()
	if err != nil {
		return result, fmt.Errorf("fetching item: %w", err)
	}

	if !exists {
		return result, core.StorageErrorNotFound
	}

	err = results.One(&result)
	if err != nil {
		return result, fmt.Errorf("finding item: %w", err)
	}

	return result, nil
}

func (c *client) GetAll(filter core.StorageFilter) (entries []models.Entry, err error) {
	condition := filterToDBCond(filter)

	collection := c.Session.Collection(entryTable)

	err = collection.Find(condition).All(&entries)
	if err != nil {
		return nil, fmt.Errorf("fetching items: %w", err)
	}

	return entries, nil
}

func (c *client) Update(update models.Entry) (updatedEntry models.Entry, err error) {
	collection := c.Session.Collection(entryTable)

	var originalGoal models.Entry

	condition := db.Cond{"id": update.Id}

	result := collection.Find(condition)

	exists, err := result.Exists()
	if err != nil {
		return updatedEntry, fmt.Errorf("fetching item: %w", err)
	}

	if !exists {
		return updatedEntry, core.StorageErrorNotFound
	}

	err = result.One(&originalGoal)
	if err != nil {
		return updatedEntry, fmt.Errorf("fetching original item: %w", err)
	}

	err = mergo.Merge(&originalGoal, update, mergo.WithOverride)
	if err != nil {
		return updatedEntry, fmt.Errorf("merging updated with original item: %w", err)
	}

	err = collection.UpdateReturning(&originalGoal)
	if err != nil {
		return updatedEntry, fmt.Errorf("updating item: %w", err)
	}

	return originalGoal, nil
}

func (c *client) Delete(id string) (err error) {
	collection := c.Session.Collection(entryTable)

	condition := db.Cond{"id": id}

	result := collection.Find(condition)

	exists, err := result.Exists()
	if err != nil {
		return fmt.Errorf("fetching item: %w", err)
	}

	if !exists {
		return core.StorageErrorNotFound
	}

	err = result.Delete()
	if err != nil {
		return fmt.Errorf("deleting item with ID %s: %w", id, err)
	}

	return nil
}

func (c *client) Open() error {
	sess, err := postgresql.Open(c.connectionURL)
	if err != nil {
		return fmt.Errorf("connecting to Postgres: %w", err)
	}

	c.Session = sess

	return c.setup()
}

func (c *client) Close() error {
	return c.Session.Close()
}

func (c *client) setup() error {
	sql := c.Session.SQL()

	_, err := sql.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id text primary key,
		startmillis text not null,
		endmillis text
	)`, entryTable))
	if err != nil {
		return fmt.Errorf("creating tables: %w", err)
	}

	return nil
}

func filterToDBCond(filter core.StorageFilter) (condition db.Cond) {
	condition = db.Cond{}

	if filter.Author != nil {
		condition["author"] = *filter.Author
	}

	return condition
}
