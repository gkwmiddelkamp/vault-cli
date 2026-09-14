package model

// Token types understood by the Previder Vault API. Only ReadWrite and
// ReadOnly may interact with secrets: ReadWrite can list, get and decrypt
// them, while ReadOnly can decrypt a secret whose id or name is already
// known. The admin types manage tokens rather than secrets.
const (
	TokenTypeReadOnly         = "ReadOnly"
	TokenTypeReadWrite        = "ReadWrite"
	TokenTypeEnvironmentAdmin = "EnvironmentAdmin"
	TokenTypeMasterAdmin      = "MasterAdmin"
)

type Token struct {
	Id            string `json:"id"`
	Description   string `json:"description,omitempty"`
	EnvironmentId string `json:"environmentId,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
	TokenType     string `json:"tokenType,omitempty"`
}
