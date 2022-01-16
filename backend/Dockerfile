FROM php:7.4-apache

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
COPY api/     /var/www/html/api
COPY classes/ /var/www/html/classes
COPY config/  /var/www/html/config
RUN  chown -R www-data:www-data /var/www/html
RUN  rm -rf /var/www/html/config/config.php

# use custom php configuration
COPY etc/php.ini /usr/local/etc/php/php.ini

# create upload folder
RUN mkdir -p /var/www/html/uploads/tmp
RUN mkdir -p /var/www/html/uploads/logos
RUN chown -R www-data:www-data /var/www/html/uploads

# RUN echo "file_uploads = On"          > /usr/local/etc/php/conf.d/booksys-file-uploading.ini
# RUN echo "upload_tmp_dir = /tmp"     >> /usr/local/etc/php/conf.d/booksys-file-uploading.ini
# RUN echo "upload_max_file_size = 2M" >> /usr/local/etc/php/conf.d/booksys-file-uploading.ini
# RUN echo "max_file_uploads = 20"     >> /usr/local/etc/php/conf.d/booksys-file-uploading.ini

# booksys specific php settings
RUN echo "log_errors = On"                     > /usr/local/etc/php/conf.d/booksys-settings.ini
RUN echo "error_log = /var/log/php-error.log" >> /usr/local/etc/php/conf.d/booksys-settings.ini

# use a custom apache configuration
COPY etc/apache2/000-default.conf /etc/apache2/sites-available/000-default.conf

# Ports
EXPOSE 80