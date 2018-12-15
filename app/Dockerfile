FROM node:8.12.0 as build-stage
WORKDIR /app
#initially copy over just package.json (and locks, if present) for dependency installation step
COPY package*.json /app/
RUN npm install
#then copy over everything, including our code, for the build step
COPY ./ /app/
RUN npm run build

#once the build-stage is complete, we can build an image for our deployment
#which is ONLY nginx(alpine) + our bundled app roughly speaking, the build
#image ends up over 1 GB while the deployable image is only about 25 MB
FROM nginx:alpine as deploy-stage
COPY --from=build-stage /app/build /usr/share/nginx/html
COPY --from=build-stage /app/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/entrypoint.sh /etc/nginx/entrypoint.sh
RUN chmod 777 /etc/nginx/entrypoint.sh
CMD /etc/nginx/entrypoint.sh