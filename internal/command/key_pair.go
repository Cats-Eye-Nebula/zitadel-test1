package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/keypair"
)

const (
	oidcUser = "OIDC"
)

func (c *Commands) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	ctx = setOIDCCtx(ctx)
	privateCrypto, publicCrypto, err := crypto.GenerateEncryptedKeyPair(c.keySize, c.keyAlgorithm)
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	privateKeyExp := time.Now().UTC().Add(c.privateKeyLifetime)
	publicKeyExp := time.Now().UTC().Add(c.publicKeyLifetime)

	//TODO: InstanceID not available here?
	keyPairWriteModel := NewKeyPairWriteModel(keyID, "system") //TODO: change with multi issuer
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSigning,
		algorithm,
		privateCrypto, publicCrypto,
		privateKeyExp, publicKeyExp))
	return err
}

func setOIDCCtx(ctx context.Context) context.Context {
	//TODO: InstanceID not available here?
	return authz.SetCtxData(ctx, authz.CtxData{UserID: oidcUser, OrgID: authz.GetInstance(ctx).InstanceID()})
}
