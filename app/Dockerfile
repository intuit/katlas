FROM node:8.12.0 as build-stage
WORKDIR /app
#initially copy over just package.json (and locks, if present) for dependency installation
COPY package*.json /app/
RUN yarn install
#then copy over our code for build
COPY ./ /app/
RUN yarn build

#once the build-stage is complete, we can build an image for our deployment
#which is ONLY nginx(alpine) + our bundled app roughly speaking, the build
#image ends up over 1 GB while the deployable image is only about 25 MB
FROM nginx:alpine
COPY --from=build-stage /app/build /usr/share/nginx/html
COPY --from=build-stage /app/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/entrypoint.sh /etc/nginx/entrypoint.sh
RUN chmod 777 /etc/nginx/entrypoint.sh
EXPOSE 80
CMD /etc/nginx/entrypoint.sh