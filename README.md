### Why
If you have an [MFA-enabled](https://aws.amazon.com/iam/details/mfa/) account on Amazon AWS, you need to refresh the token periodically, in order to use [aws cli toolkit](https://aws.amazon.com/cli/). 

The sequence of actions is:

* using the primary AWS account, request the [list of MFA devices](http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListMFADevices.html) configured for this account
* issue an STS request to [get the session token](http://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html)
* update the `~/.aws/credentials` file with the received access key, secret key and session token for the given profile

This simple flow is implemented as Go utility, that only updates the existing profile in the `~/.aws/credentials` with the access/secret/session tokens.

There is another utility [awsmfa](https://github.com/dcoker/awsmfa/) with extended functionality for AWS key management / rotation.

### How

```
Usage of ./go-aws-mfa:
  -d string
        MFA-enabled profile
  -s string
        Source (primary) profile
```

where

* `-s` specifies the IAM role that has an [MFA device configured](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html)
* `-d` specifies the target profile to add/replace the credentials to.

#### Example

`./go-aws-mfa -s user1 -d user1-mfa` will ask for the token code for MFA device configured for `user1`. Then the temporary credentials will be stored for `user1-mfa`.
In order to use that temporary account with `awscli`, you need to set the `AWS_PROFILE` environment variable to `user1-mfa` and then invoke `aws` command normally, for example:

```
AWS_PROFILE=user1-mfa aws s3 ls s3://bucket-user1/
```
