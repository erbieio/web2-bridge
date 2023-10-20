package jobs

const (
	BatchCount = 100
)

func Do() {
	go HandleMintReceipt()
	go HandleTransferReceipt()
}
