FROM scratch
ENTRYPOINT ["/steampipe-plugin-confluence.plugin"]
COPY steampipe-plugin-confluence.plugin /
