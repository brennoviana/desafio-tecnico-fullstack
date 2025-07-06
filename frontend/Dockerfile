FROM node:20-alpine

WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port 5173
EXPOSE 5173

# Run the development server
CMD ["npm", "run", "dev", "--", "--host"] 