package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"rest/config"
	"rest/internal/repo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDeps struct {
	Conf *config.Postgres
}

type postgres struct {
	pool *pgxpool.Pool
}

type PostgresI interface {
	GetWalletByID(ctx context.Context, walletID string) (*models.Wallet, error)
	MakeTransaction(ctx context.Context, walletID, operation string, amount float64) error
}

func New(ctx context.Context, dep PostgresDeps) (PostgresI, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		dep.Conf.User, dep.Conf.Password, dep.Conf.Host, dep.Conf.Port, dep.Conf.Database,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25 // Максимум соединений
	config.MinConns = 5  // Минимум соединений в пуле
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Подключение
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &postgres{
		pool: pool,
	}, nil
}

func (p *postgres) GetWalletByID(ctx context.Context, walletID string) (*models.Wallet, error) {
	var wallet models.Wallet
	query := `
		SELECT id, balance, status, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	err := p.pool.QueryRow(ctx, query, walletID).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.Status,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

func (db *postgres) MakeTransaction(ctx context.Context, walletID, operation string, amount float64) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
	var currentBalance float64
	var status string
	walletQuery := `
			SELECT balance, status
			FROM wallets
			WHERE id = $1
			FOR UPDATE
	`

	tx, err := db.pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, walletQuery, walletID).Scan(&currentBalance, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	if status != "ACTIVE" {
		return fmt.Errorf("wallet is not active (current status: %s)", status)
	}

	var newBalance float64

	switch models.OperationType(operation) {
	case models.OperationDeposit:
		newBalance = currentBalance + amount
	case models.OperationWithdraw:
		if currentBalance < amount {
			return fmt.Errorf("insufficient funds: available %.2f, requested %.2f",
				currentBalance, amount)
		}
		newBalance = currentBalance - amount
	default:
		return fmt.Errorf("invalid operation type: %s", operation)
	}

	// 4. Обновляем баланс кошелька
	updateQuery := `
			UPDATE wallets
			SET balance = $1
			WHERE id = $2
		`

	updateResult, err := tx.Exec(ctx, updateQuery, newBalance, walletID)
	if err != nil {
		// Проверяем constraint violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" { // check_violation
			return fmt.Errorf("balance cannot be negative")
		}
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	if updateResult.RowsAffected() == 0 {
		return fmt.Errorf("wallet update failed")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
