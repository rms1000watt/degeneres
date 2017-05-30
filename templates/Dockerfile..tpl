FROM scratch
COPY {{.CertsPath}} /certs
COPY ./{{.ProjectFolder }} /
EXPOSE 443
ENTRYPOINT ["./{{.ProjectFolder }}"]
