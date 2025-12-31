# AuraX Web Dashboard

Modern React dashboard for AuraX IoT Fleet Management Platform.

## Features

- 📊 Real-time fleet statistics
- 🔌 Device management
- 📦 Firmware upload and versioning
- 🚀 OTA release management
- 📈 Health monitoring

## Tech Stack

- React 18
- Vite
- Tailwind CSS
- Axios
- React Router
- Lucide Icons

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Environment

The dashboard expects the API server to be running on `http://localhost:8080`.

Configure proxy in `vite.config.js` if needed.

## Pages

- **Dashboard**: Overview with fleet statistics
- **Devices**: Device list with search and creation
- **Firmware**: Firmware upload and version management
- **Releases**: OTA release creation and monitoring

## API Integration

All API calls are made through the service layer in `src/services/api.js`.

The dashboard communicates with the AuraX backend API.
