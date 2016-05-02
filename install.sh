#/bin/bash

echo "install Black Duck Software Scanner"

cp -f /blackduck /host/etc/atomic.d/

/atomic_scanner
