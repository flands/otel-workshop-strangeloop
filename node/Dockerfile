FROM node
ADD app.js /app/
ADD tracing.js /app/
ADD package.json /app/
ADD package-lock.json /app/
WORKDIR /app
RUN npm install
CMD ["node", "-r", "./tracing.js", "app.js"]
