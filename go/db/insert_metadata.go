/* Functions for inserting metadata into Postgres tables */ 
package db 

import ( 
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertDocumentMetadata(ctx context.Context, p *pgxpool.Pool, doc Document) error { 

} // InsertDocumentMetadata