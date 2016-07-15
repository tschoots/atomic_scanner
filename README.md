# atomic_scanner


# Intro

This is a scanner ment for the the atomic platform.


# Usage

Before we can use the scanner we first have to configure a data container that contains the Black Duck hub server configuration so the scanner can connect.
This is done by the following docker command:

$docker run -ti -v /scanner -v /conf --name eng-hub blackducksoftware/atomic_scanner

$docker run -ti --rm -h $(hostname) --volumes-from eng-hub -v /etc/localtime:/etc/localtime -v $(pwd)/scanin:/scanin -v $(pwd)/scanout:/scanout ton


# Release

When a new atomic scanner has to be released the following steps should be taken


1.  move to the top directory and run the following command
`./docker_golang_build.sh` 
2.  push the created image with the follwing command
`docker push blackducksoftware/atomic` 
