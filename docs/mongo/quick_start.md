## [Enable authorization](https://docs.mongodb.com/manual/core/authorization/)

MongoDB does not enable access control by default. You can enable authorization using the [`--auth`](https://docs.mongodb.com/manual/reference/program/mongod/#std-option-mongod.--auth) or the [`security.authorization`](https://docs.mongodb.com/manual/reference/configuration-options/#mongodb-setting-security.authorization) setting. Enabling [internal authentication](https://docs.mongodb.com/manual/core/security-internal-authentication/) also enables client authorization.

### Step1:Create user for root

```mongo
db.createUser(
    {
        user: "root",
        pwd: "root",
        roles: ["root"]
    }
)
```

### Step2:Modify config file

If your use MacOS the config file position is `/usr/local/etc/mongod.conf`, Linux is `/etc/mongod.conf`

```conf
security:
  authorization: enabled
```

### Step3:Reset mongod server

```shell
$ mongod --config /usr/local/etc/mongod.conf
```

### Step4:User root login mongosh

```shell
$ mongosh -u root -p root
```