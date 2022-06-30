FROM nginx:1.23
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx/conf.d /etc/nginx/conf.d
