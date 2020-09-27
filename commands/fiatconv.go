package commands

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/urfave/cli"

	"github.com/fiatconv/integration/exchangeratesapi"
	"github.com/fiatconv/internal/fiatconv"

)

func NewFiatconvCmd() cli.Command {
	return cli.Command{
		Name: "fiatconv",
		Action: runConverter,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "base",
			},
			cli.StringFlag{
				Name: "quote",
			},
			cli.StringFlag{
				Name: "amount",
			},

		},
	}
}

func runConverter(c *cli.Context) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	exchangeratesapiAPI := exchangeratesapi.New("https://api.exchangeratesapi.io/")
	exchangeratesapiClient := fiatconv.NewExchangeRatesAPIClient(exchangeratesapiAPI)
	fiatconvService := fiatconv.NewService(exchangeratesapiClient)
	fiatconvTransport := fiatconv.NewTransport(logger, fiatconvService)

	fiatconvTransport.ConvertFromBaseToQuote(c.String("base"), c.String("quote"), c.String("amount"))
}