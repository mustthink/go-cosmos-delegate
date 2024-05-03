package parser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	tx2 "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"

	"github.com/mustthink/go-cosmos-delegate/internal/models"
)

type (
	Parser interface {
		Parse(ctx context.Context) ([]models.Transaction, error)
	}

	parser struct {
		client cmtservice.ServiceClient
	}
)

func New(conn *grpc.ClientConn) Parser {
	return &parser{
		client: cmtservice.NewServiceClient(conn),
	}
}

func (p *parser) Parse(ctx context.Context) ([]models.Transaction, error) {
	now := time.Now()
	response, err := p.client.GetLatestBlock(ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, fmt.Errorf("couldn't get latest block: %w", err)
	}

	header := response.SdkBlock.GetHeader()
	blockID := (&header).GetHeight()
	var transactions = make([]models.Transaction, 0, len(response.SdkBlock.GetData().Txs))
	for _, txBytes := range response.SdkBlock.GetData().Txs {
		var tx tx2.Tx
		err := proto.Unmarshal(txBytes, &tx)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
		}

		hash := sha256.Sum256(txBytes)
		transactionID := hex.EncodeToString(hash[:])

		msgs := tx.Body.GetMessages()
		delegateMessages, err := parseDelegateMessages(msgs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse delegate messages: %w", err)
		}

		transactions = append(transactions, models.Transaction{
			ExternalID: transactionID,
			BlockID:    blockID,
			Timestamp:  now,
			Messages:   delegateMessages,
		})
	}

	return transactions, nil
}

const delegateMessage = "/cosmos.staking.v1beta1.MsgDelegate"

func parseDelegateMessages(messages []*codecTypes.Any) ([]models.DelegateMessage, error) {
	var delegateMessages []models.DelegateMessage
	for _, msg := range messages {
		if msg.TypeUrl == delegateMessage {
			delegate := &types.MsgDelegate{}
			err := proto.Unmarshal(msg.Value, delegate)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal delegate message: %w", err)
			}

			delegateMessages = append(delegateMessages, models.DelegateMessage{
				DelegatorAddress: delegate.DelegatorAddress,
				ValidatorAddress: delegate.ValidatorAddress,
				Amount:           delegate.Amount.Amount.Int64(),
				Currency:         parseCurrency(delegate.Amount.String()),
			})
		}
	}
	return delegateMessages, nil
}

func parseCurrency(s string) string {
	re := regexp.MustCompile("[a-zA-Z]+")
	return re.FindString(s)
}
