package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"log"
	"os"
	"os/user"
)

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	srcF := flag.String("s", "", "Source (primary) profile")
	dstF := flag.String("d", "", "MFA-enabled profile")

	flag.Parse()

	if *srcF == "" || *dstF == "" {
		flag.Usage()
		os.Exit(1)
	}

	conf := &aws.Config{
		Credentials: credentials.NewSharedCredentials("", *srcF),
	}

	sess, err := session.NewSession(conf)

	fatalErr(err)

	_iam := iam.New(sess)

	devices, err := _iam.ListMFADevices(&iam.ListMFADevicesInput{})

	fatalErr(err)

	if len(devices.MFADevices) == 0 {
		log.Fatal("No MFA devices configured")
	}

	sn := devices.MFADevices[0].SerialNumber

	fmt.Printf("Using device %1s\n", *sn)

	_sts := sts.New(sess)

	fmt.Printf("Enter MFA code: ")

	r := bufio.NewReader(os.Stdin)
	code, _, err := r.ReadLine()

	fatalErr(err)

	codeStr := string(code)

	res, err := _sts.GetSessionToken(&sts.GetSessionTokenInput{
		TokenCode:    &codeStr,
		SerialNumber: sn,
	})

	fatalErr(err)

	usr, err := user.Current()

	fatalErr(err)

	filePath := usr.HomeDir + "/.aws/credentials"

	credFile, err := ini.Load(filePath)

	fatalErr(err)

	sect, err := credFile.NewSection(*dstF)

	fatalErr(err)

	sect.NewKey("aws_access_key_id", *res.Credentials.AccessKeyId)
	sect.NewKey("aws_secret_access_key", *res.Credentials.SecretAccessKey)
	sect.NewKey("aws_session_token", *res.Credentials.SessionToken)

	credFile.SaveTo(filePath)

	fatalErr(err)

	fmt.Printf("Access token updated for %1s\n", *dstF)

}
