# atomic_scanner


# Intro

This is a scanner ment for the the atomic platform.


# Usage

Before we can use the scanner we first have to configure a data container that contains the Black Duck hub server configuration so the scanner can connect.
This is done by the following docker command:

$docker run -ti -v /scanner -v /conf --name eng-hub blackducksoftware/atomic_scanner

$docker run -ti --rm -h $(hostname) --volumes-from eng-hub -v /etc/localtime:/etc/localtime -v $(pwd)/scanin:/scanin -v $(pwd)/scanout:/scanout ton

