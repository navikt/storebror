# storebror

Storebror passer på å holde Git-repositories oppdatert med korrekt konfigurasjon.

*WORK IN PROGRESS*

Akkurat nå sjekkes det om nais-manifester i et repository inneholder en seksjon
med `prometheus` som er skrudd på. Hvis ikke, setter Storebror den inn i
kildekoden, og gjør en `git commit`.

## Utviklingsmiljø

* Golang >= 1.10
* [Dep](https://github.com/golang/dep)

```
dep ensure
go build
./storebror https://github.com/path/to/repository
# eller:
./storebror file:///path/to/repository
```

## Docker-container

```
docker build -t navikt/storebror .  # eller "make"
docker run --rm -it navikt/storebror https://github.com/path/to/repository
```
