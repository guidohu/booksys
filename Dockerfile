FROM php:7.4.4-apache

COPY src/ /var/www/html/
COPY etc/apache2/000-default.conf /etc/apache2/sites-available/000-default.conf

RUN if [ -f /var/www/html/config/config.php ]; then rm /var/www/html/config/config.php; fi
RUN mkdir -p /usr/local/etc/php/conf.d
RUN docker-php-ext-install mysqli && docker-php-ext-enable mysqli

RUN apt-get update && apt-get install -y sendemail
RUN apt-get clean

# Ports
EXPOSE 80