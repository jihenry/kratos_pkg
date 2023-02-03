package sign

type rsaConfig struct {
	PrivateKey []byte
	PublicKey  []byte
	XOR        string
}

var (
	blockMap = map[string]*rsaConfig{
		"pagoda": &rsaConfig{
			PublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0laI3ZKjLIed3CGx87eRnzyzm
++9NPVrnFUvKH9Zx9AmoAdATXnMsaIYkqaNlgxMAHNuED3sqB5+fW4FYm9Fou54W
TKuyfMzqFoPRyB4Prkf//SwO+8CtHr2CfGHx1QduyJxbNSGYYhEStwjAKg8Fqhmc
hI9mIbdGBtgjmS3S4QIDAQAB
-----END PUBLIC KEY-----`),
			PrivateKey: []byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBALSVojdkqMsh53cI
bHzt5GfPLOb77009WucVS8of1nH0CagB0BNecyxohiSpo2WDEwAc24QPeyoHn59b
gVib0Wi7nhZMq7J8zOoWg9HIHg+uR//9LA77wK0evYJ8YfHVB27InFs1IZhiERK3
CMAqDwWqGZyEj2Yht0YG2COZLdLhAgMBAAECgYEAo5pE4oZxXccTmoWpM+2aVmod
tg5dGM8TQfPLPA1oDMkYznsF9eZF1d/EWAbQH7GGTz3VqmkUHlnVxVvzbUGNjx13
BqUUIPdDh8TuER20QEje1sk7IulWP5DgY+nQLV/65vuDDKUHeGolEQLVrIuaMJDV
aJpuncd8/JPP/1F1EbkCQQDaiIxxmiNbyFjAxPopOPNX5NEXZFu9yjUuIU82Q11d
Jm+B+C7F8FdUGml2xGIKY/qDx0PEKT0WK6LO9zBqThZXAkEA04uBO2M2hL6EI6l7
bhBdf7kDLVhfabH2Cs7VOi81jQfZgQv6763ayKkHeDzreKy5cIIP1WZxqVcI/uS4
2UtthwJAHwYzqg0P5//RWcydFy0WnuvFI2UEATWrxxjDfhiiMI88VV8+hKtSOoZl
Yo8OvBrlfb/URwzztyoKuwcswGrFkQJAMt1iT3NFkpl0kFaaFRbeRG2p8+dB2dou
fN7KqljbmXN/uuW0ipjU+FacMy8Ct1tgo0rCn98oCT2iLhe00pquVQJBAMW5UzZC
/WH6Nr/PxxlcWN2QxTbaxJnQp9GVfQqDSscQXCHJqqCWuYXUVETjKFO5cNH/qarL
kxJY6hZUvab5X5U=
-----END PRIVATE KEY-----`),
			XOR: "dnJAPa3ndA",
		},
	}
)
