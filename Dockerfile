FROM scratch
ADD main /
ADD configs/config.yml /configs/config.yml
CMD ["/main"]