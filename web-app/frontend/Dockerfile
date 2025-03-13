# Dependency stage
FROM node:20-alpine as deps

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json ./package.json
COPY yarn.lock ./yarn.lock

# Install dependencies
RUN yarn install --frozen-lockfile

# Build stage
FROM node:20-alpine as builder

# Set working directory
WORKDIR /app

# Copy dependencies from previous stage
COPY --from=deps /app/node_modules ./node_modules

# Copy source code
COPY . .

# Build the application
RUN yarn build

# Production stage
FROM node:20-alpine

# Define build arguments
ARG BACKEND_URL

# Set environment variables
ENV BACKEND_URL=$BACKEND_URL

WORKDIR /app

# Copy necessary files from builder
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

# Expose port 3000
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]
