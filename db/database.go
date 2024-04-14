package db

import (
	"context"
	"fmt"
	"os"

	_ "encoding/json"

	"github.com/jackc/pgx/v4"
	"github.com/rassulmagauin/webscraper/gpt"
)

// dbDriver encapsulates database operations
type DbDriver struct {
	conn *pgx.Conn
}

// NewDBDriver creates a new instance of dbDriver
func NewDBDriver() *DbDriver {
	return &DbDriver{}
}

// Connect establishes a connection to the PostgreSQL database and sets the connection handle in dbDriver
func (d *DbDriver) Connect() error {
	// Parse the connection string using pgx.ParseConfig
	dbURL := os.Getenv("DATABASE_URL")
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse DATABASE_URL: %w", err)
	}

	// Connect to the PostgreSQL server
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	fmt.Println("Connected successfully!")
	d.conn = conn
	return nil
}

// Close closes the database connection
func (d *DbDriver) Close() {
	if d.conn != nil {
		d.conn.Close(context.Background())
	}
}

func (d *DbDriver) UpdateOrCreateOffers(ctx context.Context, offers []gpt.Offer) error {
	if d.conn == nil {
		return fmt.Errorf("connection is not established")
	}

	for _, offer := range offers {
		_, err := d.conn.Exec(ctx, `
		INSERT INTO offers (cashback, condition, expiry, category, restrictions, card_type, bank_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (card_type, bank_id, condition) DO UPDATE
		SET cashback = EXCLUDED.cashback,
			condition = EXCLUDED.condition,
			expiry = EXCLUDED.expiry,
			category = EXCLUDED.category,
			restrictions = EXCLUDED.restrictions;
	`, offer.Cashback, offer.Condition, offer.Expiry, offer.Category, offer.Restrictions, offer.CardType, offer.BankID)

		if err != nil {
			// Optionally, log the error but continue processing other offers
			// This approach means one failure does not stop the rest of the process
			fmt.Printf("Error updating or creating offer for bank ID %d and card type %s: %v\n", offer.BankID, offer.CardType, err)
		}
	}

	return nil
}
