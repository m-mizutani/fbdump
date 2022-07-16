# fbdump

`fbdump` is a tool to dump a bunch of Firebase Auth user records. Firebase provides [firebase-tools](https://www.npmjs.com/package/firebase-tools) that can export firebase users. However the tool is not suitable to export a lot of users about than 1 million. Because,

- Can not retry if aborted before completed
- Can output only one large file that is breakable by abort
- Will crash after export about tens of millions of users by stack overflow of recursive call.

`fbdump` keeps dumping state for retry when aborted and can start again from aborted point. Additionally, the tool splits exported user records to multiple files for avoiding to break a file by abort.

## Install

```bash
$ go install github.com/m-mizutani/fbdump@latest
```

## Usage

### dump

```bash
$ export GOOGLE_APPLICATION_CREDENTIALS=path/to/your-service-account-key.json
$ fbdump --log-level debug dump
12:21:37.692 [debug] saved user records
"path" => "repo/0/0/000xxxxxxxxxxxxxxx.json"
...
```

### load

```bash
$ fbdump load
{"providerId":"firebase","rawId":"xxxxxxxx","CustomClaims":null,"Disabled":false,"EmailVerified":false,"ProviderUserInfo":null,"TokensValidAfterMillis":0,"UserMetadata":{"CreationTimestamp":1234567890,"LastLogInTimestamp":1234567890,"LastRefreshTimestamp":1234567890},"TenantID":"","MultiFactor":{"EnrolledFactors":null},"PasswordHash":"","PasswordSalt":""}
...
```

## License

Apache License version 2.0
