package options

import (
	"github.com/sirupsen/logrus"
	"test-billing/commons/config"
	"test-billing/pkg/utils"
)

type Options struct {
	Config     *config.Conf
	Logger     *logrus.Logger
	DBPostgres *utils.DB
}
