FROM amd64/alpine
EXPOSE 3001
WORKDIR /app
COPY mplay2-go .
CMD ["/app/mplay2-go"]
