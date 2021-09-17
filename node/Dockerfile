FROM node
ADD app.js /app/
ADD package.json /app/
ADD package-lock.json /app/
WORKDIR /app
RUN npm install
CMD ["node", "app.js"]
