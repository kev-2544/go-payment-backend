package repository

import (
    "context"
    "errors"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/kev-2544/go-payment-backend/internal/model"
)

type postgresPaymentRepository struct {
    db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository {
    return &postgresPaymentRepository{
        db: db,
    }
}

func (r *postgresPaymentRepository) CreateTransaction(
    ctx context.Context,
    tx *model.Transaction,
) error {

    query := `
        INSERT INTO transactions (
            id,
            user_id,
            amount,
            fee,
            currency,
            status,
            provider_reference
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

    _, err := r.db.Exec(
        ctx,
        query,
        tx.ID,
        tx.UserID,
        tx.Amount,
        tx.Fee,
        tx.Currency,
        tx.Status,
        tx.ProviderReference,
    )

    return err
}

func (r *postgresPaymentRepository) GetTransaction(
    ctx context.Context,
    id string,
) (*model.Transaction, error) {

    tx := &model.Transaction{}

    query := `
        SELECT id, user_id, amount, fee, currency, status,
               COALESCE(provider_reference, '')
        FROM transactions
        WHERE id = $1
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &tx.ID,
        &tx.UserID,
        &tx.Amount,
        &tx.Fee,
        &tx.Currency,
        &tx.Status,
        &tx.ProviderReference,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, err
    }

    return tx, nil
}

func (r *postgresPaymentRepository) UpdateProviderReference(
    ctx context.Context,
    transactionID string,
    providerReference string,
) error {

    _, err := r.db.Exec(
        ctx,
        `UPDATE transactions
         SET provider_reference = $1
         WHERE id = $2`,
        providerReference,
        transactionID,
    )

    return err
}
