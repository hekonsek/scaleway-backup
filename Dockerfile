FROM fedora:28

ADD scaleway-backup /usr/bin/scaleway-backup

CMD /usr/bin/scaleway-backup