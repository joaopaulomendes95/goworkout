FROM node:20-alpine

WORKDIR /app

COPY package*.json svelte.config.js vite.config.ts tsconfig.json ./
RUN npm install

COPY . .

RUN npm run build

EXPOSE 5173

ENV PORT=5173

CMD ["node", "build"]
