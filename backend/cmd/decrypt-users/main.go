// Command decrypt-users decrypts email and name fields in the users table
// that were previously stored encrypted. Run once after deploying the code
// change that removes email/name encryption.
//
// Usage: decrypt-users [--dry-run]
package main

import (
	"context"
	"flag"
	"log"

	"github.com/idtazkia/stmik-admission-api/internal/config"
	"github.com/idtazkia/stmik-admission-api/internal/model"
	"github.com/idtazkia/stmik-admission-api/internal/pkg/crypto"
	"github.com/joho/godotenv"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "only show what would be decrypted, don't update")
	flag.Parse()

	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := crypto.Init(cfg.Encryption.Key); err != nil {
		log.Fatalf("failed to initialize encryption: %v", err)
	}

	enc := crypto.Get()

	ctx := context.Background()
	if err := model.Connect(ctx, cfg.Database.DSN()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer model.Close()

	pool := model.Pool()

	rows, err := pool.Query(ctx, "SELECT id, email, name FROM users")
	if err != nil {
		log.Fatalf("failed to query users: %v", err)
	}
	defer rows.Close()

	type userRow struct {
		id, email, name string
	}
	var users []userRow
	for rows.Next() {
		var u userRow
		if err := rows.Scan(&u.id, &u.email, &u.name); err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		users = append(users, u)
	}
	rows.Close()

	for _, u := range users {
		emailDec, err := enc.DecryptDeterministic(u.email)
		if err != nil {
			emailDec = u.email // already plaintext
		}

		nameDec, err := enc.DecryptProbabilistic(u.name)
		if err != nil {
			nameDec = u.name // already plaintext
		}

		if emailDec == u.email && nameDec == u.name {
			log.Printf("[unchanged] id=%s email=%s name=%s", u.id, u.email, u.name)
			continue
		}

		log.Printf("[decrypt] id=%s email=%s name=%s", u.id, emailDec, nameDec)

		if *dryRun {
			continue
		}

		_, err = pool.Exec(ctx, "UPDATE users SET email = $1, name = $2 WHERE id = $3", emailDec, nameDec, u.id)
		if err != nil {
			log.Printf("[error] failed to update user %s: %v", u.id, err)
			continue
		}
		log.Printf("[updated] id=%s", u.id)
	}

	log.Println("Done.")
}
