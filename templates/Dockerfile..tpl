FROM scratch
COPY ./{{.ProjectFolder }} /
ENTRYPOINT ["./{{.ProjectFolder }}"]
