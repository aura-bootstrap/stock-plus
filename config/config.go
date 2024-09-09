package config

import "github.com/bootstrap-library/stock-plus/crypto"

var (
	TelegramAPITokenList = []string{
		crypto.Decrypt("Iv4F7UvO0e6W74Uj5GtyTVV5OiJJlvalJPvR0xdGjQZ3Phqyb3ILnQ+hQrkLk/p+"),
		crypto.Decrypt("iz45kxVlDBglt1J710ujM1MzMh3LFQj5Qvfc/4TpGO19tqkB1fgLm9pAalIn4/nk"),
		crypto.Decrypt("K5BrfH0/ONr2+o1cuWikzVSxccL1GJn2/v6q7Wu8lwh68H4MZfqSxzuNj+413YUw"),
	}
	TelegramFirstName = crypto.Decrypt("YfS9C5G3V9loLwij1ZSBpwdTimK1Nz5chM1FjnVSCIA=")
	TelegramProxy     = crypto.Decrypt("s+C0xb5Oi+Fh/LqiIuMpzibqZUC2ROGizhJhSC5IiXFGqSDkvkg825pv/7JSccMqgVkL2Sx/mKGD9Pb091GiqA==")
)
