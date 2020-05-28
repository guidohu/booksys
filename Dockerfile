FROM php:7.4.4-apache

# activate mysqli extension for PHP
RUN mkdir -p /usr/local/etc/php/conf.d
RUN docker-php-ext-install mysqli && docker-php-ext-enable mysqli

# install packages to send emails
RUN apt-get update && apt-get install -y sendemail \
    libnet-ssleay-perl \
    libio-socket-ssl-perl
RUN apt-get clean

# enable CSP header
RUN a2enmod headers

# copy our UI application code
COPY src/ /var/www/html/
# use a custom apache configuration
COPY etc/apache2/000-default.conf /etc/apache2/sites-available/000-default.conf

# Ports
EXPOSE 80