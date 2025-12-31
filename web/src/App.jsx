import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Devices from './pages/Devices'
import Firmware from './pages/Firmware'
import Releases from './pages/Releases'

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/devices" element={<Devices />} />
          <Route path="/firmware" element={<Firmware />} />
          <Route path="/releases" element={<Releases />} />
        </Routes>
      </Layout>
    </Router>
  )
}

export default App
