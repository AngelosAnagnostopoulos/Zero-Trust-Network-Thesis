go build
sudo ./pep_server -host "pep_server" -srvcert "./pep_server.crt"  -srvkey "./pep_server.key" -cacert "./ThesisCA.crt" -port 443 -certopt 4