#a simple reverse proxy

run:

    docker network create mynetwork
    docker run -d --name db --net mynetwork klaemo/couchdb:1.6.1
    docker run --rm -it -p 5984:8888 --net mynetwork gesellix/reverse-proxy -port :8888 -target http://db:5984

run (legacy):

    docker run --rm -it -p 5984:8888 --link couchdb:db gesellix/reverse-proxy -port :8888 -target http://db:5984
