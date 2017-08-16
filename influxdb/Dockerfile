FROM influxdb:1.2

# Shamelessly stolen

MAINTAINER ruslan.ledesma@digitalonus.com

# Create databases
COPY start-and-initialize.sh /
CMD ["/start-and-initialize.sh"]