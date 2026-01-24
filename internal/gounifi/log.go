// Local addition to support vendored code that references 'log'
// This file is NOT from go-unifi and won't be overwritten by sync

package gounifi

import "github.com/sirupsen/logrus"

var log = logrus.New()
