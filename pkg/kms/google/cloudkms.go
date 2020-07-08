package google

import (
	"context"
	"fmt"
	"github.com/square/go-jose"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type Client struct {
	Project    string
	Location   string
	KeyRing    string
	KeyName    string
	KeyVersion int
}

func (c *Client) keyFullName() string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d",
		c.Project,
		c.Location,
		c.KeyRing,
		c.KeyName,
		c.KeyVersion)
}

func (c *Client) DecryptKey(encryptedKey []byte, header jose.Header) ([]byte, error) {
	// TODO (immutableT) Examine the header to see if the payload was encrypted via the parameters
	// supported by the CryptoKey.

	return c.decryptAsymmetric(encryptedKey)
}

func (c *Client) decryptAsymmetric(ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %v", err)
	}

	req := &kmspb.AsymmetricDecryptRequest{
		Name:       c.keyFullName(),
		Ciphertext: ciphertext,
	}

	result, err := client.AsymmetricDecrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ciphertext: %v", err)
	}
	return result.Plaintext, nil
}
