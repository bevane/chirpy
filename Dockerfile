FROM debian:stable-slim

COPY chirpy /bin/chirpy
COPY assets /bin/assets
COPY index.html /bin/index.html

ENV PORT=8080
WORKDIR "/bin"
CMD ["./chirpy"]
