package handler

import "fmt"

const (
	tagSuffix = "99999999999999999999999"

	// TagPreorderTrytes is the trytes for pre-order tag
	TagPreorderTrytes = "ZBYB" + tagSuffix
	// TagRegisterTrytes is the trytes for register tag
	TagRegisterTrytes = "ACQB" + tagSuffix
	// TagRenewTrytes is the trytes for renew tag
	TagRenewTrytes = "ACXB" + tagSuffix
	// TagUpdateTrytes is the trytes for update tag
	TagUpdateTrytes = "DCNB" + tagSuffix
	// TagTransferTrytes is the trytes for transfer tag
	TagTransferTrytes = "CCPB" + tagSuffix
	// TagRevokeTrytes is the trytes for revoke tag
	TagRevokeTrytes = "ACEC" + tagSuffix

	// TagPreorder is the tag for pre-ordering the domain
	TagPreorder = "PO"
	// TagRegister is the tag for registering the domain
	TagRegister = "RG"
	// TagRenew is the tag for renewing the domain
	TagRenew = "RN"
	// TagUpdate is the tag for updating the domain
	TagUpdate = "UD"
	// TagTransfer is the tag for transfering the domain
	TagTransfer = "TF"
	// TagRevoke is the tag for revoking the domain
	TagRevoke = "RV"
)

// NewTxHandler is used for handling the incoming new transaction by ZMQ listener
func NewTxHandler(txContent []string) error {
	var err error
	switch tag := txContent[12]; tag {
	case TagPreorderTrytes:
		fmt.Println("Preorder transaction checking")
		err = PreorderHandler(&txContent)
	case TagRegisterTrytes:
		fmt.Println("Register transaction checking")
		err = RegisterHandler(&txContent)
	case TagRenewTrytes:
		fmt.Println("Renew transaction checking")
		err = RenewHandler()
	case TagUpdateTrytes:
		fmt.Println("Update transaction checking")
		err = UpdateHandler()
	case TagTransferTrytes:
		fmt.Println("Transfer transaction checking")
		err = TransferHandler()
	case TagRevokeTrytes:
		fmt.Println("Revoke transaction checking")
		err = RevokeHandler()
	default:
		fmt.Println("Other transaction, skipped")
		err = nil
	}

	return err
}

// ConfirmedTxHandler is used for handling the incoming confirmed transaction by ZMQ listener
func ConfirmedTxHandler() {}

// CheckConfirmedTxMsg is used for checking the message of a confirmed transaction
func CheckConfirmedTxMsg() {}

// PreorderHandler is used for handling the new transaction for pre-ordering the domain name
func PreorderHandler(txContent *[]string) error {
	return nil
}

// RegisterHandler is used for handling the new transaction for registering the domain name
func RegisterHandler(txContent *[]string) error {
	return nil
}

// RenewHandler is used for handling the new transaction for renewing the domain name
func RenewHandler() error {
	return nil
}

// UpdateHandler is used for handling the new transaction for updating the domain name
func UpdateHandler() error {
	return nil
}

// TransferHandler is used for handling the new transaction for transfering the domain name
func TransferHandler() error {
	return nil
}

// RevokeHandler is used for handling the new transaction for revoking the domain name
func RevokeHandler() error {
	return nil
}
