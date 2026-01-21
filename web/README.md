# MCP Task Tracker - Vue.js UI

Modern, responsive Vue.js frontend for the MCP Task Tracker application.

## Features

- ✨ **Modern Vue 3** with Composition API
- 🎨 **Clean, responsive design** that works on all devices
- 🚀 **Fast development** with Vite
- 📦 **State management** with Pinia
- 🔄 **Real-time updates** via WebSocket
- 🛣️ **Client-side routing** with Vue Router
- 📡 **API integration** with Axios

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **Vite** - Next generation frontend tooling
- **Vue Router** - Official router for Vue.js
- **Pinia** - State management for Vue
- **Axios** - Promise-based HTTP client

## Project Structure

```
web/
├── index.html              # HTML entry point
├── package.json            # Dependencies and scripts
├── vite.config.js         # Vite configuration
├── src/
│   ├── main.js            # Application entry point
│   ├── App.vue            # Root component
│   ├── assets/
│   │   └── styles.css     # Global styles
│   ├── router/
│   │   └── index.js       # Route definitions
│   ├── stores/
│   │   ├── projectStore.js  # Project state management
│   │   ├── agentStore.js    # Agent state management
│   │   └── taskStore.js     # Task state management
│   ├── services/
│   │   └── api.js         # API client
│   ├── composables/
│   │   └── useWebSocket.js  # WebSocket composable
│   └── views/
│       ├── Dashboard.vue   # Dashboard page
│       ├── Projects.vue    # Projects list
│       ├── ProjectDetail.vue  # Project details
│       ├── Agents.vue      # Agents list
│       └── Tasks.vue       # Tasks list
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
cd web
npm install
```

### Development

```bash
# Start development server (with hot-reload)
npm run dev
```

The application will be available at http://localhost:3000

### Build for Production

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

The production build will be output to the `dist/` directory.

## Configuration

The Vite configuration includes:

- **Proxy**: Automatically proxies API requests to `http://localhost:8080`
- **WebSocket**: Proxies WebSocket connections to the backend
- **Hot Module Replacement**: Fast refresh during development

## API Integration

The app connects to the MCP backend API at:
- Development: http://localhost:8080/api (via Vite proxy)
- Production: Same origin (configured in deployment)

## Available Pages

### Dashboard
- Overview statistics
- Recent projects, agents, and tasks
- Real-time connection status

### Projects
- List all projects
- Create new projects
- View project details
- Manage project agents and tasks

### Agents
- View all registered agents
- See agent status and details
- Filter by role and status

### Tasks
- List all tasks
- Filter by status and priority
- View task details and assignments

## State Management

The app uses Pinia stores for state management:

- **projectStore** - Manages project data
- **agentStore** - Manages agent data
- **taskStore** - Manages task data

Each store provides:
- Reactive state
- Actions for data fetching and mutations
- Error handling
- Loading states

## WebSocket Integration

Real-time updates are handled via WebSocket:

```javascript
import { useWebSocket } from '@/composables/useWebSocket'

const { connect, on, send, isConnected } = useWebSocket()

// Connect to WebSocket
connect()

// Listen for events
on('task_updated', (data) => {
  console.log('Task updated:', data)
})

// Send messages
send({ type: 'subscribe', channel: 'tasks' })
```

## Styling

The app uses vanilla CSS with CSS variables for theming:

- Consistent color palette
- Responsive design
- Modern UI components
- Smooth animations

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Deployment

### Static Hosting

After building, deploy the `dist/` directory to any static hosting service:

- Netlify
- Vercel
- GitHub Pages
- AWS S3 + CloudFront
- Any static web server

### Docker

The existing Dockerfile serves both the API and static files. The Vue build can be integrated:

```dockerfile
# Build Vue.js app
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Copy built files to final image
COPY --from=frontend-builder /app/web/dist /app/web/dist
```

## Development Tips

### Hot Reload
Vite provides instant hot module replacement - save any file and see changes immediately.

### API Proxy
All `/api` requests are automatically proxied to the backend, no CORS issues during development.

### Vue DevTools
Install Vue DevTools browser extension for debugging:
- Chrome: https://chrome.google.com/webstore
- Firefox: https://addons.mozilla.org/firefox

## Troubleshooting

### Port Already in Use
If port 3000 is in use, Vite will automatically try the next available port.

### API Connection Issues
Ensure the MCP backend is running on port 8080:
```bash
docker-compose up -d
```

### WebSocket Not Connecting
Check that the backend WebSocket endpoint is accessible at `/ws`.

## Contributing

When adding new features:

1. Create new components in `src/components/`
2. Add routes in `src/router/index.js`
3. Use stores for shared state
4. Follow Vue 3 Composition API patterns
5. Keep components focused and reusable

## License

Same as the main MCP Task Tracker project.
