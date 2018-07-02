FROM fedora:28

ADD bin/scaleway-backup /usr/bin/scaleway-backup

CMD /usr/bin/scaleway-backup