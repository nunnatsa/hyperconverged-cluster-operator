module github.com/kubevirt/hyperconverged-cluster-operator/tools/release-notes

go 1.21

toolchain go1.21.7

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/joho/godotenv v1.4.0
	github.com/kubevirt/hyperconverged-cluster-operator/tools/release-notes/git v0.0.0-00010101000000-000000000000
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230828082145-3c4c8a2d2371 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-git/go-git/v5 v5.11.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-github/v32 v32.1.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/skeema/knownhosts v1.2.1 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

replace github.com/kubevirt/hyperconverged-cluster-operator/tools/release-notes/git => ./git

// FIX: Denial of service in golang.org/x/text/language
replace golang.org/x/text => golang.org/x/text v0.7.0

// FIX: Uncontrolled Resource Consumption
replace golang.org/x/net => golang.org/x/net v0.9.0

// FIX: Use of a Broken or Risky Cryptographic Algorithm in golang.org/x/crypto/ssh
replace golang.org/x/crypto => golang.org/x/crypto v0.8.0
