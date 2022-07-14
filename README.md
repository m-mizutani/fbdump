# fbdump

Dump bunch of Firebase Auth user record

## Usage

### dump

```bash
$ export GOOGLE_APPLICATION_CREDENTIALS=path/to/your-service-account-key.json
$ fbdump --log-level dump
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

### features

- Keep dump state to `state.json` for retry

## License

Apache License version 2.0
