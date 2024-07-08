package crypto

import (
	"crypto/elliptic"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUnmarshalWebKeyConfig(t *testing.T) {
	type args struct {
		data       []byte
		configType WebKeyConfigType
	}
	tests := []struct {
		name       string
		args       args
		wantConfig WebKeyConfig
		wantErr    error
	}{
		{
			name: "rsa",
			args: args{
				[]byte(`{"bits":2048, "hasher":0}`),
				WebKeyConfigTypeRSA,
			},
			wantConfig: &WebKeyRSAConfig{
				Bits:   RSABits2048,
				Hasher: RSAHasherSHA256,
			},
		},
		{
			name: "ecdsa",
			args: args{
				[]byte(`{"curve":0}`),
				WebKeyConfigTypeECDSA,
			},
			wantConfig: &WebKeyECDSAConfig{
				Curve: EllipticCurveP256,
			},
		},
		{
			name: "ed25519",
			args: args{
				[]byte(`{}`),
				WebKeyConfigTypeED25519,
			},
			wantConfig: &WebKeyED25519Config{},
		},
		{
			name: "unknown type error",
			args: args{
				[]byte(`{"curve":0}`),
				99,
			},
			wantErr: zerrors.ThrowInternal(nil, "CRYPT-Eig8ho", "Errors.Internal"),
		},
		{
			name: "unmarshal error",
			args: args{
				[]byte(`~~`),
				WebKeyConfigTypeED25519,
			},
			wantErr: zerrors.ThrowInternal(nil, "CRYPT-waeR0N", "Errors.Internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := UnmarshalWebKeyConfig(tt.args.data, tt.args.configType)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, gotConfig, tt.wantConfig)
		})
	}
}

func TestWebKeyRSAConfig_Alg(t *testing.T) {
	type fields struct {
		Hasher RSAHasher
	}
	tests := []struct {
		name   string
		fields fields
		want   jose.SignatureAlgorithm
	}{
		{
			name: "SHA256",
			fields: fields{
				Hasher: RSAHasherSHA256,
			},
			want: jose.RS256,
		},
		{
			name: "SHA384",
			fields: fields{
				Hasher: RSAHasherSHA384,
			},
			want: jose.RS384,
		},
		{
			name: "SHA512",
			fields: fields{
				Hasher: RSAHasherSHA512,
			},
			want: jose.RS512,
		},
		{
			name: "default",
			fields: fields{
				Hasher: 99,
			},
			want: jose.RS256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := WebKeyRSAConfig{
				Hasher: tt.fields.Hasher,
			}
			got := c.Alg()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWebKeyECDSAConfig_Alg(t *testing.T) {
	type fields struct {
		Curve EllipticCurve
	}
	tests := []struct {
		name   string
		fields fields
		want   jose.SignatureAlgorithm
	}{
		{
			name: "P256",
			fields: fields{
				Curve: EllipticCurveP256,
			},
			want: jose.ES256,
		},
		{
			name: "P384",
			fields: fields{
				Curve: EllipticCurveP384,
			},
			want: jose.ES384,
		},
		{
			name: "P512",
			fields: fields{
				Curve: EllipticCurveP512,
			},
			want: jose.ES512,
		},
		{
			name: "default",
			fields: fields{
				Curve: 99,
			},
			want: jose.ES256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := WebKeyECDSAConfig{
				Curve: tt.fields.Curve,
			}
			got := c.Alg()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWebKeyECDSAConfig_GetCurve(t *testing.T) {
	type fields struct {
		Curve EllipticCurve
	}
	tests := []struct {
		name   string
		fields fields
		want   elliptic.Curve
	}{
		{
			name:   "P256",
			fields: fields{EllipticCurveP256},
			want:   elliptic.P256(),
		},
		{
			name:   "P384",
			fields: fields{EllipticCurveP384},
			want:   elliptic.P384(),
		},
		{
			name:   "P512",
			fields: fields{EllipticCurveP512},
			want:   elliptic.P521(),
		},
		{
			name:   "default",
			fields: fields{99},
			want:   elliptic.P256(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := WebKeyECDSAConfig{
				Curve: tt.fields.Curve,
			}
			got := c.GetCurve()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_generateEncryptedWebKey(t *testing.T) {
	type args struct {
		keyID     string
		genConfig WebKeyConfig
	}
	tests := []struct {
		name          string
		args          args
		assertPrivate func(t *testing.T, got *jose.JSONWebKey)
		assertPublic  func(t *testing.T, got *jose.JSONWebKey)
		wantErr       error
	}{
		{
			name: "RSA",
			args: args{
				keyID: "keyID",
				genConfig: &WebKeyRSAConfig{
					Bits:   RSABits2048,
					Hasher: RSAHasherSHA256,
				},
			},
			assertPrivate: assertJSONWebKey("keyID", "RS256", "sig", false),
			assertPublic:  assertJSONWebKey("keyID", "RS256", "sig", true),
		},
		{
			name: "ECDSA",
			args: args{
				keyID: "keyID",
				genConfig: &WebKeyECDSAConfig{
					Curve: EllipticCurveP256,
				},
			},
			assertPrivate: assertJSONWebKey("keyID", "ES256", "sig", false),
			assertPublic:  assertJSONWebKey("keyID", "ES256", "sig", true),
		},
		{
			name: "ED25519",
			args: args{
				keyID:     "keyID",
				genConfig: &WebKeyED25519Config{},
			},
			assertPrivate: assertJSONWebKey("keyID", "EdDSA", "sig", false),
			assertPublic:  assertJSONWebKey("keyID", "EdDSA", "sig", true),
		},
		{
			name: "invalid config",
			args: args{
				keyID:     "keyID",
				genConfig: nil,
			},
			wantErr: zerrors.ThrowInternalf(nil, "CRYPT-aeW6x", "Errors.Internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrivate, gotPublic, err := generateWebKey(tt.args.keyID, tt.args.genConfig)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.assertPrivate != nil {
				tt.assertPrivate(t, gotPrivate)
			}
			if tt.assertPublic != nil {
				tt.assertPublic(t, gotPublic)
			}
		})
	}
}

func assertJSONWebKey(keyID, algorithm, use string, isPublic bool) func(t *testing.T, got *jose.JSONWebKey) {
	return func(t *testing.T, got *jose.JSONWebKey) {
		assert.NotNil(t, got)
		assert.NotNil(t, got.Key, "key")
		assert.Equal(t, keyID, got.KeyID, "keyID")
		assert.Equal(t, algorithm, got.Algorithm, "algorithm")
		assert.Equal(t, use, got.Use, "user")
		assert.Equal(t, isPublic, got.IsPublic(), "isPublic")
	}
}