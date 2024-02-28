FROM node:18-alpine AS development
ENV NODE_ENV development
# Add a work directory
WORKDIR /app
# Cache and Install dependencies
COPY package.json .
COPY package-lock.json .
RUN yarn install
# Copy app files
COPY . .
COPY env.sample .env
# Expose port
EXPOSE 3000
# Start the app
CMD [ "yarn", "start" ]

