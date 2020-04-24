package main

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/eoscanada/eos-go"
    "github.com/kelseyhightower/envconfig"
    "github.com/rs/zerolog/log"
    "io"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
)

func readOffset(r io.Reader) (uint64, error) {
    log.Debug().Msg("reading offset")
    var offset uint64
    _, err := fmt.Fscan(r, &offset)
    return offset, err
}

func writeOffset(w io.Writer, offset uint64) error {
    log.Debug().Msgf("writing offset, value: %v", offset)
    _, err := fmt.Fprint(w, offset)
    return err
}

func readWIF(filename string) string {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Panic().Msg(err.Error())
    }
    wif := strings.TrimSpace(strings.TrimSuffix(string(content), "\n"))
    return wif
}

func readRsa(filename string) (*rsa.PrivateKey, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    block, _ := pem.Decode(data)
    key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return key, err
}

func readConfigFile(cfg *Config, path string) {
    _, err := toml.DecodeFile(path, cfg)
    if err != nil {
        log.Panic().Msg(err.Error())
    }
}

func readEnv(cfg *Config) {
    err := envconfig.Process("", cfg)
    if err != nil {
        log.Panic().Msg(err.Error())
    }
}

func getConfigPath(envVar, defaultValue string) string {
    cfgPath, isSet := os.LookupEnv(envVar)
    if isSet {
        return cfgPath
    }
    return defaultValue
}

func getAddr(port int) string {
    return ":" + strconv.Itoa(port)
}

func rsaSign(digest eos.Checksum256, key *rsa.PrivateKey) (string, error) {
    sign, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
    if err != nil {
        return "", err
    }

    // contract requires base64 string
    return base64.StdEncoding.EncodeToString(sign), nil
}
