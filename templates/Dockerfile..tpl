FROM scratch
COPY ./{{.ProjectFolder }} /
EXPOSE 443
ENTRYPOINT ["./{{.ProjectFolder }}"]
