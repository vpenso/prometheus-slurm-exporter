rm -f go.mod
go get -v -u github.com/prometheus/client_golang/prometheus
go get -v -u github.com/sirupsen/logrus
go get -v -u gopkg.in/alecthomas/kingpin.v2
go get -v -u github.com/m3db/prometheus_common/log
mkdir -p ${GOPATH}/src/github.com/prometheus/common/log/
cp -r ${GOPATH}/src/github.com/m3db/prometheus_common/log/* ${GOPATH}/src/github.com/prometheus/common/log/