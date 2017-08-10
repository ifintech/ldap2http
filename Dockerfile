FROM alpine:3.6

# File Author / Maintainer
MAINTAINER lvyalin lvyalin.yl@gmail.com

COPY ./bin/ldap-auth /usr/local/bin/ldap-auth
RUN chmod +x /usr/local/bin/ldap-auth

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]