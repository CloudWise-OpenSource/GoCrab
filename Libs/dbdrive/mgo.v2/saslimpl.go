//+build sasl

package mgo

import (
	"github.com/CloudWise-OpenSource/GoCrab/Libs/dbdrive/mgo.v2/sasl"
)

func saslNew(cred Credential, host string) (saslStepper, error) {
	return sasl.New(cred.Username, cred.Password, cred.Mechanism, cred.Service, host)
}
